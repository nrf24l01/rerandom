package task

type SheetRow struct {
	Id            uint   `json:"id"`
	Fraction      uint   `json:"fraction"`
	FractionFrom  uint   `json:"fraction_from"`
	FractionTo    uint   `json:"fraction_to"`
	LastName      string `json:"last_name"`
	FirstName     string `json:"first_name"`
	Alive         bool   `json:"alive"`
}

type Action struct {
	Type  uint // 1 - change fraction, 2 - change alive
	RowId uint // row id
	Param uint // new fraction or alive value
}

type Sheet struct {
	ClearUsers    []SheetRow
	Actions       []Action
	ModifiedUsers []SheetRow
}