package twitter

import (
	"encoding/json"
	"net/url"
	"time"

	"github.com/ChimeraCoder/anaconda"
	log "github.com/Sirupsen/logrus"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"github.com/volatiletech/sqlboiler/queries"

	"github.com/Bnei-Baruch/mdb/api"
	"github.com/Bnei-Baruch/mdb/events"
	"github.com/Bnei-Baruch/mdb/models"
	"github.com/Bnei-Baruch/mdb/utils"
)

func ImportLatestTweets() {
	clock, emitter := Init()

	utils.Must(importLatestTweets(emitter))

	Shutdown()
	log.Info("Success")
	log.Infof("Total run time: %s", time.Now().Sub(clock).String())
}

func importLatestTweets(emitter *events.BufferedEmitter) error {
	// initialize twitter api
	accessToken := viper.GetString("twitter.access-token")
	accessTokenSecret := viper.GetString("twitter.access-token-secret")
	consumerKey := viper.GetString("twitter.consumer-key")
	consumerSecret := viper.GetString("twitter.consumer-secret")
	twitterAPI := anaconda.NewTwitterApiWithCredentials(accessToken, accessTokenSecret, consumerKey, consumerSecret)
	defer twitterAPI.Close()

	// fetch last imported tweet ID per account
	sinceIDs, err := getSinceIDs()
	if err != nil {
		return errors.Wrap(err, "getSinceIDs")
	}

	// process each account
	for k, v := range sinceIDs {
		log.Infof("Fetching user timeline for [%s, %s]", k, v)

		user := api.TWITTER_USERS_REGISTRY.ByUsername[k]
		timeline, err := twitterAPI.GetUserTimeline(url.Values{
			"user_id":  []string{user.AccountID},
			"since_id": []string{v},
		})
		if err != nil {
			return errors.Wrapf(err, "twitterAPI.GetUserTimeline [%s, %s]", k, v)
		}

		log.Infof("%s has %d new tweets in his timeline", k, len(timeline))
		for i := range timeline {
			if err := saveTweetToDB(&timeline[i], user); err != nil {
				log.Errorf("Error saving tweet to DB: %s", err.Error())

				jsonb, err := json.Marshal(timeline[i])
				if err != nil {
					return errors.Wrapf(err, "json.Marshal")
				}
				log.Warn(string(jsonb))
			} else {
				emitter.Emit(events.TweetCreateEvent(&models.TwitterTweet{TwitterID: timeline[i].IdStr}))
			}

		}
	}

	return nil
}

func getSinceIDs() (map[string]string, error) {
	rows, err := queries.Raw(mdb, `select distinct on (u.id)
  u.username,
  t.twitter_id
from twitter_tweets t
  inner join twitter_users u on t.user_id = u.id
order by u.id, t.tweet_at desc`).Query()
	if err != nil {
		return nil, errors.Wrap(err, "queries.Raw")
	}
	defer rows.Close()

	m := make(map[string]string)
	for rows.Next() {
		var username, sinceID string
		if err := rows.Scan(&username, &sinceID); err != nil {
			return nil, errors.Wrap(err, "rows.Scan")
		}
		m[username] = sinceID
	}
	if err := rows.Err(); err != nil {
		return nil, errors.Wrap(err, "rows.Err")
	}

	return m, nil
}
