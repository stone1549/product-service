package common

import (
	"github.com/shopspring/decimal"
)

type Product struct {
	Id               string           `json:"id"`
	Name             string           `json:"name"`
	DisplayImage     *string          `json:"displayImage"`
	Thumbnail        *string          `json:"thumbnail"`
	Price            *decimal.Decimal `json:"price"`
	Description      *string          `json:"description"`
	ShortDescription *string          `json:"shortDescription"`
	Quantity         int              `json:"quantity"`
}
