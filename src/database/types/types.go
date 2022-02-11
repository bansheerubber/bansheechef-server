package types

//go:generate database-type IngredientType
type IngredientType struct {
	Id int
	Name string
	Barcode string
	MaxAmount float32
	IsVolume bool
	UnitCount int
	ImageId int
}
