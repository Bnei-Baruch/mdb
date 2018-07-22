package blog

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/pkg/errors"
	"github.com/robbiet480/go-wordpress"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"

	"github.com/Bnei-Baruch/mdb/models"
	"github.com/Bnei-Baruch/mdb/utils"
)

func Import() {
	clock := Init()

	utils.Must(doImport())

	Shutdown()
	log.Info("Success")
	log.Infof("Total run time: %s", time.Now().Sub(clock).String())
}

func doImport() error {
	skipCount := 0

	walkFn := func(path string, info os.FileInfo, err error) error {
		f, err := os.Open(path)
		if err != nil {
			return errors.Wrapf(err, "os.Open %s", path)
		}

		var post wordpress.Post
		err = json.NewDecoder(f).Decode(&post)
		if err != nil {
			return errors.Wrapf(err, "json.Decode %s", path)
		}

		if LESSON_RE.MatchString(post.Title.Rendered) ||
			//CLIP_RE.MatchString(post.Title.Rendered) ||
			TWITTER_RE.MatchString(post.Title.Rendered) ||
			DECLAMATION_RE.MatchString(post.Title.Rendered) {
			skipCount++
			return nil
		}

		//_, err = prepare(&post)
		blogPost, err := prepare(&post)
		if err != nil {
			return errors.Wrapf(err, "prepare %s", path)
		}

		blogPost.BlogID = 1
		err = blogPost.Insert(mdb)
		if err != nil {
			return errors.Wrapf(err, "blogPost.Insert %s", path)
		}

		return nil
	}

	err := traverse(walkFn)
	if err != nil {
		fmt.Printf("traverse error: %v\n", err)
	}

	log.Infof("skipCount: %d", skipCount)

	return nil
}

func prepare(post *wordpress.Post) (*models.BlogPost, error) {
	ctxNode := html.Node{
		Type:     html.ElementNode,
		DataAtom: atom.Body,
		Data:     "body",
	}
	var sb strings.Builder

	nodes, err := html.ParseFragment(strings.NewReader(post.Content.Rendered), &ctxNode)
	if err != nil {
		return nil, errors.Wrapf(err, "html.Parse %d", post.ID)
	}
	for i := range nodes {
		node := nodes[i]

		// skip script nodes
		if isScript(node) {
			continue
		}

		// stop at share div (usually the last node)
		if isShareDaddy(node) {
			break
		}

		cleanNode(node)
		html.Render(&sb, node)
	}

	bp := &models.BlogPost{
		WPID:     int64(post.ID),
		PostedAt: post.DateGMT.Time,
		Title:    post.Title.Rendered,
		Content:  sb.String(),
	}

	return bp, nil
}

func isShareDaddy(node *html.Node) bool {
	if node.DataAtom != atom.Div {
		return false
	}

	for _, a := range node.Attr {
		if a.Key == "class" {
			classes := strings.Fields(strings.TrimSpace(a.Val))
			for i := range classes {
				if classes[i] == "sharedaddy" {
					return true
				}
			}
			break
		}
	}

	return false
}

func isScript(node *html.Node) bool {
	return node.DataAtom == atom.Script
}

func cleanNode(node *html.Node) {
	// clean children
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		if isScript(c) {
			node.RemoveChild(c)
		}
	}
}
