package api

import (
	"fmt"
	"strconv"
	"time"
)

type (
	// Common

	Operation struct {
		Station    string `json:"station" binding:"required"`
		User       string `json:"user" binding:"required,email"`
		WorkflowID string `json:"workflow_id"`
	}

	File struct {
		FileName  string     `json:"file_name" binding:"required,max=255"`
		Sha1      string     `json:"sha1" binding:"required,len=40,hexadecimal"`
		Size      int64      `json:"size" binding:"required"`
		CreatedAt *Timestamp `json:"created_at" binding:"required"`
		Type      string     `json:"type" binding:"max=16"`
		SubType   string     `json:"sub_type" binding:"max=16"`
		MimeType  string     `json:"mime_type" binding:"max=255"`
		Language  string     `json:"language" binding:"omitempty,len=2"`
	}

	AVFile struct {
		File
		Duration float64 `json:"duration"`
	}

	Rename struct {
		Sha1     string     `json:"sha1" binding:"required,len=40,hexadecimal"`
		FileName string     `json:"file_name" binding:"required,max=255"`
	}

	// Operations

	CaptureStartRequest struct {
		Operation
		FileName      string `json:"file_name" binding:"max=255"`
		CaptureSource string `json:"capture_source"`
		CollectionUID string `json:"collection_uid"`
	}

	CaptureStopRequest struct {
		Operation
		AVFile
		CaptureSource string `json:"capture_source"`
		CollectionUID string `json:"collection_uid"`
		Part          string `json:"part"`
	}

	DemuxRequest struct {
		Operation
		Sha1          string `json:"sha1" binding:"required,len=40,hexadecimal"`
		Original      AVFile `json:"original"`
		Proxy         AVFile `json:"proxy"`
		CaptureSource string `json:"capture_source"`
	}

	TrimRequest struct {
		Operation
		OriginalSha1  string    `json:"original_sha1" binding:"required,len=40,hexadecimal"`
		ProxySha1     string    `json:"proxy_sha1" binding:"required,len=40,hexadecimal"`
		Original      AVFile    `json:"original"`
		Proxy         AVFile    `json:"proxy"`
		In            []float64 `json:"in"`
		Out           []float64 `json:"out"`
		CaptureSource string    `json:"capture_source"`
	}

	SendRequest struct {
		Operation
		Original     Rename    `json:"original"`
		Proxy        Rename    `json:"proxy"`
		WorkflowType string `json:"workflow_type"`
	}

	UploadRequest struct {
		Operation
		AVFile
		Url string `json:"url" binding:"required"`
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

// A time.Time like structure with Unix timestamp JSON marshalling
type Timestamp struct {
	time.Time
}

func (t *Timestamp) MarshalJSON() ([]byte, error) {
	ts := t.Time.Unix()
	stamp := fmt.Sprint(ts)

	return []byte(stamp), nil
}

func (t *Timestamp) UnmarshalJSON(b []byte) error {
	ts, err := strconv.Atoi(string(b))
	if err != nil {
		return err
	}

	t.Time = time.Unix(int64(ts), 0)

	return nil
}
