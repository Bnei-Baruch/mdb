package batch

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	log "github.com/Sirupsen/logrus"

	"github.com/Bnei-Baruch/mdb/utils"
)

type FFPstreamTags struct {
	Language string `json:"language"`
}

type FFPstream struct {
	Index          int           `json:"index"`
	CodecName      string        `json:"codec_name"`
	CodecLongName  string        `json:"codec_long_name"`
	CodecType      string        `json:"codec_type"`
	CodecTimeBase  string        `json:"codec_time_base"`
	CodecTagString string        `json:"codec_tag_string"`
	CodecTag       string        `json:"codec_tag"`
	SampleFmt      string        `json:"sample_fmt"`
	SampleRate     string        `json:"sample_rate"`
	Channels       int           `json:"channels"`
	BitsPerSample  int           `json:"bits_per_sample"`
	RFrameRate     string        `json:"r_frame_rate"`
	AvgFrameRate   string        `json:"avg_frame_rate"`
	TimeBase       string        `json:"time_base"`
	StartPts       int           `json:"start_pts"`
	StartTime      string        `json:"start_time"`
	DurationTs     int           `json:"duration_ts"`
	Duration       string        `json:"duration"`
	BitRate        string        `json:"bit_rate"`
	Tags           FFPstreamTags `json:"tags"`

	Profile            string `json:"profile"`
	Width              int    `json:"width"`
	Height             int    `json:"height"`
	CodedWidth         int    `json:"coded_width"`
	CodedHeight        int    `json:"coded_height"`
	HasBFrames         int    `json:"has_b_frames"`
	SampleAspectRatio  string `json:"sample_aspect_ratio"`
	DisplayAspectRatio string `json:"display_aspect_ratio"`
	PixFmt             string `json:"pix_fmt"`
	Level              int    `json:"level"`
	ChromaLocation     string `json:"chroma_location"`
	Refs               int    `json:"refs"`
}

type FFPformat struct {
	Filename       string                 `json:"filename"`
	NbStreams      int                    `json:"nb_streams"`
	NbPrograms     int                    `json:"nb_programs"`
	FormatName     string                 `json:"format_name"`
	FormatLongName string                 `json:"format_long_name"`
	StartTime      string                 `json:"start_time"`
	Duration       string                 `json:"duration"`
	Size           string                 `json:"size"`
	BitRate        string                 `json:"bit_rate"`
	ProbeScore     int                    `json:"probe_score"`
	Tags           map[string]interface{} `json:"tags"`
}

type FFprobeMetadata struct {
	Streams []FFPstream `json:"streams"`
	Format  FFPformat   `json:"format"`
}

func ImportFFprobeMetadata() {
	baseDir := "/home/edos/projects/wmv_meta"

	fileList := make([]string, 0)
	err := filepath.Walk(baseDir, func(path string, f os.FileInfo, err error) error {
		if strings.HasSuffix(path, "json") {
			fileList = append(fileList, path)
		}
		return err
	})
	utils.Must(err)

	badJson := 0
	noVideo := 0
	log.Infof("%d json files", len(fileList))
	for _, file := range fileList {
		data, err := readFFProbeJson(file)
		if err != nil {
			badJson++
			continue
		}

		var video *FFPstream
		for _, s := range data.Streams {
			if s.CodecType == "video" {
				video = &s
			}
		}

		if video != nil {
			log.Infof("%s %dx%d %s %s", video.DisplayAspectRatio, video.Width, video.Height, video.RFrameRate, data.Format.Filename)
		} else {
			noVideo++
		}

	}
	log.Infof("%d bad json files", badJson)
	log.Infof("%d no video files", noVideo)
}

func readFFProbeJson(filename string) (FFprobeMetadata, error) {
	raw, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	var c FFprobeMetadata
	err = json.Unmarshal(raw, &c)
	return c, err
}
