package mdb

import (
	"time"
)


type User struct {
	ID        uint64 `gorm:"primary_key"`
	Email     string `gorm:"type:varchar(64);unique_index"`
	Name     string `gorm:"type:char(32)"`
	Phone     string `gorm:"type:varchar(32)"`
	Comments     string `gorm:"type:varchar(255)"`
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
	UserID           uint64 `gorm:"index"`
	CreatedAt        time.Time
}

func (IString) TableName() string {
	return "strings"
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
