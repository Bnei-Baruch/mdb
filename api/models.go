package api

type (
	// Common

	FileKey struct {
		FileName  string `json:"file_name" binding:"max=255"`
		Sha1      string `json:"sha1" binding:"max=40"`
		CreatedAt int64  `json:"created_at"`
		UID       string `json:"uid"`
	}

	Operation struct {
		Station    string `json:"station" binding:"required"`
		User       string `json:"user" binding:"required"`
		WorkflowID string `json:"workflow_id"`
	}

	FileUpdate struct {
		FileKey
		Size uint64 `json:"size" binding:"required"`
	}

	// Operations

	CaptureStartRequest struct {
		Operation
		FileKey
		CollectionUID  string `json:"collection_uid" binding:"max=8"`
		CaptureSource string `json:"capture_source"`
	}

	CaptureStopRequest struct {
		Operation
		FileKey
		CollectionUID  string `json:"collection_uid" binding:"max=8"`
		CaptureSource string `json:"capture_source"`
		Sha1          string `json:"sha1" binding:"required,max=40"`
		Size          uint64 `json:"size" binding:"required"`
		Part          string `json:"part"`
	}

	DemuxRequest struct {
		Operation
		FileKey
		Original FileUpdate
		Proxy    FileUpdate
	}

	SendRequest struct {
		Operation
		FileKey
		Dest FileUpdate
	}

	UploadRequest struct {
		Operation
		FileUpdate
		Url      string  `json:"url" binding:"required"`
		Duration uint64  `json:"duration"`
		Existing FileKey `binding:"structonly"`
	}

	// simple CRUD
	CreateCollectionRequest struct {
		Type        string `json:"type" binding:"required"`
		UID         string `json:"uid" binding:"max=8"`
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
		Language    string `json:"language" binding:"max=2"`
	}
)
