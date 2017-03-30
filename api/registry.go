package api

import "github.com/Bnei-Baruch/mdb/models"

var (
	CONTENT_TYPE_REGISTRY   = &ContentTypeRegistry{}
	OPERATION_TYPE_REGISTRY = &OperationTypeRegistry{}
	SOURCE_TYPE_REGISTRY = &SourceTypeRegistry{}

	// kmedia - select concat('"',code3,'": "',locale,'",') from languages;
	LANG_MAP = map[string]string {
		"ENG": "en",
		"HEB": "he",
		"RUS": "ru",
		"SPA": "es",
		"ITA": "it",
		"GER": "de",
		"DUT": "nl",
		"FRE": "fr",
		"POR": "pt",
		"TRK": "tr",
		"POL": "pl",
		"ARB": "ar",
		"HUN": "hu",
		"FIN": "fi",
		"LIT": "lt",
		"JPN": "ja",
		"BUL": "bg",
		"GEO": "ka",
		"NOR": "no",
		"SWE": "sv",
		"HRV": "hr",
		"CHN": "zh",
		"FAR": "fa",
		"RON": "ro",
		"HIN": "hi",
		"MKD": "mk",
		"SLV": "sl",
		"LAV": "lv",
		"SLK": "sk",
		"CZE": "cs",
		"UKR": "ua",
	}
)

type ContentTypeRegistry struct {
	ByName map[string]*models.ContentType
	ByID map[int64]*models.ContentType  // DO NOT REMOVE: Used by ETL's in archive site
}

func (r *ContentTypeRegistry) Init() error {
	types, err := models.ContentTypesG().All()
	if err != nil {
		return err
	}

	r.ByName = make(map[string]*models.ContentType)
	r.ByID = make(map[int64]*models.ContentType)
	for _, t := range types {
		r.ByName[t.Name] = t
		r.ByID[t.ID] = t
	}

	return nil
}

type OperationTypeRegistry struct {
	ByName map[string]*models.OperationType
}

func (r *OperationTypeRegistry) Init() error {
	types, err := models.OperationTypesG().All()
	if err != nil {
		return err
	}

	r.ByName = make(map[string]*models.OperationType)
	for _, t := range types {
		r.ByName[t.Name] = t
	}

	return nil
}


type SourceTypeRegistry struct {
	ByName map[string]*models.SourceType
	ByID map[int64]*models.SourceType
}

func (r *SourceTypeRegistry) Init() error {
	types, err := models.SourceTypesG().All()
	if err != nil {
		return err
	}

	r.ByName = make(map[string]*models.SourceType)
	r.ByID = make(map[int64]*models.SourceType)
	for _, t := range types {
		r.ByName[t.Name] = t
		r.ByID[t.ID] = t
	}

	return nil
}