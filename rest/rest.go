package rest

type (
    Operation struct {
        Station     string  `json:"station" binding:"required"`
        User        string  `json:"user" binding:"required"`
	    Type        string  `json:"type" binding:"required"`
    }

    CaptureStart struct {
        Operation
        FileName    string  `json:"file_name" binding:"required,max=255"`
        CaptureID   string  `json:"capture_id" binding:"required,max=255"`
    }

    CaptureStop struct {
        CaptureStart
        Sha1        string  `json:"sha1" binding:"required,max=40"`
        Size        uint64  `json:"size" binding:"required"`
    }

    FileUpdate struct {
        FileName    string  `json:"file_name" binding:"required,max=255"`
        Sha1        string  `json:"sha1" binding:"required,max=40"`
        Size        uint64  `json:"size" binding:"required"`
    }

    Demux struct {
        Operation
        Sha1        string  `json:"sha1" binding:"required,max=40"`
        Original    FileUpdate
        Proxy       FileUpdate
    }

    Send struct {
        Operation
        Sha1        string  `json:"sha1" binding:"required,max=40"`
        Dest        FileUpdate
    }
)
