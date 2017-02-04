package rest

type (
	// Helping structures.
	FileKey struct {
		FileName  string `json:"file_name" binding:"max=255"`
		Sha1      string `json:"sha1" binding:"max=40"`
		CreatedAt int64  `json:"created_at"`
	}

	Operation struct {
		Station string `json:"station" binding:"required"`
		User    string `json:"user" binding:"required"`
	}

	FileUpdate struct {
		FileKey
		Size uint64 `json:"size" binding:"required"`
	}

	// Operation structures.
	CaptureStart struct {
		Operation
		FileName      string `json:"file_name" binding:"required,max=255"`
		CreatedAt     int64  `json:"created_at"`
		CaptureID     string `json:"capture_id" binding:"required,max=255"`
		CaptureSource string `json:"capture_source" binding:"required"`
	}

	CaptureStop struct {
		CaptureStart
		Sha1 string `json:"sha1" binding:"required,max=40"`
		Size uint64 `json:"size" binding:"required"`
		Part string `json:"part" binding:"required"`
	}

	Demux struct {
		Operation
		FileKey
		Original FileUpdate
		Proxy    FileUpdate
	}

	Send struct {
		Operation
		FileKey
		Dest FileUpdate
	}

	Upload struct {
		Operation
		FileUpdate
		Url      string `json:"url" binding:"required"`
		Duration uint64 `json:"duration"`
        Existing FileKey `binding:"structonly"`
	}
)
