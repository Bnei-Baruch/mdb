package blog

import (
	"context"
	"fmt"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/pkg/errors"
	"github.com/robbiet480/go-wordpress"
	"github.com/spf13/viper"
	"github.com/volatiletech/sqlboiler/queries"
	"github.com/volatiletech/sqlboiler/queries/qm"

	"github.com/Bnei-Baruch/mdb/models"
	"github.com/Bnei-Baruch/mdb/utils"
)

func ImportLatest() {
	clock := Init()

	utils.Must(importLatest())

	Shutdown()
	log.Info("Success")
	log.Infof("Total run time: %s", time.Now().Sub(clock).String())
}

func importLatest() error {
	// load blogs
	blogs, err := models.Blogs(mdb).All()
	if err != nil {
		return errors.Wrap(err, "load blogs")
	}

	// fetch last imported posted_at per blog
	lastImported, err := getLastImported()
	if err != nil {
		return errors.Wrap(err, "getLastImported")
	}
	log.Infof("lastImported")
	for k, v := range lastImported {
		log.Infof("%d => %s", k, v.Format(time.RFC3339))
	}

	// import latest in each blog
	for i := range blogs {
		b := blogs[i]
		lastTS, ok := lastImported[b.ID]
		if !ok {
			log.Infof("skipping %s", b.Name)
			continue
		}

		if err := importLastFromBlog(b, lastTS); err != nil {
			log.Errorf("importLastFromBlog %s: %s", b.Name, err.Error())
		}
	}

	return nil
}

func importLastFromBlog(b *models.Blog, lastTS time.Time) error {
	log.Infof("Importing latest posts from %s [%s]", b.Name, lastTS.Format(time.RFC3339))

	wpConfig := viper.GetStringMapString(fmt.Sprintf("wordpress.%s", b.Name))
	client, err := NewWordpressClient(wpConfig["url"], wpConfig["username"], wpConfig["password"])
	if err != nil {
		return errors.Wrap(err, "NewWordpressClient")
	}

	after := lastTS.AddDate(0, 0, -3)
	postFilter := getBlogPostFilter(b.ID)

	page := 1
	perPage := 100
	skipCount := 0
	newPosts := make([]*models.BlogPost, 0)
	for {
		log.Infof("Page %d [%d skipped]", page, skipCount)
		posts, resp, err := client.Posts.List(context.Background(), &wordpress.PostListOptions{
			ListOptions: wordpress.ListOptions{
				Page:    page,
				PerPage: perPage,
			},
			After: &after,
		})
		if err != nil {
			return errors.Wrapf(err, "Posts.List %d", page)
		}

		for _, post := range posts {

			exist, err := models.BlogPosts(mdb,
				qm.Where("blog_id = ? and wp_id = ?", b.ID, post.ID)).
				Exists()
			if err != nil {
				log.Errorf("Check exists %d %d: %s", b.ID, post.ID, err.Error())
				continue
			}
			if exist {
				log.Infof("Post exists %d %d. Skipping", b.ID, post.ID)
				skipCount++
				continue
			}

			blogPost, err := prepare(post)
			if err != nil {
				log.Errorf("Prepare post %d %d: %s", b.ID, post.ID, err.Error())
				continue
			}

			log.Infof("Insert new post %s [%d]", post.Title.Rendered, post.ID)
			blogPost.BlogID = b.ID
			blogPost.Filtered = !postFilter.IsPass(post)
			err = blogPost.Insert(mdb)
			if err != nil {
				log.Errorf("Insert post to DB %d %d: %s", b.ID, post.ID, err.Error())
				continue
			}

			newPosts = append(newPosts, blogPost)
		}

		page = resp.NextPage
		if page < 1 {
			break
		}
	}

	// make relative links of new posts
	err = loadLinkMap()
	if err != nil {
		return errors.Wrap(err, "loadLinkMap")
	}

	for i := range newPosts {
		if err := cleanPost(newPosts[i]); err != nil {
			log.Errorf("cleanPost %d %d: %s", b.ID, newPosts[i].ID, err.Error())
			continue
		}
	}

	return nil
}

func getLastImported() (map[int64]time.Time, error) {
	rows, err := queries.Raw(mdb, `select distinct on (b.id)
  b.id,
  p.posted_at
from blog_posts p
  inner join blogs b on p.blog_id = b.id
order by b.id, p.posted_at desc`).Query()
	if err != nil {
		return nil, errors.Wrap(err, "queries.Raw")
	}
	defer rows.Close()

	m := make(map[int64]time.Time)
	for rows.Next() {
		var id int64
		var postedAt time.Time
		if err := rows.Scan(&id, &postedAt); err != nil {
			return nil, errors.Wrap(err, "rows.Scan")
		}
		m[id] = postedAt
	}
	if err := rows.Err(); err != nil {
		return nil, errors.Wrap(err, "rows.Err")
	}

	return m, nil
}
