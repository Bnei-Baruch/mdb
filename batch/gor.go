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

	"github.com/Bnei-Baruch/mdb/api"
	"github.com/Bnei-Baruch/mdb/utils"
)

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

func ReadRequestsLog() {
	path := "requests.log"
	file, err := os.Open(path)
	utils.Must(err)
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
	fmt.Printf("rMap %d\n", len(rMap))

	convertWErr := make([]*Request, 0)
	for _, v := range rMap {
		if strings.HasPrefix(v.Payload[0], "POST /operations/convert") &&
			strings.HasSuffix(v.Response[0], "Internal Server Error") {
			convertWErr = append(convertWErr, v)
		}
	}

	fmt.Printf("convertWErr %d\n", len(convertWErr))
	sort.Slice(convertWErr, func(i, j int) bool {
		return convertWErr[i].Meta.Timestamp.Before(convertWErr[j].Meta.Timestamp)
	})

	client := &http.Client{}
	for i := range convertWErr {
		r := convertWErr[i]
		strPayload := r.Payload[len(r.Payload)-1]
		var body api.ConvertRequest
		err := json.Unmarshal([]byte(strPayload), &body)
		if err != nil {
			fmt.Printf("error parsing request body: %s\n", strPayload)
		}

		req, err := http.NewRequest("POST",
			"http://localhost:8080/operations/convert",
			strings.NewReader(strPayload))
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		if err != nil {
			panic(err)
		}

		if resp.StatusCode != http.StatusOK {
			fmt.Println("response Status:", resp.Status)
			fmt.Println("response Headers:", resp.Header)
			body, _ := ioutil.ReadAll(resp.Body)
			fmt.Println("response Body:", string(body))
			resp.Body.Close()
		}
	}

}

func parseMeta(line string) Meta {
	meta := strings.Split(line, " ")
	id, _ := strconv.Atoi(meta[0])
	ts, _ := strconv.ParseInt(meta[2], 10, 64)
	return Meta{
		Type:      id,
		ID:        meta[1],
		Timestamp: time.Unix(ts, 0),
	}
}
