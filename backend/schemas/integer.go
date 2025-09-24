package schemas

type IntegerRequest struct {
	Num    int    `json:"num" query:"num" validate:"required,min=1,max=10000"`
	Min    int    `json:"min" query:"min" validate:"required"`
	Max    int    `json:"max" query:"max" validate:"required"`
	Col    int    `json:"col" query:"col" validate:"omitempty,min=1"`
	Base   int    `json:"base" query:"base" validate:"omitempty,oneof=2 8 10 16"`
	Format string `json:"format" query:"format" validate:"omitempty,oneof=plain html json"`
	Rnd    string `json:"rnd" query:"rnd" validate:"omitempty,oneof=new id.identifier"`
	Cl     string `json:"cl" query:"cl" validate:"omitempty,oneof=w b"`
}
