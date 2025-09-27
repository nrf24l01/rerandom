package schemas

type TaskRequest struct {
	Type  uint `json:"type"`  // 1 - change fraction, 2 - change alive
	RowId uint `json:"row_id"`
	Param uint `json:"param"` // new fraction or alive value
}