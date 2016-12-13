package models

import (
	"time"
	"database/sql/driver"
	"encoding/json"

    "github.com/jinzhu/gorm"
)

type (

	User struct {
        gorm.Model

		ID        uint64 `gorm:"primary_key"`
		Email     string `gorm:"type:varchar(64);unique_index"`
		Name      string `gorm:"type:char(32)"`
		Phone     string `gorm:"type:varchar(32)"`
		Comments  string `gorm:"type:varchar(255)"`
		CreatedAt time.Time
		UpdatedAt time.Time
		DeletedAt *time.Time `gorm:"index"`
	}

	OperationType struct {
		ID uint64 `gorm:"primary_key"`
		TranslatedContent
	}

	Operation struct {
		ID        uint64 `gorm:"primary_key"`
		Type      OperationType
		TypeID    uint64
		CreatedAt time.Time
		Station   string `gorm:"type:varchar(255)"`
		Details   string `gorm:"type:varchar(255)"`
		User      User
		UserID    uint64
	}

	TranslatedContent struct {
		NameID        uint64
		DescriptionID uint64
		Name          StringTranslation `gorm:"ForeignKey:NameID;AssociationForeignKey:Text"`
		Description   StringTranslation `gorm:"ForeignKey:DescriptionID;AssociationForeignKey:Text"`
	}

	Collection struct {
		ID           uint64
		UID          string
		TypeID       uint64
		CreateAt     time.Time
		Properties   JsonB
		ExternalID   string
		ContentUnits []ContentUnit `gorm:"many2many:collections_content_units;"`
		TranslatedContent
		/*Name          string
		Lang          string
		AssetType     string
		Date          time.Time
		Size          uint
		Containers    []Container `gorm:"many2many:containers_file_assets;"`*/
	}

	ContentUnit struct {
		ID          uint64
		UID         string
		TypeID      uint64
		TranslatedContent
		CreatedAt   time.Time
		Properties  JsonB
		Files       []File
		Collections []Collection `gorm:"many2many:collections_content_units;"`
	}

	File struct {
		ID              uint64
		UID             string
		Name            string
		Size            uint64
		Type            string
		SubType         string
		MimeType        string
		Sha1            []byte
		OperationID     uint64
		ContentUnitID   uint64
		ContentUnit     ContentUnit
		ParentID        uint64
		CreateAt        time.Time
		Language        string
		BackupCount     uint
		FirstBackupTime time.Time
		Properties      JsonB
	}

	StringTranslation struct {
		ID               uint64
		Language         string
		Text             string
		OriginalLanguage string
		User             User
		UserID           uint64
		CreatedAt        time.Time
	}
)

type JsonB map[string]interface{}

func (StringTranslation) TableName() string {
	return "strings"
}

func (j JsonB) Value() (driver.Value, error) {
	valueString, err := json.Marshal(j)
	return string(valueString), err
}

func (j *JsonB) Scan(value interface{}) error {
	if err := json.Unmarshal(value.([]byte), &j); err != nil {
		return err
	}
	return nil
}

