package blog

import (
	"encoding/json"
	"net/url"
	"os"
	"regexp"
	"sort"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/pkg/errors"
	"github.com/robbiet480/go-wordpress"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
	"jaytaylor.com/html2text"

	"github.com/Bnei-Baruch/mdb/utils"
)

var LESSON_RE = regexp.MustCompile("(?i)^(Уроки и лекции|Утренний урок|Вечерний урок|Урок по Книге Зоар|Lección diaria de Cabalá|Lección diaria de Cábala|Lección diaria de la Cabalá|Daily Kabbalah Lesson|Evening Zohar Lesson|שיעור הקבלה היומי|שיעור וירטואלי|שיעור זוהר|שיעורי הקבלה, סדנאות חיבור ושיחות)")
var CLIP_RE = regexp.MustCompile("(?i)клип")
var TWITTER_RE = regexp.MustCompile("(?i)(Мои мысли в Twitter|Mis pensamientos en Twitter|Mis pensamiento en Twitter|My Thoughts On Twitter|המחשבות שלי ב &#8211; Twitter)")
var DECLAMATION_RE = regexp.MustCompile("(?i)^(Радио-версия|Audio Version Of The Blog|גרסת אודיו)")
var PROGRAMS_RE = regexp.MustCompile("(?i)^(Una nueva vida|Una vida nueva|Good Environment|New Life|חיים חדשים|טעימות משיעור הקבלה היומי)")

func Analyze() {
	clock := Init()

	utils.Must(doAnalyze())

	Shutdown()
	log.Info("Success")
	log.Infof("Total run time: %s", time.Now().Sub(clock).String())
}

func doAnalyze() error {
	cssClasses := make(map[string]int64)
	linkHosts := make(map[string]int64)
	tagCounts := make(map[string]int64)
	titleGroups := make(map[string][]string)
	termsWCMap := make(map[string]*BucketHistorgram)

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

		ctxNode := html.Node{
			Type:     html.ElementNode,
			DataAtom: atom.Body,
			Data:     "body",
		}
		nodes, err := html.ParseFragment(strings.NewReader(post.Content.Rendered), &ctxNode)
		if err != nil {
			return errors.Wrapf(err, "html.Parse %d", post.ID)
		}
		for i := range nodes {
			//html.Render(os.Stdout, nodes[i])
			node := nodes[i]
			mergeHist(cssClasses, extractCssClasses(node))
			mergeHist(linkHosts, linkAnalysis(node))
			mergeHist(tagCounts, tagAnalysis(node))
		}

		// *** Title Groups
		title := post.Title.Rendered
		if LESSON_RE.MatchString(title) {
			titleGroups["lesson"] = append(titleGroups["lesson"], title)
		} else if CLIP_RE.MatchString(title) {
			titleGroups["clip"] = append(titleGroups["clip"], title)
		} else if TWITTER_RE.MatchString(title) {
			titleGroups["twitter"] = append(titleGroups["twitter"], title)
		} else if DECLAMATION_RE.MatchString(title) {
			titleGroups["declamation"] = append(titleGroups["declamation"], title)
		} else if PROGRAMS_RE.MatchString(title) {
			titleGroups["programs"] = append(titleGroups["programs"], title)
		} else {
			titleGroups["other"] = append(titleGroups["other"], title)
		}

		return nil
	}

	err := traverse(walkFn)
	if err != nil {
		return errors.Wrap(err, "traverse error")
	}

	log.Infof("cssClasses has %d entries", len(cssClasses))
	type KV struct {
		k string
		v int64
	}
	kvs := make([]KV, 0)
	for k, v := range cssClasses {
		kvs = append(kvs, KV{k: k, v: v})
		//log.Infof("%s\t%d", k, v)
	}
	sort.Slice(kvs, func(i, j int) bool {
		return kvs[i].v < kvs[j].v
	})
	for i := range kvs {
		log.Infof("%s\t%d", kvs[i].k, kvs[i].v)
	}

	log.Info("\n\n\n\n\n\n\n\n\n")
	log.Infof("linkHosts has %d entries", len(linkHosts))
	kvs = make([]KV, 0)
	for k, v := range linkHosts {
		kvs = append(kvs, KV{k: k, v: v})
		//log.Infof("%s\t%d", k, v)
	}
	sort.Slice(kvs, func(i, j int) bool {
		return kvs[i].v < kvs[j].v
	})
	for i := range kvs {
		log.Infof("%s\t%d", kvs[i].k, kvs[i].v)
	}

	log.Info("\n\n\n\n\n\n\n\n\n")
	log.Infof("tagCounts has %d entries", len(tagCounts))
	kvs = make([]KV, 0)
	for k, v := range tagCounts {
		kvs = append(kvs, KV{k: k, v: v})
		//log.Infof("%s\t%d", k, v)
	}
	sort.Slice(kvs, func(i, j int) bool {
		return kvs[i].v < kvs[j].v
	})
	for i := range kvs {
		log.Infof("%s\t%d", kvs[i].k, kvs[i].v)
	}

	log.Info("\n\n\n\n\n\n\n\n\n")
	for k, v := range titleGroups {
		log.Infof("%s\t%d", k, len(v))
	}

	//groups := []string{"lesson", "clip", "twitter", "declamation", "programs", "other"}
	//for i := range groups {
	//	titles := titleGroups[groups[i]]
	//	sort.Slice(titles, func(i, j int) bool {
	//		return titles[i] < titles[j]
	//	})
	//	for j := range titles {
	//		log.Info(titles[j])
	//	}
	//}

	log.Info("\n\n\n\n\n\n\n\n\n")
	sum := 0
	for k, v := range termsWCMap {
		log.Infof("Term WC Histogram: %s [%d]", k, len(v.buckets))
		v.Dump()
		psum := 0
		for _, vv := range v.buckets {
			psum += vv
		}
		log.Infof("psum: %d", psum)
		sum += psum
	}
	log.Infof("total: %d", sum)

	return nil
}

func extractCssClasses(node *html.Node) map[string]int64 {
	classMap := make(map[string]int64)

	for c := node.FirstChild; c != nil; c = c.NextSibling {
		mergeHist(classMap, extractCssClasses(c))
	}

	for _, a := range node.Attr {
		if a.Key == "class" {
			classes := strings.Fields(strings.TrimSpace(a.Val))
			for i := range classes {
				classMap[classes[i]]++
			}
			break
		}
	}

	return classMap
}

func linkAnalysis(node *html.Node) map[string]int64 {
	hostsMap := make(map[string]int64)

	for c := node.FirstChild; c != nil; c = c.NextSibling {
		mergeHist(hostsMap, linkAnalysis(c))
	}

	if node.DataAtom == atom.A {
		for _, a := range node.Attr {
			if a.Key == "href" {
				pUrl, err := url.Parse(strings.TrimSpace(a.Val))
				if err != nil {
					log.Errorf("url.Parse: %s", err.Error())
					break
				}

				hostsMap[pUrl.Host]++
				break
			}
		}
	}

	return hostsMap
}

func tagAnalysis(node *html.Node) map[string]int64 {
	tagMap := make(map[string]int64)

	for c := node.FirstChild; c != nil; c = c.NextSibling {
		mergeHist(tagMap, tagAnalysis(c))
	}

	tagMap[node.DataAtom.String()]++

	return tagMap
}

func mergeHist(a, b map[string]int64) {
	for k, v := range b {
		a[k] += v
	}
}

func htmlWordCount(inputHTML string) (int, error) {
	text, err := html2text.FromString(inputHTML, html2text.Options{OmitLinks: true})
	if err != nil {
		return 0, errors.Wrap(err, "html2text.FromString")
	}
	//log.Info(text)
	words := strings.Fields(text)
	return len(words), nil
}

type BucketHistorgram struct {
	bucketSize int
	buckets    map[int]int
}

func NewBucketHistorgram(bucketSize int) *BucketHistorgram {
	return &BucketHistorgram{
		bucketSize: bucketSize,
		buckets:    make(map[int]int),
	}
}

func (h *BucketHistorgram) Add(val int) {
	h.buckets[val%h.bucketSize]++
}

func (h *BucketHistorgram) Dump() {
	keys := make([]int, len(h.buckets))
	i := 0
	for k := range h.buckets {
		keys[i] = k
		i++
	}
	sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })
	for i := range keys {
		log.Infof("%d\t%d", keys[i], h.buckets[keys[i]])
	}

}

type BlogPostFilter interface {
	IsPass(*wordpress.Post) bool
}

type PassFilter struct{}

func (f *PassFilter) IsPass(*wordpress.Post) bool {
	return true
}

type TitleBasedFilter struct{}

func (f *TitleBasedFilter) IsPass(post *wordpress.Post) bool {
	title := post.Title.Rendered
	return !(LESSON_RE.MatchString(title) ||
		// CLIP_RE.MatchString(post.Title.Rendered) ||
		TWITTER_RE.MatchString(title) ||
		PROGRAMS_RE.MatchString(title) ||
		DECLAMATION_RE.MatchString(title))
}

var titleBasedFilter = new(TitleBasedFilter)

func getBlogPostFilter(blogID int64) BlogPostFilter {
	return titleBasedFilter
}
