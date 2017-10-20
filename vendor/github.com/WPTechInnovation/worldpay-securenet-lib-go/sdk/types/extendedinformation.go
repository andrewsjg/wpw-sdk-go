package types

// ExtendedInformation is further metadata about a transaction
type ExtendedInformation struct {
}

// TypeOfGoods means the type of item being purchased
type TypeOfGoods string

const (

	// Digital digital goods
	Digital TypeOfGoods = "DIGITAL"
	// Physical physical goods
	Physical TypeOfGoods = "PHYSICAL"
)
