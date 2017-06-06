package api

import (
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
	"github.com/vattle/sqlboiler/boil"

	"github.com/Bnei-Baruch/mdb/bindata"
	"github.com/Bnei-Baruch/mdb/models"
	"gopkg.in/nullbio/null.v6"
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
	return i18ns, nil
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

	i18ns := make([]*models.ContentUnitI18n, 0)
	for k, v := range names {
		i18n := &models.ContentUnitI18n{
			ContentUnitID: cu.ID,
			Language:      k,
			Name:          null.StringFrom(v),
		}
		i18ns = append(i18ns, i18n)
	}

	return i18ns, nil
}

func (d GenericDescriber) DescribeCollection(c *models.Collection) ([]*models.CollectionI18n, error) {

	i18nKey := fmt.Sprintf("content_type.%s", CONTENT_TYPE_REGISTRY.ByID[c.TypeID].Name)
	names, err := GetI18ns(i18nKey)
	if err != nil {
		return nil, errors.Wrapf(err, "Get I18ns")
	}

	i18ns := make([]*models.CollectionI18n, 0)
	if len(names) == 0 {
		return i18ns, nil
	}

	for k, v := range names {
		i18n := &models.CollectionI18n{
			CollectionID: c.ID,
			Language:     k,
			Name:         null.StringFrom(v),
		}
		i18ns = append(i18ns, i18n)
	}

	return i18ns, nil
}

type LessonPartDescriber struct{}

func (d LessonPartDescriber) DescribeContentUnit(exec boil.Executor,
	cu *models.ContentUnit,
	metadata CITMetadata) ([]*models.ContentUnitI18n, error) {

	var err error
	var names map[string]string

	if metadata.Part.Valid && metadata.Part.Int == 0 {
		names, err = GetI18ns("lesson_preparation")
		if err != nil {
			return nil, errors.Wrapf(err, "Get I18ns")
		}
	}

	if len(names) == 0 {
		return new(GenericDescriber).DescribeContentUnit(exec, cu, metadata)
	}

	i18ns := make([]*models.ContentUnitI18n, 0)
	for k, v := range names {
		i18n := &models.ContentUnitI18n{
			ContentUnitID: cu.ID,
			Language:      k,
			Name:          null.StringFrom(v),
		}
		i18ns = append(i18ns, i18n)
	}

	return i18ns, nil
}

var CUDescribers = map[string]ContentUnitDescriber{
	CT_LESSON_PART: LessonPartDescriber{},
}

var CDescribers = map[string]CollectionDescriber{}

func DescribeContentUnit(exec boil.Executor, cu *models.ContentUnit, metadata CITMetadata) error {
	describer, ok := CUDescribers[CONTENT_TYPE_REGISTRY.ByID[cu.TypeID].Name]
	if !ok {
		describer = GenericDescriber{}
	}

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
