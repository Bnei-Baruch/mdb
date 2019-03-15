package twitter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"time"

	"github.com/ChimeraCoder/anaconda"
	log "github.com/Sirupsen/logrus"
	"github.com/pkg/errors"
	"github.com/volatiletech/sqlboiler/queries/qm"
	"gopkg.in/volatiletech/null.v6"

	"github.com/Bnei-Baruch/mdb/models"
	"github.com/Bnei-Baruch/mdb/utils"
)

var dumps = map[string]string{
	"Michael_Laitman": "importer/twitter/data/twitter-2018-07-04-Michael_Laitman",
	"laitman_co_il":   "importer/twitter/data/twitter-2018-07-03-laitman_co_il",
	"laitman":         "importer/twitter/data/twitter-2018-07-05-laitman",
	"laitman_es":      "importer/twitter/data/twitter-2018-07-05-laitman_es",
}

func ImportDumps() {
	clock, _ := Init()

	for k, v := range dumps {
		utils.Must(importDump(k, v))
	}

	Shutdown()
	log.Info("Success")
	log.Infof("Total run time: %s", time.Now().Sub(clock).String())
}

func importDump(username, dir string) error {
	err := cleanTwitterDump(dir)
	if err != nil {
		return errors.Wrapf(err, "clean dump: %s", username)
	}

	tweets, err := readTwitterDump(dir)
	if err != nil {
		return errors.Wrapf(err, "clean dump: %s", username)
	}
	log.Infof("%s has %d tweets", username, len(tweets))

	user, err := models.TwitterUsers(mdb, qm.Where("username = ?", username)).One()
	if err != nil {
		return errors.Wrapf(err, "lookup username: %s", username)
	}

	for i := range tweets {
		if err := saveTweetToDB(tweets[i], user); err != nil {
			return errors.Wrapf(err, "Save tweet to DB: %s %d", username, i)
		}
	}

	return nil
}

func saveTweetToDB(t *anaconda.Tweet, user *models.TwitterUser) error {
	ts, err := t.CreatedAtTime()
	if err != nil {
		return errors.Wrapf(err, "Tweet.CreatedAtTime()")
	}

	jsonb, err := json.Marshal(t)
	if err != nil {
		return errors.Wrapf(err, "json.Marshal")
	}

	tx, err := mdb.Begin()
	utils.Must(err)

	mt := models.TwitterTweet{
		UserID:    user.ID,
		TwitterID: t.IdStr,
		FullText:  t.FullText,
		TweetAt:   ts,
		Raw:       null.JSONFrom(jsonb),
	}

	err = mt.Upsert(tx, true, []string{"twitter_id"}, []string{"full_text", "tweet_at", "raw"})
	if err != nil {
		utils.Must(tx.Rollback())
		return errors.Wrapf(err, "Upsert to DB")
	} else {
		utils.Must(tx.Commit())
	}

	return nil
}

func Analyze() {
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})

	dir := dumps["laitman_es"]
	err := cleanTwitterDump(dir)
	utils.Must(err)

	tweets, err := readTwitterDump(dir)
	utils.Must(err)

	// sort tweets by created_at
	sort.Slice(tweets, func(i, j int) bool {
		tsA, err := tweets[i].CreatedAtTime()
		if err != nil {
			log.Errorf("%d\tTime.Parse Error %s\n", i, tweets[i].CreatedAt)
			return true
		}

		tsB, err := tweets[j].CreatedAtTime()
		if err != nil {
			log.Errorf("%d\tTime.Parse Error %s\n", j, tweets[j].CreatedAt)
			return false
		}

		return tsB.Before(tsA)
	})

	log.Infof("Tweets: %d\n", len(tweets))

	hashtags := make(map[string]int)
	langs := make(map[string]int)
	wMedia := make([]*anaconda.Tweet, 0)
	wUrls := make([]*anaconda.Tweet, 0)
	wUserMentions := make([]*anaconda.Tweet, 0)
	for i, t := range tweets {
		ts, err := t.CreatedAtTime()
		if err != nil {
			fmt.Printf("%d\tTime.Parse Error %s\n", i, t.CreatedAt)
			continue
		}

		//if i > 2000 {
		//	continue
		//}

		langs[t.Lang]++

		for _, ht := range t.Entities.Hashtags {
			hashtags[ht.Text]++
		}

		if len(t.Entities.Media) > 0 {
			wMedia = append(wMedia, t)
		}
		if len(t.Entities.Urls) > 0 {
			wUrls = append(wUrls, t)
		}
		if len(t.Entities.User_mentions) > 0 {
			wUserMentions = append(wUserMentions, t)
		}

		fmt.Printf("%d\t%s\t%d\t%d\t%s\n", i, ts.Format(time.RFC3339), t.FavoriteCount, t.RetweetCount, t.FullText)
	}

	log.Infof("HashTags histogram has %d keys", len(hashtags))
	for k, v := range hashtags {
		log.Infof("%s\t%d", k, v)
	}
	log.Infof("Langs histogram has %d keys", len(langs))
	for k, v := range langs {
		log.Infof("%s\t%d", k, v)
	}

	log.Infof("%d tweets with media", len(wMedia))
	for i := range wMedia {
		log.Info(wMedia[i].IdStr)
	}
	log.Infof("%d tweets with Urls", len(wUrls))
	for i := range wUrls {
		log.Info(wUrls[i].IdStr)
	}
	log.Infof("%d tweets with User Mentions", len(wUserMentions))
	for i := range wUserMentions {
		log.Info(wUserMentions[i].IdStr)
	}
}

func cleanTwitterDump(dir string) error {
	// read original file
	raw, err := ioutil.ReadFile(filepath.Join(dir, "tweet.js"))
	if err != nil {
		return errors.Wrap(err, "ioutil.ReadFile")
	}
	raw = bytes.TrimPrefix(raw, []byte("window.YTD.tweet.part0 = "))

	var data []map[string]interface{}
	err = json.Unmarshal(raw, &data)
	if err != nil {
		return errors.Wrap(err, "json.Unmarshal")
	}

	// mutate data
	log.Infof("Original tweets %d", len(data))
	for i, item := range data {
		if err := cleanRawArchiveTweet(item); err != nil {
			log.Warnf("Error cleaning tweet %d: %s", i, err.Error())
		}
	}

	// save results to clean file

	f, err := os.OpenFile(filepath.Join(dir, "tweet_clean.js"), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return errors.Wrap(err, "os.OpenFile")
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	err = enc.Encode(data)
	if err != nil {
		return errors.Wrap(err, "json.Encode")
	}

	return nil
}

func strFloatToInt(s string) (int, error) {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, errors.Wrap(err, "ParseFloat")
	}
	return int(f), nil
}

func strFloatToInt64(s string) (int64, error) {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, errors.Wrap(err, "ParseFloat")
	}
	return int64(f), nil
}

// Twitter archive data dump does not conform to the Twitter API v1.1
// So we have to clean up a lot of mess in the json
func cleanRawArchiveTweet(item map[string]interface{}) error {
	for k, v := range item {
		var clean interface{}
		var err error
		switch k {
		case "display_text_range", "indices", "aspect_ratio":
			vval := v.([]interface{})
			nVal := make([]int, len(vval))
			for i := range vval {
				nVal[i], err = strFloatToInt(vval[i].(string))
			}
			clean = nVal

		case "favorite_count",
			"retweet_count",
			"w",
			"h",
			"bitrate",
			"size":
			if vv, ok := v.(string); ok {
				clean, err = strFloatToInt(vv)
			}

		case "id",
			"in_reply_to_user_id",
			"in_reply_to_status_id",
			"quoted_status_id",
			"duration_millis",
			"source_status_id":
			if vv, ok := v.(string); ok {
				clean, err = strFloatToInt64(vv)
			}

		case "entities",
			"extended_entities",
			"extended_tweet",
			"quoted_status",
			"retweeted_status",
			"sizes",
			"medium",
			"thumb",
			"small",
			"large",
			"video_info":
			err = cleanRawArchiveTweet(v.(map[string]interface{}))
			clean = v

		case "urls", "hashtags", "user_mentions", "media", "variants":
			vval := v.([]interface{})
			for i := range vval {
				err = cleanRawArchiveTweet(vval[i].(map[string]interface{}))
			}
			clean = v

		case "coordinates":
			if vv, ok := v.(map[string]interface{}); ok { // struct
				err = cleanRawArchiveTweet(vv)
				clean = v
			} else if vv, ok := v.([]interface{}); ok { // []float64
				nVal := make([]float64, len(vv))
				for i := range vv {
					nVal[i], err = strconv.ParseFloat(vv[i].(string), 64)
				}
				clean = nVal
			}
		default:
			continue
		}

		if err != nil {
			return errors.Wrapf(err, "%s => %s", k, v)
		} else {
			item[k] = clean
		}
	}

	return nil
}

func readTwitterDump(dir string) ([]*anaconda.Tweet, error) {
	raw, err := ioutil.ReadFile(filepath.Join(dir, "tweet_clean.js"))
	if err != nil {
		return nil, errors.Wrap(err, "ioutil.ReadFile")
	}

	var tweets []*anaconda.Tweet
	err = json.Unmarshal(raw, &tweets)
	if err != nil {
		return nil, errors.Wrap(err, "json.Unmarshal")
	}

	return tweets, nil
}
