package models

import (
	"time"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
)

type (
	User struct {
		ID        uint64 `json:"omitempty",gorm:"primary_key"`
		Email     string `json:"omitempty",gorm:"type:varchar(64);unique_index"`
		Name      string `json:"omitempty",gorm:"type:char(32)"`
		Phone     string `json:"omitempty",gorm:"type:varchar(32)"`
		Comments  string `json:"omitempty",gorm:"type:varchar(255)"`
		CreatedAt time.Time `json:"omitempty"`
		UpdatedAt time.Time `json:"omitempty"`
		DeletedAt *time.Time `json:"omitempty",gorm:"index"`
	}

	OperationType struct {
		ID          uint64 `json:"omitempty",gorm:"primary_key"`
        Name        string `json:"omitempty"`
        Description string `json:"omitempty"`
	}

	Operation struct {
		ID        uint64        `json:"omitempty",gorm:"primary_key"`
        UID       string        `json:"omitempty"`
		Type      OperationType `json:"omitempty",gorm:"ForeignKey:TypeID"`
		TypeID    uint64        `json:"omitempty"`
		CreatedAt time.Time     `json:"omitempty"`
		Station   string        `json:"omitempty",gorm:"type:varchar(255)"`
		Details   string        `json:"omitempty",gorm:"type:varchar(255)"`
		User      User          `json:"omitempty"`
		UserID    uint64        `json:"omitempty"`
	}

	TranslatedContent struct {
		NameID        uint64            `json:"omitempty"`
		DescriptionID uint64            `json:"omitempty"`
		Name          StringTranslation `json:"omitempty",gorm:"ForeignKey:NameID;AssociationForeignKey:ID"`
		Description   StringTranslation `json:"omitempty",gorm:"ForeignKey:DescriptionID;AssociationForeignKey:ID"`
	}

	Collection struct {
		ID           uint64         `json:"omitempty"`
		UID          string         `json:"omitempty"`
		TypeID       uint64         `json:"omitempty"`
        Type         ContentType    `json:"omitempty",gorm:"ForeignKey:TypeID"`
		CreatedAt    time.Time      `json:"omitempty"`
		Properties   JsonB          `json:"omitempty"`
		ExternalID   string         `json:"omitempty"`
		ContentUnits []ContentUnit  `json:"omitempty",gorm:"many2many:collections_content_units;AssociationForeignKey:ID;ForeignKey:ID;"`
		TranslatedContent
	}

    ContentType struct {
        ID          uint64  `json:"omitempty"`
        Name        string  `json:"omitempty"`
        Description string  `json:"omitempty"`
    }

	ContentUnit struct {
		ID          uint64      `json:"omitempty"`
		UID         string      `json:"omitempty"`
        TypeID      uint64      `json:"omitempty"`
        Type        ContentType `json:"omitempty",gorm:"ForeignKey:TypeID;"`
		TranslatedContent       `json:"omitempty"`
		CreatedAt   time.Time   `json:"omitempty"`
		Properties  JsonB       `json:"omitempty"`
		Files       []File
		Collections []Collection `gorm:"many2many:collections_content_units;AssociationForeignKey:ID;ForeignKey:ID;"`
	}

    CollectionsContentUnit struct {
        CollectionID    uint64      `json:"omitempty"`
        Collection      Collection  `json:"omitempty",gorm:"ForeignKey:CollectionID;"`
        ContentUnitID   uint64      `json:"omitempty"`
        ContentUnit     ContentUnit `json:"omitempty",gorm:"ForeignKey:ContentUnitID;"`
        Name            string      `json:"omitempty"`
    }

	File struct {
		ID              uint64          `json:"omitempty"`
		UID             string          `json:"omitempty"`
		Name            string          `json:"omitempty"`
		Size            uint64          `json:"omitempty"`
		Type            string          `json:"omitempty"`
		SubType         string          `json:"omitempty"`
		MimeType        string          `json:"omitempty"`
		Sha1            sql.NullString  `json:"omitempty"`
		OperationID     uint64          `json:"omitempty"`
        Operation       Operation       `json:"omitempty",gorm:"ForeignKey:OperationID"`
		ContentUnitID   uint64          `json:"omitempty"`
		ContentUnit     ContentUnit     `json:"omitempty",gorm:"ForeignKey:ContentUnitID"`
        ParentID        sql.NullInt64   `json:"omitempty"`
		CreatedAt       time.Time       `json:"omitempty"`
		Language        string          `json:"omitempty"`
		BackupCount     uint            `json:"omitempty"`
		FirstBackupTime time.Time       `json:"omitempty"`
		Properties      JsonB           `json:"omitempty"`
	}

	StringTranslation struct {
		ID               uint64     `json:"omitempty"`
		Language         string     `json:"omitempty"`
		Text             string     `json:"omitempty"`
		OriginalLanguage string     `json:"omitempty"`
		User             User       `json:"omitempty"`
		UserID           uint64     `json:"omitempty"`
		CreatedAt        time.Time  `json:"omitempty"`
	}
)

type JsonB map[string]interface{}

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