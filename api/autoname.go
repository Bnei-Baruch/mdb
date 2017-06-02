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
	DescribeCollection(boil.Executor, *models.Collection, CITMetadata) ([]*models.CollectionI18n, error)
}

type GenericDescriber struct{}

func (d GenericDescriber) DescribeContentUnit(exec boil.Executor,
	cu *models.ContentUnit,
	metadata CITMetadata) ([]*models.ContentUnitI18n, error) {

	return nil, nil
}

func (d GenericDescriber) DescribeCollection(exec boil.Executor,
	cu *models.Collection,
	metadata CITMetadata) ([]*models.CollectionI18n, error) {
	return nil, nil
}

type LessonPartDescriber struct{}

func (d LessonPartDescriber) DescribeContentUnit(exec boil.Executor,
	cu *models.ContentUnit,
	metadata CITMetadata) ([]*models.ContentUnitI18n, error) {

	if metadata.Part.Valid && metadata.Part.Int == 0 {
		i18ns, err := GetI18ns("lesson_preparation")
		if err != nil {
			return nil, errors.Wrapf(err, "Get I18ns")
		}

		cui81ns := make([]*models.ContentUnitI18n, 0)
		for k, v := range i18ns {
			i18n := &models.ContentUnitI18n{
				ContentUnitID: cu.ID,
				Language:      k,
				Name:          null.StringFrom(v),
			}
			cui81ns = append(cui81ns, i18n)
		}
		return cui81ns, nil
	}

	return nil, nil
}

var CUDescribers = map[string]ContentUnitDescriber{
	CT_LESSON_PART: LessonPartDescriber{},
}

var CDescribers = map[string]CollectionDescriber{
	CT_DAILY_LESSON: GenericDescriber{},
}

func AutonameContentUnit(exec boil.Executor, cu *models.ContentUnit, metadata CITMetadata) error {
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

func AutonameCollection(exec boil.Executor, c *models.Collection, metadata CITMetadata) error {
	describer, ok := CDescribers[CONTENT_TYPE_REGISTRY.ByID[c.TypeID].Name]
	if !ok {
		describer = GenericDescriber{}
	}

	i18ns, err := describer.DescribeCollection(exec, c, metadata)
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
