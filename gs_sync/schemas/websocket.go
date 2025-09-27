package schemas

type TaskRequest struct {
	PreExcluded []uint `json:"pre_excluded"`
	Excluded    []uint `json:"excluded"`
}

type UserDrop struct {
	RowId        uint    `json:"row_id"`
	FirstName    string  `json:"first_name"`
	LastName     string  `json:"last_name"`
	FractionFrom uint    `json:"fraction_from"`
	FractionTo   uint    `json:"fraction_to"`
	Fraction     uint    `json:"fraction"`
	MaxFraction  uint    `json:"max_fraction"`
}