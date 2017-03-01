package api

type (
	// Common

	FileKey struct {
		FileName  string `json:"file_name" binding:"max=255"`
		Sha1      string `json:"sha1" binding:"omitempty,len=40,hexadecimal"`
		CreatedAt int64  `json:"created_at"`
		UID       string `json:"uid" binding:"omitempty,len=8,base64"`
	}

	Operation struct {
		Station    string `json:"station" binding:"required"`
		User       string `json:"user" binding:"required,email"`
		WorkflowID string `json:"workflow_id"`
	}

	FileUpdate struct {
		FileKey
		Size int64 `json:"size" binding:"required"`
	}

	// Operations

	CaptureStartRequest struct {
		Operation
		FileKey
		CollectionUID string `json:"collection_uid" binding:"omitempty,base64"`
		CaptureSource string `json:"capture_source"`
	}

	CaptureStopRequest struct {
		Operation
		FileKey
		CollectionUID string `json:"collection_uid" binding:"omitempty,base64"`
		CaptureSource string `json:"capture_source"`
		Sha1          string `json:"sha1" binding:"required,len=40,hexadecimal"`
		Size          int64 `json:"size" binding:"required"`
		ContentType   string `json:"content_type"`
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
		Duration int64  `json:"duration"`
		Existing FileKey `binding:"structonly"`
	}

	// simple CRUD
	CreateCollectionRequest struct {
		Type        string `json:"type" binding:"required"`
		UID         string `json:"uid" binding:"len=8"`
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
		Language    string `json:"language" binding:"len=2"`
	}
)
