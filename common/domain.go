package common

import (
	"github.com/shopspring/decimal"
	"time"
)

type Product struct {
	Id               string           `json:"id"`
	Name             string           `json:"name"`
	DisplayImage     *string          `json:"displayImage"`
	Thumbnail        *string          `json:"thumbnail"`
	Price            *decimal.Decimal `json:"price"`
	Description      *string          `json:"description"`
	ShortDescription *string          `json:"shortDescription"`
	QtyInStock       int              `json:"qtyInStock"`
	CreatedAt        *time.Time       `json:"createdAt"`
	UpdatedAt        *time.Time       `json:"updatedAt"`
}
