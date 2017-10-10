package api

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/pkg/errors"
	"github.com/vattle/sqlboiler/boil"
	"gopkg.in/nullbio/null.v6"

	"github.com/Bnei-Baruch/mdb/bindata"
	"github.com/Bnei-Baruch/mdb/models"
	"github.com/Bnei-Baruch/mdb/utils"
	"github.com/vattle/sqlboiler/queries/qm"
)

var I18n map[string]map[string]string

func init() {
	data, err := bindata.Asset("data/i18n.json")
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(data, &I18n)
	if err != nil {
		panic(err)
	}
}

type MissingI18n struct {
	Key      string
	Language string
}

func (e MissingI18n) Error() string {
	return fmt.Sprintf("Missing I18n %s[%s]", e.Key, e.Language)
}

func GetI18ns(key string) (map[string]string, error) {
	i18ns, ok := I18n[key]
	if !ok {
		return nil, errors.Errorf("Unknown i18n key: %s", key)
	}

	// we create a clone so users won't change our own instance
	var newMap = make(map[string]string,len(i18ns))
	for k,v := range i18ns {
		newMap[k] = v
	}

	return newMap, nil
}

func T(key, language string) (string, error) {
	i18ns, err := GetI18ns(key)
	if err != nil {
		return "", err
	}
	val, ok := i18ns[language]
	if ok {
		return val, nil
	}
	return "", MissingI18n{Key: key, Language: language}
}

type ContentUnitDescriber interface {
	DescribeContentUnit(boil.Executor, *models.ContentUnit, CITMetadata) ([]*models.ContentUnitI18n, error)
}

type CollectionDescriber interface {
	DescribeCollection(*models.Collection) ([]*models.CollectionI18n, error)
}

type GenericDescriber struct{}

func (d GenericDescriber) DescribeContentUnit(exec boil.Executor,
	cu *models.ContentUnit,
	metadata CITMetadata) ([]*models.ContentUnitI18n, error) {

	names := map[string]string{
		LANG_HEBREW:  metadata.FinalName,
		LANG_ENGLISH: metadata.FinalName,
		LANG_RUSSIAN: metadata.FinalName,
	}

	return makeCUI18ns(cu.ID, names), nil
}

func (d GenericDescriber) DescribeCollection(c *models.Collection) ([]*models.CollectionI18n, error) {
	i18nKey := fmt.Sprintf("content_type.%s", CONTENT_TYPE_REGISTRY.ByID[c.TypeID].Name)
	names, err := GetI18ns(i18nKey)
	if err != nil {
		return nil, errors.Wrapf(err, "Get I18ns")
	}

	if len(names) == 0 {
		return make([]*models.CollectionI18n, 0), nil
	}

	return makeCI18ns(c.ID, names), nil
}

type LessonPartDescriber struct{}

func (d LessonPartDescriber) DescribeContentUnit(exec boil.Executor,
	cu *models.ContentUnit,
	metadata CITMetadata) ([]*models.ContentUnitI18n, error) {

	var err error
	var names map[string]string

	// Use collection number if applicable
	var cNumber *int
	if metadata.CollectionUID.Valid && metadata.Number.Valid {
		cNumber = &metadata.Number.Int
	}

	if metadata.Part.Valid && metadata.Part.Int == 0 {
		names, err = GetI18ns("autoname.lesson_preparation")
		if err != nil {
			return nil, errors.Wrap(err, "Get I18ns")
		}

		if cNumber != nil {
			for k, v := range names {
				names[k] = fmt.Sprintf("%s %d", v, *cNumber)
			}
		}
	} else if metadata.Major != nil {
		idx := metadata.Major.Idx
		switch metadata.Major.Type {
		case "source":
			if idx >= len(metadata.Sources) {
				log.Warnf("metadata.major index out of bounds got %d but only %d elements in sources",
					idx, len(metadata.Sources))
			}
			names, err = nameBySourceUID(exec, metadata.Sources[idx])
			if err != nil {
				return nil, errors.Wrap(err, "Name by source")
			}
			break
		case "tag":
			if idx >= len(metadata.Tags) {
				log.Warnf("metadata.major index out of bounds got %d but only %d elements in tags",
					idx, len(metadata.Tags))
			}
			names, err = nameByTagUID(exec, metadata.Tags[idx], cNumber)
			if err != nil {
				return nil, errors.Wrap(err, "Name by tag")
			}
			break
		default:
			log.Warnf("Unknown metadata.major type %s", metadata.Major.Type)
		}
	} else {
		// no Major info from metadata
		// give names by what we have in DB
		// may be used by batch processes
		err = cu.L.LoadTags(exec, true, cu)
		if err != nil {
			return nil, errors.Wrap(err, "Load tags from DB")
		}
		if len(cu.R.Tags) > 0 {
			names, err = nameByTagUID(exec, cu.R.Tags[0].UID, nil)
			if err != nil {
				return nil, errors.Wrap(err, "Name by tag")
			}
		} else {
			err = cu.L.LoadSources(exec, true, cu)
			if err != nil {
				return nil, errors.Wrap(err, "Load sources from DB")
			}
			if len(cu.R.Sources) > 0 {
				names, err = nameBySourceUID(exec, cu.R.Sources[0].UID)
				if err != nil {
					return nil, errors.Wrap(err, "Name by source")
				}
			}
		}
	}

	// make sure major languages has something
	genericNames := map[string]string{
		LANG_HEBREW:  metadata.FinalName,
		LANG_ENGLISH: metadata.FinalName,
		LANG_RUSSIAN: metadata.FinalName,
	}
	names = mergeMaps(genericNames, names)

	return makeCUI18ns(cu.ID, names), nil
}

type MealDescriber struct{}

func (d MealDescriber) DescribeContentUnit(exec boil.Executor,
	cu *models.ContentUnit,
	metadata CITMetadata) ([]*models.ContentUnitI18n, error) {

	names, err := GetI18ns("autoname.meal")
	if err != nil {
		return nil, errors.Wrap(err, "Get I18ns")
	}

	return makeCUI18ns(cu.ID, names), nil
}

type VideoProgramChapterDescriber struct{}

func (d VideoProgramChapterDescriber) DescribeContentUnit(exec boil.Executor,
	cu *models.ContentUnit,
	metadata CITMetadata) ([]*models.ContentUnitI18n, error) {

	if !metadata.CollectionUID.Valid {
		return new(GenericDescriber).DescribeContentUnit(exec, cu, metadata)
	}

	collection, err := models.Collections(exec,
		qm.Where("uid=?", metadata.CollectionUID.String),
		qm.Load("CollectionI18ns")).One()
	if err != nil {
		return nil, errors.Wrap(err, "Lookup collection in DB")
	}

	var names map[string]string
	for _, language := range ALL_LANGS {
		i18n := getCollectionI18n(collection, language)
		if i18n == nil {
			continue
		}
		names[language] = fmt.Sprintf("%s %d", i18n.Name.String, metadata.Number.Int)
	}

	return makeCUI18ns(cu.ID, names), nil
}

var CUDescribers = map[string]ContentUnitDescriber{
	CT_LESSON_PART:           new(LessonPartDescriber),
	CT_VIDEO_PROGRAM_CHAPTER: new(VideoProgramChapterDescriber),
	CT_MEAL:                  new(MealDescriber),
}

var CDescribers = map[string]CollectionDescriber{}

type EventPartDescriber struct{}

func (d EventPartDescriber) DescribeContentUnit(exec boil.Executor,
	cu *models.ContentUnit,
	metadata CITMetadata) ([]*models.ContentUnitI18n, error) {

	collection, err := models.Collections(exec,
		qm.Where("uid=?", metadata.CollectionUID.String),
		qm.Load("CollectionI18ns")).One()
	if err != nil {
		return nil, errors.Wrap(err, "Lookup collection in DB")
	}

	var names map[string]string
	ct := CONTENT_TYPE_REGISTRY.ByID[cu.TypeID].Name
	switch ct {
	case CT_LESSON_PART, CT_FULL_LESSON:
		cuI18ns, err := new(LessonPartDescriber).DescribeContentUnit(exec, cu, metadata)
		if err != nil {
			return nil, err
		}

		names = make(map[string]string)
		for i := range cuI18ns {
			i18n := cuI18ns[i]
			names[i18n.Language] = i18n.Name.String
		}

	default:
		names, err := GetI18ns(fmt.Sprintf("content_type.%s", ct))
		if err != nil {
			return nil, errors.Wrap(err, "Get I18ns")
		}

		// add event part number as suffix
		if !metadata.Number.Valid {
			for k, v := range names {
				names[k] = fmt.Sprintf("%s %d", v, metadata.Number.Int)
			}
		}
	}

	// prefix names with event name
	var finalNames = make(map[string]string)
	for k, v := range names {
		i18n := getCollectionI18n(collection, k)
		if i18n == nil || !i18n.Name.Valid {
			continue
		}
		finalNames[k] = fmt.Sprintf("%s. %s", i18n.Name.String, v)
	}

	return makeCUI18ns(cu.ID, finalNames), nil
}

func DescribeContentUnit(exec boil.Executor, cu *models.ContentUnit, metadata CITMetadata) error {
	var describer ContentUnitDescriber

	// do we need a special describer based on collection type ?
	if metadata.CollectionUID.Valid {
		collection, err := models.Collections(exec,
			qm.Where("uid=?", metadata.CollectionUID.String)).One()
		if err != nil {
			return errors.Wrap(err, "Lookup collection in DB")
		}

		ct := CONTENT_TYPE_REGISTRY.ByID[collection.TypeID].Name
		if ct == CT_CONGRESS || ct == CT_HOLIDAY || ct == CT_UNITY_DAY || ct == CT_PICNIC {
			describer = new(EventPartDescriber)
		}
	}

	// no special describer based on collection type,
	// use describer based on CU type
	if describer == nil {
		var ok bool
		describer, ok = CUDescribers[CONTENT_TYPE_REGISTRY.ByID[cu.TypeID].Name]
		if !ok {
			describer = GenericDescriber{}
		}
	}

	// describe
	i18ns, err := describer.DescribeContentUnit(exec, cu, metadata)
	if err != nil {
		return errors.Wrap(err, "Auto naming content unit")
	}

	// TODO: reuse relevant repo method when possible
	err = cu.AddContentUnitI18ns(exec, true, i18ns...)
	if err != nil {
		return errors.Wrap(err, "Save to DB")
	}

	return nil
}

func DescribeCollection(exec boil.Executor, c *models.Collection) error {
	describer, ok := CDescribers[CONTENT_TYPE_REGISTRY.ByID[c.TypeID].Name]
	if !ok {
		describer = GenericDescriber{}
	}

	i18ns, err := describer.DescribeCollection(c)
	if err != nil {
		return errors.Wrap(err, "Auto naming collection")
	}

	// TODO: reuse relevant repo method when possible
	err = c.AddCollectionI18ns(exec, true, i18ns...)
	if err != nil {
		return errors.Wrap(err, "Save to DB")
	}

	return nil
}

func makeCI18ns(id int64, names map[string]string) []*models.CollectionI18n {
	i18ns := make([]*models.CollectionI18n, 0)
	for k, v := range names {
		i18n := &models.CollectionI18n{
			CollectionID: id,
			Language:     k,
			Name:         null.StringFrom(v),
		}
		i18ns = append(i18ns, i18n)
	}
	return i18ns
}

func makeCUI18ns(id int64, names map[string]string) []*models.ContentUnitI18n {
	i18ns := make([]*models.ContentUnitI18n, 0)
	for k, v := range names {
		i18n := &models.ContentUnitI18n{
			ContentUnitID: id,
			Language:      k,
			Name:          null.StringFrom(v),
		}
		i18ns = append(i18ns, i18n)
	}
	return i18ns
}

type sourceNamer interface {
	GetName(*models.Author, []*models.Source) (map[string]string, error)
}

func nameBySourceUID(exec boil.Executor, uid string) (map[string]string, error) {
	s, err := FindSourceByUID(exec, uid)
	if err != nil {
		return nil, errors.Wrapf(err, "Find source in DB")
	}

	path, err := FindSourcePath(exec, s.ID)
	if err != nil {
		return nil, errors.Wrapf(err, "Load source path from DB")
	}

	err = s.L.LoadSourceI18ns(exec, false, &path)
	if err != nil {
		return nil, errors.Wrapf(err, "Load sources i18ns from DB")
	}

	root := path[len(path)-1]
	author, err := FindAuthorBySourceID(exec, root.ID)
	if err != nil {
		return nil, errors.Wrapf(err, "Load author from DB")
	}

	// Find the first matching namer, bottom up, for this source
	var namer sourceNamer
	var ok bool
	for _, source := range path {
		if namer, ok = sourceNamers[source.UID]; ok {
			break
		}
	}
	if namer == nil {
		namer = new(PlainNamer)
	}

	// reverse path
	for left, right := 0, len(path)-1; left < right; left, right = left+1, right-1 {
		path[left], path[right] = path[right], path[left]
	}

	return namer.GetName(author, path)
}

type PlainNamer struct{}

// <author>, <path...>
func (n PlainNamer) GetName(author *models.Author, path []*models.Source) (map[string]string, error) {
	names := make(map[string]string)
	for _, language := range ALL_LANGS {
		vals := make([]string, 0)

		// author name
		ai18n := getAuthorI18n(author, language)
		if ai18n == nil || !ai18n.Name.Valid {
			continue
		}
		vals = append(vals, ai18n.Name.String)

		// sources path names
		for _, s := range path {
			i18n := getSourceI18n(s, language)
			if i18n == nil || !i18n.Name.Valid {
				break
			}
			vals = append(vals, i18n.Name.String)
		}

		// skip if we don't have all i18ns
		if len(vals) != 1+len(path) {
			continue
		}

		names[language] = strings.Join(vals, ". ")
	}

	return names, nil
}

type PrefaceNamer struct{}

// <author>. <leaf node in path>
func (n PrefaceNamer) GetName(author *models.Author, path []*models.Source) (map[string]string, error) {
	names := make(map[string]string)
	for _, language := range ALL_LANGS {
		vals := make([]string, 0)

		// author name
		ai18n := getAuthorI18n(author, language)
		if ai18n == nil || !ai18n.Name.Valid {
			continue
		}
		vals = append(vals, ai18n.Name.String)

		// leaf node name
		i18n := getSourceI18n(path[len(path)-1], language)
		if i18n == nil || !i18n.Name.Valid {
			continue
		}
		vals = append(vals, i18n.Name.String)

		names[language] = strings.Join(vals, ". ")
	}

	return names, nil
}

type LettersNamer struct {
	PrefaceNamer
}

// <author>. <leaf node in path - cleaned of (...) suffixes>
func (n LettersNamer) GetName(author *models.Author, path []*models.Source) (map[string]string, error) {
	names := make(map[string]string)

	baseNames, err := n.PrefaceNamer.GetName(author, path)
	if err != nil {
		return nil, err
	}

	for k, v := range baseNames {
		a := strings.Split(v, "(")
		names[k] = strings.TrimSpace(a[0])
	}

	return names, nil
}

type RBRecordsNamer struct{}

// <author>. Record <position>. <leaf node in path>
func (n RBRecordsNamer) GetName(author *models.Author, path []*models.Source) (map[string]string, error) {
	recordI18ns, err := GetI18ns("autoname.record")
	if err != nil {
		return nil, errors.Wrap(err, "Get I18ns")
	}

	names := make(map[string]string)
	for language, recordI18n := range recordI18ns {
		vals := make([]string, 0)

		// author name
		ai18n := getAuthorI18n(author, language)
		if ai18n == nil || !ai18n.Name.Valid {
			continue
		}
		vals = append(vals, ai18n.Name.String)

		// leaf node position & name
		leaf := path[len(path)-1]
		if !leaf.Position.Valid {
			continue
		}
		vals = append(vals, fmt.Sprintf("%s %d", recordI18n, leaf.Position.Int))

		i18n := getSourceI18n(leaf, language)
		if i18n == nil || !i18n.Name.Valid {
			continue
		}
		vals = append(vals, i18n.Name.String)

		names[language] = strings.Join(vals, ". ")
	}

	return names, nil
}

type RBArticlesNamer struct{}

// <author>. <leaf.name>. <number (year)>
func (n RBArticlesNamer) GetName(author *models.Author, path []*models.Source) (map[string]string, error) {
	names := make(map[string]string)

	// source.name is expected to be "<year-number> <non-i18n title>"
	// we strip the year-number part and transform it a little
	leaf := path[len(path)-1]
	yearNum := strings.Split(leaf.Name, " ")[0]
	yearNum = yearNum[1 : len(yearNum)-1]
	x := strings.Split(yearNum, "-")
	year, num := x[0], strings.TrimLeft(x[1], "0")
	if len(x) > 2 {
		num = num + "-" + x[2]
	}
	suffix := fmt.Sprintf("%s (%s)", num, year)

	for _, language := range ALL_LANGS {
		vals := make([]string, 0)

		// author name
		ai18n := getAuthorI18n(author, language)
		if ai18n == nil || !ai18n.Name.Valid {
			continue
		}
		vals = append(vals, ai18n.Name.String)

		// article name
		i18n := getSourceI18n(leaf, language)
		if i18n == nil || !i18n.Name.Valid {
			continue
		}
		vals = append(vals, i18n.Name.String)

		// num (year) suffix
		vals = append(vals, suffix)

		names[language] = strings.Join(vals, ". ")
	}

	return names, nil
}

type ShamatiNamer struct{}

// <author>. path[0].name, <path[1]'s number (heb gimatria)>. path[1].name
func (n ShamatiNamer) GetName(author *models.Author, path []*models.Source) (map[string]string, error) {
	names := make(map[string]string)

	// source.name is expected to be "<number> <non-i18n title>"
	// we strip the number part and transform it a little
	pLast := path[len(path)-1]
	num, err := strconv.Atoi(strings.Split(pLast.Name, " ")[0])
	if err != nil {
		return nil, errors.Wrapf(err, "Bad source.name %s", pLast.Name)
	}

	for _, language := range ALL_LANGS {
		vals := make([]string, 0)

		// author name
		ai18n := getAuthorI18n(author, language)
		if ai18n == nil || !ai18n.Name.Valid {
			continue
		}
		vals = append(vals, ai18n.Name.String)

		// collection name (Shamati)
		p0 := getSourceI18n(path[0], language)
		if p0 == nil || !p0.Name.Valid {
			continue
		}

		// article number cleaned and maybe hebrew letters numbering
		if language == LANG_HEBREW {
			vals = append(vals, fmt.Sprintf("%s, %s", p0.Name.String, utils.NumToHebrew(uint16(num))))
		} else {
			vals = append(vals, fmt.Sprintf("%s, %d", p0.Name.String, num))
		}

		// article name
		i18n := getSourceI18n(pLast, language)
		if i18n == nil || !i18n.Name.Valid {
			continue
		}
		vals = append(vals, i18n.Name.String)

		names[language] = strings.Join(vals, ". ")
	}

	return names, nil
}

type ZoharNamer struct{}

// <path[0]>, <path[1].description>, <path[2:]>
func (n ZoharNamer) GetName(author *models.Author, path []*models.Source) (map[string]string, error) {
	names := make(map[string]string)
	for _, language := range ALL_LANGS {
		vals := make([]string, 0)

		// sources path names
		for i, s := range path {
			i18n := getSourceI18n(s, language)
			if i18n == nil {
				break
			}

			var val string
			if i == 1 && i18n.Description.Valid {
				val = i18n.Description.String
			}
			if i != 1 && i18n.Name.Valid {
				val = i18n.Name.String
			}
			if val == "" {
				break
			}

			vals = append(vals, val)
		}

		// skip if we don't have all i18ns
		if len(vals) != len(path) {
			continue
		}

		names[language] = strings.Join(vals, ". ")
	}

	return names, nil
}

// Specs:
// https://docs.google.com/spreadsheets/d/1zY-MQlbZl9nIJA8MUaE0-LPWYU9qNRJk4svix0R5Gv0/edit?usp=sharing
var sourceNamers = map[string]sourceNamer{
	"L2jMWyce": new(PrefaceNamer), // BaalHaSulam Prefaces
	"SJDw9tHs": new(PrefaceNamer), // Rabash Prefaces
	"qMeV5M3Y": new(PrefaceNamer), // BaalHaSulam Articles
	"DVSS0xAR": new(LettersNamer), // BaalHaSulam Letters
	"b8SHlrfH": new(LettersNamer), // Rabash Letters
	//"xtKmrbb9": new(PlainNamer), // BaalHaSulam TES
	"qMUUn22b": new(ShamatiNamer),    // BaalHaSulam Shamati
	"rQ6sIUZK": new(RBArticlesNamer), // Rabash Articles
	"2GAdavz0": new(RBRecordsNamer),  // Rabash Records
	"QUBP2DYe": new(PrefaceNamer),    // Michael Laitman Articles
	"AwGBQX2L": new(ZoharNamer),      // Zohar La'am
}

var TAGS_WITH_CONST_NAMES = [...]string{"lfq8P1du", "iWPWpSNy"}

func nameByTagUID(exec boil.Executor, uid string, cNumber *int) (map[string]string, error) {

	// check constant names first
	for i := range TAGS_WITH_CONST_NAMES {
		if uid == TAGS_WITH_CONST_NAMES[i] {
			return GetI18ns(fmt.Sprintf("autoname.%s", uid))
		}
	}

	// Load tag details from DB
	t, err := FindTagByUID(exec, uid)
	if err != nil {
		return nil, errors.Wrapf(err, "Find tag in DB")
	}

	path, err := FindTagPath(exec, t.ID)
	if err != nil {
		return nil, errors.Wrapf(err, "Load tag path from DB")
	}

	err = t.L.LoadTagI18ns(exec, false, &path)
	if err != nil {
		return nil, errors.Wrapf(err, "Load tag i18ns from DB")
	}
	t = path[0]

	// create names
	names := make(map[string]string)

	// lesson topic || holiday has a different format
	if root := path[len(path)-1]; root.UID == "mS7hrYXK" || root.UID == "1nyptSIo" {
		prefixes, err := GetI18ns("autoname.lesson_by_topic_tag")
		if err != nil {
			return nil, errors.Wrap(err, "Get I18ns")
		}

		for k, v := range prefixes {
			// insert collection number if we have it
			var vv string
			if cNumber == nil {
				vv = fmt.Sprintf(v, "")
			} else {
				vv = fmt.Sprintf(v, strconv.Itoa(*cNumber))
			}

			i18n := getTagI18n(t, k)
			if i18n != nil && i18n.Label.Valid {
				names[k] = fmt.Sprintf("%s \"%s\"", vv, i18n.Label.String)
			}
		}
	} else {
		for _, language := range ALL_LANGS {
			i18n := getTagI18n(t, language)
			if i18n != nil && i18n.Label.Valid {
				names[language] = i18n.Label.String
			}
		}
	}

	return names, nil
}

func getAuthorI18n(author *models.Author, language string) *models.AuthorI18n {
	for _, i18n := range author.R.AuthorI18ns {
		if i18n.Language == language {
			return i18n
		}
	}
	return nil
}

func getSourceI18n(source *models.Source, language string) *models.SourceI18n {
	for _, i18n := range source.R.SourceI18ns {
		if i18n.Language == language {
			return i18n
		}
	}
	return nil
}

func getTagI18n(tag *models.Tag, language string) *models.TagI18n {
	for _, i18n := range tag.R.TagI18ns {
		if i18n.Language == language {
			return i18n
		}
	}
	return nil
}

func getCollectionI18n(collection *models.Collection, language string) *models.CollectionI18n {
	for _, i18n := range collection.R.CollectionI18ns {
		if i18n.Language == language {
			return i18n
		}
	}
	return nil
}

func mergeMaps(a, b map[string]string) map[string]string {
	c := make(map[string]string)
	for k, v := range a {
		c[k] = v
	}
	for k, v := range b {
		c[k] = v
	}
	return c
}
