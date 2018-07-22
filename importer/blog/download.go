package blog

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/pkg/errors"
	"github.com/robbiet480/go-wordpress"
	"github.com/spf13/viper"

	"github.com/Bnei-Baruch/mdb/utils"
)

func Download() {
	clock := Init()

	utils.Must(doDownload())

	Shutdown()
	log.Info("Success")
	log.Infof("Total run time: %s", time.Now().Sub(clock).String())
}

func doDownload() error {
	wpConfig := viper.GetStringMapString(fmt.Sprintf("wordpress.%s", blogName))
	client, err := NewWordpressClient(wpConfig["url"], wpConfig["username"], wpConfig["password"])
	if err != nil {
		return errors.Wrap(err, "NewWordpressClient")
	}

	page := 1
	perPage := 100
	skipCount := 0
	for {
		log.Infof("Page %d [%d skipped]", page, skipCount)
		posts, resp, err := client.Posts.List(context.Background(), &wordpress.PostListOptions{
			ListOptions: wordpress.ListOptions{
				Page:    page,
				PerPage: perPage,
			},
		})
		if err != nil {
			return errors.Wrapf(err, "Posts.List %d", page)
		}

		for _, post := range posts {
			if LESSON_RE.MatchString(post.Title.Rendered) ||
				//CLIP_RE.MatchString(post.Title.Rendered) ||
				TWITTER_RE.MatchString(post.Title.Rendered) ||
				DECLAMATION_RE.MatchString(post.Title.Rendered) {
				skipCount++
			}

			if err := saveToFile(post); err != nil {
				log.Errorf("saveToFile %d: %s", post.ID, err.Error())
			}
		}

		page = resp.NextPage
		if page < 1 {
			break
		}
	}

	return nil
}

func saveToFile(post *wordpress.Post) error {
	// create output file
	idStr := strconv.Itoa(post.ID)
	dir := fmt.Sprintf("importer/blog/data/%s/%s", blogName, idStr[len(idStr)-2:])
	if err := os.MkdirAll(dir, os.FileMode(0755)); err != nil {
		return errors.Wrapf(err, "os.MkdirAll %s", dir)
	}

	path := fmt.Sprintf("%s/%s.json", dir, idStr)
	output, err := os.Create(path)
	if err != nil {
		return errors.Wrapf(err, "os.Create %s", path)
	}
	defer output.Close()

	err = json.NewEncoder(output).Encode(post)
	if err != nil {
		return errors.Wrap(err, "json.Encode")
	}

	return nil
}
