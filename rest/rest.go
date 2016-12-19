package rest

type (
    CaptureStart struct {
	    Type        string  `json:"type" binding:"required"`
        Station     string  `json:"station" binding:"required"`
        User        string  `json:"user" binding:"required"`
        FileName    string  `json:"file_name" binding:"required,max=255"`
        CaptureID   string  `json:"capture_id" binding:"required,max=255"`
    }

    CaptureStop struct {
        CaptureStart
        Sha1        string  `json:"sha1" binding:"required,max=40"`
        Size        uint64  `json:"size" binding:"required"`
    }
)
