package batch

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"

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
	Meta     Meta
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
	rMap, err := readLog("requests.log")
	utils.Must(err)
	//utils.Must(printFailed(rMap))
	utils.Must(replayTranscodeWErr(rMap))
	//utils.Must(replayConvertWErr(rMap))
}

func parseMeta(line string) Meta {
	meta := strings.Split(line, " ")
	id, _ := strconv.Atoi(meta[0])
	ts, _ := strconv.ParseInt(meta[2], 10, 64)
	return Meta{
		Type:      id,
		ID:        meta[1],
		Timestamp: time.Unix(0, ts),
	}
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
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		if strings.TrimSpace(line) == "" {
			continue
		}

		if strings.HasPrefix(line, "1") {
			meta := parseMeta(line)
			r := Request{
				Meta:     meta,
				Payload:  make([]string, 0),
				Response: make([]string, 0),
				req:      true,
			}
			rMap[meta.ID] = &r
			current = &r
		} else if strings.HasPrefix(line, "2") {
			meta := parseMeta(line)
			rMap[meta.ID].req = false
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
			fmt.Println("response Status:", resp.Status)
			fmt.Println("response Headers:", resp.Header)
			body, _ := ioutil.ReadAll(resp.Body)
			fmt.Println("response Body:", string(body))
			resp.Body.Close()
		}
	}

	return nil
}

func replayTranscodeWErr(rMap map[string]*Request) error {
	filter := andFilters(
		payloadHasPrefixFilter("POST /operations/transcode"),
		responseHasSuffixFilter("500 Internal Server Error"),
	)
	wErr := filterRequests(rMap, filter)
	fmt.Printf("wErr %d\n", len(wErr))

	client := &http.Client{}
	for i := range wErr {
		r := wErr[i]
		strPayload := r.Payload[len(r.Payload)-1]
		var body api.TranscodeRequestSuccess
		err := json.Unmarshal([]byte(strPayload), &body)
		if err != nil {
			return errors.Wrapf(err, "json.Unmarshal %s", r.Meta.ID)
		}

		req, err := http.NewRequest("POST",
			"http://app.mdb.bbdomain.org/operations/transcode",
			strings.NewReader(strPayload))
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		if err != nil {
			return errors.Wrapf(err, "http client.Do %s", r.Meta.ID)
		}

		if resp.StatusCode != http.StatusOK {
			fmt.Println("response Status:", resp.Status)
			fmt.Println("response Headers:", resp.Header)
			body, _ := ioutil.ReadAll(resp.Body)
			fmt.Println("response Body:", string(body))
			resp.Body.Close()
		}
	}

	return nil
}

func printFailed(rMap map[string]*Request) error {
	filter := andFilters(
		afterFilter(time.Date(2017, 9, 12, 0, 0, 0, 0, time.UTC)),
		payloadHasPrefixFilter("POST /operations/transcode"),
		responseHasSuffixFilter("500 Internal Server Error"),
	)

	failed := filterRequests(rMap, filter)
	fmt.Printf("failed %d\n", len(failed))

	for i := range failed {
		failed[i].dump()
	}

	return nil
}
