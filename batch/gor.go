package batch

import (
	"bufio"
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/pkg/errors"
	"github.com/spf13/viper"

	"github.com/Bnei-Baruch/mdb/api"
	"github.com/Bnei-Baruch/mdb/utils"
)

// A utility file to scan requests.log, filter relevant requests and re-run them

type Meta struct {
	Type      int
	ID        string
	Timestamp time.Time
}

type Request struct {
	Meta     *Meta
	Payload  []string
	Response []string
	req      bool
}

func (r *Request) dump() {
	fmt.Printf("Request %s %s\n", r.Meta.ID, r.Meta.Timestamp.Format(time.RFC3339))
	fmt.Println("Payload")
	for i := range r.Payload {
		fmt.Println(r.Payload[i])
	}
	fmt.Println("Response")
	for i := range r.Response {
		fmt.Println(r.Response[i])
	}
}

func ReadRequestsLog() {
	rMap, err := readLog("requests.log.combined")
	utils.Must(err)
	fmt.Printf("len(rMap) %d\n", len(rMap))
	//utils.Must(printFiltered(rMap))
	utils.Must(replayInsertWErr(rMap))
	//utils.Must(replayTranscodeWErr(rMap))
	//utils.Must(replayConvertWErr(rMap))
}

func parseMeta(line string) (*Meta, error) {
	meta := strings.Split(line, " ")
	if len(meta) < 3 {
		return nil, errors.Errorf("Expected > 3 fields got %d", len(meta))
	}

	t, err := strconv.Atoi(meta[0])
	if err != nil {
		return nil, errors.Errorf("type should be int got: %s", meta[0])
	}

	ts, err := strconv.ParseInt(meta[2], 10, 64)
	if err != nil {
		return nil, errors.Errorf("timestamp should be int64 got: %s", meta[2])
	}

	return &Meta{
		Type:      t,
		ID:        meta[1],
		Timestamp: time.Unix(0, ts),
	}, nil
}

// scanning here is rather dumb:
// blank lines are ignored
// lines starting with "1" marks a beginning of a new request
// lines starting with "2" are marks a beginning of a response
// other lines are content of either a request or response
func readLog(path string) (map[string]*Request, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, errors.Wrap(err, "os.Open")
	}
	defer file.Close()

	rMap := make(map[string]*Request)

	var current *Request
	var ln int64
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		ln += 1

		if strings.TrimSpace(line) == "" {
			continue
		}

		if strings.HasPrefix(line, "1") {
			meta, err := parseMeta(line)
			if err != nil {
				log.Fatalf("Malformed request [Line %d]: %s", ln, err.Error())
			}

			r := Request{
				Meta:     meta,
				Payload:  make([]string, 0),
				Response: make([]string, 0),
				req:      true,
			}
			rMap[meta.ID] = &r
			current = &r
		} else if strings.HasPrefix(line, "2") {
			meta, err := parseMeta(line)
			if err != nil {
				log.Fatalf("Malformed response [Line %d]: %s", ln, err.Error())
			}

			if val, ok := rMap[meta.ID]; ok {
				val.req = false
			} else {
				log.Errorf("No Request for Response [Line %d]: %s", ln, meta.ID)
			}
		} else {
			if current == nil {
				continue
			}

			if current.req {
				current.Payload = append(current.Payload, line)
			} else {
				current.Response = append(current.Response, line)
			}
		}
	}

	return rMap, nil
}

func filterRequests(rMap map[string]*Request, filter func(*Request) bool) []*Request {
	requests := make([]*Request, 0)
	for _, v := range rMap {
		if filter(v) {
			requests = append(requests, v)
		}
	}

	sort.Slice(requests, func(i, j int) bool {
		return requests[i].Meta.Timestamp.Before(requests[j].Meta.Timestamp)
	})

	return requests
}

type FilterFunc func(r *Request) bool

func andFilters(filters ...FilterFunc) FilterFunc {
	return func(r *Request) bool {
		for i := range filters {
			if !filters[i](r) {
				return false
			}
		}
		return true
	}
}

func orFilters(filters ...FilterFunc) FilterFunc {
	return func(r *Request) bool {
		for i := range filters {
			if filters[i](r) {
				return true
			}
		}
		return false
	}
}

func beforeFilter(day time.Time) FilterFunc {
	return func(r *Request) bool {
		return r.Meta.Timestamp.Before(day)
	}
}

func afterFilter(day time.Time) FilterFunc {
	return func(r *Request) bool {
		return r.Meta.Timestamp.After(day)
	}
}

func payloadHasPrefixFilter(s string) FilterFunc {
	return func(r *Request) bool {
		return strings.HasPrefix(r.Payload[0], s)
	}
}

func responseHasSuffixFilter(s string) FilterFunc {
	return func(r *Request) bool {
		return strings.HasSuffix(r.Response[0], s)
	}
}

func responseExcludeSuffixFilter(s string) FilterFunc {
	return func(r *Request) bool {
		return !strings.HasSuffix(r.Response[0], s)
	}
}

func replayConvertWErr(rMap map[string]*Request) error {
	convertWErr := filterRequests(rMap, func(r *Request) bool {
		return strings.HasPrefix(r.Payload[0], "POST /operations/convert") &&
			strings.HasSuffix(r.Response[0], "Internal Server Error")
	})
	fmt.Printf("convertWErr %d\n", len(convertWErr))

	client := &http.Client{}
	for i := range convertWErr {
		r := convertWErr[i]
		strPayload := r.Payload[len(r.Payload)-1]
		var body api.ConvertRequest
		err := json.Unmarshal([]byte(strPayload), &body)
		if err != nil {
			return errors.Wrapf(err, "json.Unmarshal %s", r.Meta.ID)
		}

		req, err := http.NewRequest("POST",
			"http://app.mdb.bbdomain.org/operations/convert",
			strings.NewReader(strPayload))
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		if err != nil {
			return errors.Wrapf(err, "http client.Do %s", r.Meta.ID)
		}

		if resp.StatusCode != http.StatusOK {
			dumpHttpResponse(resp)
		}
	}

	return nil
}

func replayTranscodeWErr(rMap map[string]*Request) error {
	filter := andFilters(
		payloadHasPrefixFilter("POST /operations/transcode"),
		responseHasSuffixFilter("400 Bad Request"),
	)
	wErr := filterRequests(rMap, filter)
	fmt.Printf("wErr %d\n", len(wErr))

	client := &http.Client{}
	for i := range wErr {
		r := wErr[i]
		strPayload := r.Payload[len(r.Payload)-1]
		var body api.TranscodeRequest
		err := json.Unmarshal([]byte(strPayload), &body)
		if err != nil {
			fmt.Errorf("json.Unmarshal %s : %s", r.Meta.ID, err.Error())
		}
		if body.User == "" {
			body.User = "operator@dev.com"
		}
		if body.Station == "" {
			body.Station = "files.kabbalahmedia.info"
		}

		bodyPayload, err := json.Marshal(body)
		if err != nil {
			return errors.Wrapf(err, "json.Marshal %s", r.Meta.ID)
		}

		req, err := http.NewRequest("POST",
			"http://app.mdb.bbdomain.org/operations/transcode",
			bytes.NewReader(bodyPayload))
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		if err != nil {
			return errors.Wrapf(err, "http client.Do %s", r.Meta.ID)
		}

		if resp.StatusCode != http.StatusOK {
			dumpHttpResponse(resp)
		}
	}

	return nil
}

func replayInsertWErr(rMap map[string]*Request) error {
	log.Info("Setting up connection to MDB")
	mdb, err := sql.Open("postgres", viper.GetString("mdb.url"))
	utils.Must(err)
	utils.Must(mdb.Ping())
	defer mdb.Close()

	filter := andFilters(
		payloadHasPrefixFilter("POST /operations/insert"),
		responseHasSuffixFilter("400 Bad Request"),
	)
	wErr := filterRequests(rMap, filter)
	fmt.Printf("wErr %d\n", len(wErr))

	client := &http.Client{}
	for i := range wErr {
		r := wErr[i]
		strPayload := r.Payload[len(r.Payload)-1]

		strPayload = strings.Replace(strPayload, "\"publisher_uid\":\"null\"", "\"publisher_uid\":\"\"", 1)

		var body api.InsertRequest
		err := json.Unmarshal([]byte(strPayload), &body)
		if err != nil {
			fmt.Errorf("json.Unmarshal %s : %s", r.Meta.ID, err.Error())
		}

		_, _, err = api.FindFileBySHA1(mdb, body.Sha1)
		if err != nil {
			if _, ok := err.(api.FileNotFound); ok {
				body.Mode = "new"
			} else {
				log.Fatalf("Lookup file by sha1 %s: %s", body.Sha1, err.Error())
			}
		} else {
			body.Mode = "rename"
		}

		bodyPayload, err := json.Marshal(body)
		if err != nil {
			return errors.Wrapf(err, "json.Marshal %s", r.Meta.ID)
		}

		req, err := http.NewRequest("POST",
			"http://app.mdb.bbdomain.org/operations/insert",
			bytes.NewReader(bodyPayload))
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		if err != nil {
			return errors.Wrapf(err, "http client.Do %s", r.Meta.ID)
		}

		if resp.StatusCode != http.StatusOK {
			dumpHttpResponse(resp)
			fmt.Println(strPayload)
		}
	}

	return nil
}

func printFiltered(rMap map[string]*Request) error {
	//filter := andFilters(
	//	afterFilter(time.Date(2017, 5, 31, 0, 0, 0, 0, time.UTC)),
	//	beforeFilter(time.Date(2017, 6, 2, 0, 0, 0, 0, time.UTC)),
	//	responseExcludeSuffixFilter("200 OK"),
	//)
	filter := andFilters(
		//afterFilter(time.Date(2017, 9, 12, 0, 0, 0, 0, time.UTC)),
		payloadHasPrefixFilter("POST /operations/insert"),
		responseHasSuffixFilter("400 Bad Request"),
	)

	failed := filterRequests(rMap, filter)
	fmt.Printf("filtered %d\n", len(failed))

	for i := range failed {
		failed[i].dump()
	}

	return nil
}

func dumpHttpResponse(resp *http.Response) {
	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))
	resp.Body.Close()
}
