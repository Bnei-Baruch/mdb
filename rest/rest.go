package rest

type CaptureStart struct {
	Type      string  `json:"type"`
	Station   string  `json:"station"`
	User      string  `json:"user"`
	FileName  string  `json:"file_name" binding:"required,max=25"`
	CaptureID string  `json:"capture_id" binding:"required,max=255"`
}

