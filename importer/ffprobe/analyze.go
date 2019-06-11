package ffprobe

import (
	"encoding/hex"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/Bnei-Baruch/mdb/models"
	log "github.com/Sirupsen/logrus"
	"github.com/agnivade/levenshtein"
	"github.com/pkg/errors"
	"github.com/volatiletech/sqlboiler/queries/qm"

	"github.com/Bnei-Baruch/mdb/common"
	"github.com/Bnei-Baruch/mdb/utils"
)

const (
	BASE_PATH  = "importer/ffprobe/data/"
	AUDIO_FILE = BASE_PATH + "audio.csv"
	VIDEO_FILE = BASE_PATH + "video.csv"
)

type FFPData struct {
	Duration    float64
	BitRate     int
	AspectRatio string
	Resolution  string
	VideoSize   string
}

func (d *FFPData) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"duration":     d.Duration,
		"bit_rate":     d.BitRate,
		"aspect_ratio": d.AspectRatio,
		"resolution":   d.Resolution,
		"video_size":   d.VideoSize,
	}
}

type AugmentedFile struct {
	*models.File
	*FFPData
}

var LANG_RE *regexp.Regexp
var OT_RE = regexp.MustCompile("(?i)^[ot]$")
var BITRATE_RE = regexp.MustCompile("(?i)^(24k|96k|128k|hd)$")

func init() {
	i := 0
	keys := make([]string, len(common.LANG_MAP)-1)
	for k := range common.LANG_MAP {
		if k == "" {
			continue
		}
		keys[i] = k
		i++
	}
	LANG_RE = regexp.MustCompile(fmt.Sprintf("(?i)^(%s)$", strings.Join(keys, "|")))
}

func (af *AugmentedFile) NormalizedName() string {
	name := af.Name

	if idx := strings.LastIndex(name, "."); idx > 0 {
		name = name[:idx]
	}

	parts := strings.Split(name, "_")
	for i := range parts {
		part := strings.TrimSpace(parts[i])
		if OT_RE.MatchString(part) {
			parts[i] = "ot"
		} else if BITRATE_RE.MatchString(part) {
			parts[i] = "br"
		} else if LANG_RE.MatchString(part) {
			parts[i] = "lang"
		}
	}

	name = strings.Join(parts, "_")
	return strings.Replace(name, "_br", "", -1)
}

func Analyze() {
	clock, _ := Init()

	utils.Must(doAnalyze())

	Shutdown()
	log.Info("Success")
	log.Infof("Total run time: %s", time.Now().Sub(clock).String())
}

func doAnalyze() error {
	data, err := loadData()
	if err != nil {
		return errors.Wrap(err, "load csv data")
	}
	log.Infof("%d total files in csv", len(data))

	err = compareWithMDB(data)
	if err != nil {
		return errors.Wrap(err, "compare with mdb")
	}

	return nil
}

func compareWithMDB(ffpData map[string]*FFPData) error {
	fCount, err := models.Files(mdb).Count()
	if err != nil {
		return errors.Wrapf(err, "Load files count")
	}
	log.Infof("MDB has %d files", fCount)

	pageSize := 5000
	page := 0
	//fMap := make(map[string]*models.File, fCount)
	cuFiles := make(map[int64][]*AugmentedFile, 50000)
	for page*pageSize < int(fCount) {
		log.Infof("Loading page #%d", page)
		s := page * pageSize

		files, err := models.Files(mdb,
			qm.Offset(s),
			qm.Limit(pageSize)).
			All()
		if err != nil {
			return errors.Wrapf(err, "Load files page %d", page)
		}
		for i := range files {
			f := files[i]
			sha1 := hex.EncodeToString(f.Sha1.Bytes)
			ffp, ok := ffpData[sha1]
			if !ok {
				continue
			}

			//props := make(map[string]interface{})
			//if f.Properties.Valid {
			//	err := json.Unmarshal(f.Properties.JSON, &props)
			//	if err != nil {
			//		return errors.Wrapf(err, "json.Unmarshal file properties [%d]", f.ID)
			//	}
			//}
			//
			//if fDuration, ok := props["duration"]; ok && fDuration.(float64) > 0 {
			//	if math.Abs(ffp.Duration - fDuration.(float64)) > 3 {
			//		log.Infof("duration_epsilon [%d] mdb: %f\tffprobe: %f\tdiff: %f",
			//			f.ID, fDuration.(float64), ffp.Duration, math.Abs(ffp.Duration - fDuration.(float64)))
			//	}
			//} else{
			//	log.Infof("duration_missing [%d] %t\t%f\t%f", f.ID, ok, fDuration, ffp.Duration)
			//}

			if !f.ContentUnitID.Valid {
				continue
			}

			v, ok := cuFiles[f.ContentUnitID.Int64]
			if !ok {
				v = make([]*AugmentedFile, 0)
			}
			v = append(v, &AugmentedFile{File: f, FFPData: ffp})
			cuFiles[f.ContentUnitID.Int64] = v

			//fMap[sha1] = files[i]
		}
		page++
	}

	clipRE := regexp.MustCompile("(?i)(clip|promo|5min)")

	for cuID, augFiles := range cuFiles {
		if len(augFiles) < 2 {
			continue
		}

		filesByPattern := make(map[string][]*AugmentedFile)
		for i := range augFiles {
			k := augFiles[i].NormalizedName()
			v, ok := filesByPattern[k]
			if !ok {
				v = make([]*AugmentedFile, 0)
			}
			filesByPattern[k] = append(v, augFiles[i])
		}

		if len(filesByPattern) < 2 {
			continue
		}

		//log.Infof("CU [%d] has %d patterns", cuID, len(filesByPattern))
		wClip := make([]string, 0)
		woClip := make([]string, 0)
		for k := range filesByPattern {
			//log.Infof("%d\t%s", len(v), k)
			if clipRE.MatchString(k) {
				wClip = append(wClip, k)
			} else {
				woClip = append(woClip, k)
			}
		}

		if len(wClip) > 0 && len(woClip) > 0 {
			log.Infof("CU [%d] has mixed patterns: %d clips %d wo clips", cuID, len(wClip), len(woClip))
			for i := range wClip {
				log.Infof("%s", wClip[i])
			}
			for i := range woClip {
				log.Infof("%s", woClip[i])
			}
		}

		//durations := make([]float64, len(augFiles))
		//for i := range augFiles {
		//	durations[i] = augFiles[i].Duration
		//}
		//
		//breaks := jenks.NaturalBreaks(durations, len(durations))
		////log.Infof("Breaks [%d]: %v", cuID, breaks)
		//j := 0
		//ddup := breaks[:1]
		//for i := 1; i < len(breaks); i++ {
		//	if breaks[i]-ddup[j] < 3 {
		//		continue
		//	}
		//	ddup = append(ddup, breaks[i])
		//	j++
		//}
		//
		//if len(ddup) < 2 {
		//	continue
		//}
		//log.Infof("DDuped Breaks [%d]: %v", cuID, ddup)
		//
		//augFilesByDuration := make(map[float64][]*AugmentedFile)
		//groupedFileCount := 0
		//for i := range ddup {
		//	duration := ddup[i]
		//	augFilesByDuration[duration] = make([]*AugmentedFile,0)
		//	for j := range augFiles {
		//		if math.Abs(duration - augFiles[j].Duration) < 3 {
		//			augFilesByDuration[duration] = append(augFilesByDuration[duration], augFiles[j])
		//			groupedFileCount++
		//		}
		//	}
		//}
		//if len(ddup) != len(augFilesByDuration) {
		//	log.Warnf("durations breakpoints and file groups differ: %d != %d", len(ddup), len(augFilesByDuration))
		//}
		//if groupedFileCount != len(augFiles) {
		//	log.Warnf("grouped files count differ: %d != %d", len(augFiles), groupedFileCount)
		//}
		//
		//for k, v := range augFilesByDuration {
		//	fileNames := make(map[string]bool)
		//	for i:=0; i< len(v); i++ {
		//		normalizedName := v[i].Name
		//		if idx := strings.LastIndex(normalizedName, "."); idx>0 {
		//			normalizedName = normalizedName[:idx]
		//		}
		//		fileNames[normalizedName] = true
		//	}
		//	keys := make([]string, len(fileNames))
		//	i := 0
		//	for k := range fileNames {
		//		keys[i] = k
		//		i++
		//	}
		//	pDist := pairwiseDistance(keys)
		//	log.Infof("pDist [%d] %f: %v", len(pDist), k, pDist)
		//}
	}

	return nil
}

func pairwiseDistance(data []string) [][]int {
	n := len(data)
	dist := make([][]int, n)
	for i := 0; i < n; i++ {
		dist[i] = make([]int, n)
		for j := 0; j < n; j++ {
			if j == i {
				continue
			}
			dist[i][j] = levenshtein.ComputeDistance(data[i], data[j])
		}
	}

	return dist
}

func loadData() (map[string]*FFPData, error) {
	records, err := utils.ReadCSV(AUDIO_FILE)
	if err != nil {
		return nil, errors.Wrap(err, "Read audio file")
	}

	data := make(map[string]*FFPData, 500000)
	for i := range records {
		r := records[i]

		duration, err := strconv.ParseFloat(r[2], 64)
		if err != nil {
			log.Errorf("Bad Duration, audio [%d]: %s", i, r[2])
		}
		bitrate, err := strconv.Atoi(r[1])
		if err != nil {
			log.Errorf("Bad BitRate, audio [%d]: %s", i, r[1])
		}

		data[r[0]] = &FFPData{
			Duration: duration,
			BitRate:  bitrate,
		}
	}

	records, err = utils.ReadCSV(VIDEO_FILE)
	if err != nil {
		return nil, errors.Wrap(err, "Read video file")
	}

	for i := range records {
		r := records[i]

		duration, err := strconv.ParseFloat(r[4], 64)
		if err != nil {
			log.Errorf("Bad Duration, video [%d]: %s", i, r[4])
		}

		aspectRatio := r[1]
		if aspectRatio == "" {
			//log.Warnf("No AspectRatio, video [%d]: %s", i, aspectRatio)
		}

		resolution := r[2]
		if resolution == "" {
			//log.Warnf("No Resolution, video [%d]: %s", i, resolution)
		}

		videoSize := r[3]
		if videoSize == "" {
			//log.Warnf("No VideoSize, video [%d]: %s", i, r)
		}

		data[r[0]] = &FFPData{
			Duration:    duration,
			AspectRatio: aspectRatio,
			Resolution:  resolution,
			VideoSize:   videoSize,
		}
	}

	return data, nil
}
