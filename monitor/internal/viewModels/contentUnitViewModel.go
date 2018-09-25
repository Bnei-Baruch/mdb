package viewModels

import (
	"time"

	"gopkg.in/volatiletech/null.v6"
)

type ContentUnitViewModel struct {
	ID         int64     `boil:"id" json:"id" toml:"id" yaml:"id"`
	UID        string    `boil:"uid" json:"uid" toml:"uid" yaml:"uid"`
	TypeID     int64     `boil:"type_id" json:"type_id" toml:"type_id" yaml:"type_id"`
	TypeName   string    `boil:"type_name" json:"type_name" toml:"type_name" yaml:"type_name"`
	CreatedAt  time.Time `boil:"created_at" json:"created_at" toml:"created_at" yaml:"created_at"`
	Properties null.JSON `boil:"properties" json:"properties,omitempty" toml:"properties" yaml:"properties,omitempty"`
	Secure     int16     `boil:"secure" json:"secure" toml:"secure" yaml:"secure"`
	Published  bool      `boil:"published" json:"published" toml:"published" yaml:"published"`
}
