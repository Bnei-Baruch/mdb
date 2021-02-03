package common

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Bnei-Baruch/mdb/models"
	"github.com/Bnei-Baruch/mdb/utils"
	log "github.com/Sirupsen/logrus"
	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/queries/qm"
	"gopkg.in/volatiletech/null.v6"
	"strings"
	"time"
)

// Return standard language or LANG_UNKNOWN
//
// 	if len(lang) = 2 we assume it's an MDB language code and check KNOWN_LANGS.
// 	if len(lang) = 3 we assume it's a workflow / kmedia lang code and check LANG_MAP.
func StdLang(lang string) string {
	switch len(lang) {
	case 2:
		if l := strings.ToLower(lang); KNOWN_LANGS.MatchString(l) {
			return l
		}
	case 3:
		if l, ok := LANG_MAP[strings.ToUpper(lang)]; ok {
			return l
		}
	}

	return LANG_UNKNOWN
}

func CreateCUTypeSource(s *models.Source, mdb boil.Executor) (*models.ContentUnit, error) {
	if has := haveSourceCUTypeSource(s.UID, mdb); has {
		return nil, errors.New(fmt.Sprintf("Have CU type %s for this source %s", CT_SOURCE, s.UID))
	}

	var props map[string]interface{}

	p, err := json.Marshal(props)
	if err != nil {
		return nil, err
	}

	cu := &models.ContentUnit{
		UID:        s.UID,
		TypeID:     CONTENT_TYPE_REGISTRY.ByName[CT_SOURCE].ID,
		Secure:     SEC_PUBLIC,
		Published:  true,
		Properties: null.JSONFrom(p),
		CreatedAt:  time.Now(),
	}
	log.Debug("CreateCUTypeSource - before insert unit", cu)
	err = cu.Insert(mdb)
	log.Debug("CreateCUTypeSource - after insert unit", cu, err)
	utils.Must(err)

	err = cu.AddSources(mdb, false, s)
	utils.Must(err)
	return cu, nil
}

func haveSourceCUTypeSource(suid string, mdb boil.Executor) bool {
	tid := CONTENT_TYPE_REGISTRY.ByName[CT_SOURCE].ID
	has, err := models.ContentUnits(
		mdb,
		qm.Where("content_units.type_id = ? AND s.uid = ?", tid, suid),
		qm.InnerJoin("content_units_sources as cus ON cus.content_unit_id = content_units.id"),
		qm.InnerJoin("sources as s ON cus.source_id = s.id"),
	).Exists()
	utils.Must(err)
	return has
}
