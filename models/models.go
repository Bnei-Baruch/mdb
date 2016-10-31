package models

import (
	"time"
	"database/sql/driver"
	"encoding/json"
)

type JSONB map[string]interface{}

func (j JSONB) Value() (driver.Value, error) {
	valueString, err := json.Marshal(j)
	return string(valueString), err
}

func (j *JSONB) Scan(value interface{}) error {
	if err := json.Unmarshal(value.([]byte), &j); err != nil {
		return err
	}
	return nil
}

type User struct {
	ID        uint64 `gorm:"primary_key"`
	Email     string `gorm:"type:varchar(64);unique_index"`
	Name      string `gorm:"type:char(32)"`
	Phone     string `gorm:"type:varchar(32)"`
	Comments  string `gorm:"type:varchar(255)"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `gorm:"index"`
}

type IString struct {
	ID               uint64 `gorm:"primary_key"`
	Language         string `gorm:"primary_key;type:char(2)"`
	Text             string
	OriginalLanguage string `gorm:"type:char(2)"`
	User             User
	UserID           uint64
	CreatedAt        time.Time
}

func (IString) TableName() string {
	return "strings"
}

type OperationType struct {
	ID          uint64 `gorm:"primary_key"`
	Name        string `gorm:"type:varchar(32);unique_index"`
	Description string `gorm:"type:varchar(255)"`
}

type Operation struct {
	ID        uint64 `gorm:"primary_key"`
	Type      OperationType
	TypeID    uint64
	CreatedAt time.Time
	Station   string `gorm:"type:varchar(255)"`
	Details   string `gorm:"type:varchar(255)"`
	User      User
	UserID    uint64
}

type ContentType struct {
	ID          uint64 `gorm:"primary_key"`
	Name        string `gorm:"type:varchar(32);unique_index"`
	Description string `gorm:"type:varchar(255)"`
}

type Collection struct {
	ID            uint64 `gorm:"primary_key"`
	UID           string `gorm:"type:char(10);unique_index"`
	Type          ContentType
	TypeID        uint64
	Name          IString
	NameID        uint64 `gorm:"column:name"`
	Description   IString
	DescriptionID uint64 `gorm:"column:description"`
	CreatedAt     time.Time
	Properties    JSONB `gorm:"type:jsonb"`
	ExternalID    string `gorm:"type:varchar(255);unique_index"`
}

type ContentUnit struct {
	ID            uint64 `gorm:"primary_key"`
	UID           string `gorm:"type:char(10);unique_index"`
	Type          ContentType
	TypeID        uint64
	Name          IString
	NameID        uint64 `gorm:"column:name"`
	Description   IString
	DescriptionID uint64 `gorm:"column:description"`
	CreatedAt     time.Time
	Properties    JSONB `gorm:"type:jsonb"`
}

type File struct {
	ID              uint64 `gorm:"primary_key"`
	UID             string `gorm:"type:char(10);unique_index"`
	Name            string `gorm:"type:varchar(255)"`
	Size            uint64
	Type            string `gorm:"type:varchar(16)"`
	Subtype         string `gorm:"type:varchar(16)"`
	MimeType        string `gorm:"type:varchar(255)"`
	Sha1            []byte `gorm:"type:bytea;column:SHA_1"`
	//Operation Operation
	ContentUnit     ContentUnit
	ContentUnitID   uint64
	Parent          *File
	ParentID        uint64
	CreatedAt       time.Time
	Language        string `gorm:"type:char(2)"`
	BackupCount     int8 `gorm:"type:smallint"`
	FirstBackupTime time.Time
	Properties      JSONB `gorm:"type:jsonb"`
}

type Person struct {
	ID            uint64 `gorm:"primary_key"`
	UID           string `gorm:"type:char(10);unique_index"`
	Name          IString
	NameID        uint64 `gorm:"column:name"`
	Description   IString
	DescriptionID uint64 `gorm:"column:description"`
}

type ContentRoles struct {
	ID            uint64 `gorm:"primary_key"`
	Name          IString
	NameID        uint64 `gorm:"column:name"`
	Description   IString
	DescriptionID uint64 `gorm:"column:description"`
}

type Tag struct {
	ID          uint64 `gorm:"primary_key"`
	Label       IString
	LabelID     uint64 `gorm:"column:label"`
	Description string `gorm:"type:varchar(255)"`
	Parent      *Tag
	ParentID    uint64
	CreatedAt   time.Time
}
