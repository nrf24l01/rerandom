package schemas

type PredictAddRequest struct {
	Answ      int `json:"answ" validate:"required"`
	Min       *int `json:"min,omitempty"`
	Max       *int `json:"max,omitempty"`
	DropCount *int `json:"drop_count,omitempty"`
}

type PredictEditRequest struct {
	Answ      *int `json:"answ,omitempty"`
	Min       *int `json:"min,omitempty"`
	Max       *int `json:"max,omitempty"`
	DropCount *int `json:"drop_count,omitempty"`
}

type PredictResponse struct {
	UUID        string `json:"uuid"`
	Answ        int    `json:"answ"`
	Min         *int   `json:"min,omitempty"`
	Max         *int   `json:"max,omitempty"`
	Dropped     int    `json:"dropped"`
	MaxDrops    int    `json:"max_drops"`
	Finished    bool   `json:"finished"`
	Added       int64  `json:"added"`
	LastDropped *int64 `json:"last_dropped,omitempty"`
}
