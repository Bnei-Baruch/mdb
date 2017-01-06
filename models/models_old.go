package models

import (
	"time"
)

type (
	FileAsset struct {
		ID         int
		Name       string
		Lang       string
		AssetType  string
		Date       time.Time
		Size       int
		Containers []Container `gorm:"many2many:containers_file_assets;"`
	}

	Container struct {
		ID           int
		Name         string
		CreatedAt    time.Time
		UpdatedAt    time.Time
		Lang         string
		PlaytimeSecs int
		FileAssets   []FileAsset `gorm:"many2many:containers_file_assets;"`
		Descriptions []ContainerDescription
	}

	ContainerDescription struct {
		Id            int
		Container     Container
		ContainerID   int
		ContainerDesc string
		Lang          string
	}
)
