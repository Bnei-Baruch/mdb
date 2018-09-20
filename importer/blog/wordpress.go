package blog

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/pkg/errors"
	"github.com/robbiet480/go-wordpress"
	"golang.org/x/oauth2"

	"github.com/Bnei-Baruch/mdb/version"
)

func NewWordpressClient(wpUrl, username, password string) (*wordpress.Client, error) {
	// get a fresh JWT access token
	jwtToken, err := getJWTToken(wpUrl, username, password)
	if err != nil {
		return nil, errors.Wrap(err, "getJWTToken")
	}

	// Create HttpClient
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: jwtToken},
	)
	tc := oauth2.NewClient(ctx, ts)

	client, err := wordpress.NewClient(fmt.Sprintf("%swp-json/", wpUrl), tc)
	if err != nil {
		return nil, errors.Wrap(err, "wordpress.NewClient")
	}

	return client, nil
}

func getJWTToken(wpUrl, username, password string) (string, error) {
	vals := url.Values{
		"username": []string{username},
		"password": []string{password},
	}

	req, err := http.NewRequest(http.MethodPost,
		fmt.Sprintf("%swp-json/jwt-auth/v1/token", wpUrl),
		strings.NewReader(vals.Encode()))
	if err != nil {
		return "", errors.Wrap(err, "http.NewRequest")
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", fmt.Sprintf("MDB_%s", version.Version))
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", errors.Wrap(err, "http.Do")
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", errors.Wrap(err, "ioutil.ReadAll body")
	}

	if resp.StatusCode != http.StatusOK {
		log.Warnf("wordpress.getJWTToken bad status code [%d]\n%s\n", resp.StatusCode, string(body))
	}

	var bodyJson map[string]interface{}
	err = json.Unmarshal(body, &bodyJson)
	if err != nil {
		return "", errors.Wrap(err, "json.Unmarshal")
	}

	if v, ok := bodyJson["token"]; ok {
		if vv, ok := v.(string); ok {
			return vv, nil
		} else {
			return "", errors.Errorf("JWT token is not a string, got: %v", v)
		}
	} else {
		return "", errors.Errorf("Invalid JWT token endpoint response: %v", bodyJson)
	}
}
