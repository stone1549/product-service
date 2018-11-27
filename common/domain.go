package common

import (
	"github.com/pkg/errors"
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

type OrderByKey string

const (
	OrderByCreated     OrderByKey = "created"
	OrderByCreatedDesc OrderByKey = "createdDesc"
	OrderByUpdated     OrderByKey = "updated"
	OrderByUpdatedDesc OrderByKey = "updatedDesc"
	OrderByName        OrderByKey = "name"
	OrderByNameDesc    OrderByKey = "nameDesc"
	OrderByPrice       OrderByKey = "price"
	OrderByPriceDesc   OrderByKey = "priceDesc"
)

func (obk OrderByKey) Supported() bool {
	switch obk {
	case OrderByCreated:
		fallthrough
	case OrderByCreatedDesc:
		fallthrough
	case OrderByUpdated:
		fallthrough
	case OrderByUpdatedDesc:
		fallthrough
	case OrderByName:
		fallthrough
	case OrderByNameDesc:
		fallthrough
	case OrderByPrice:
		fallthrough
	case OrderByPriceDesc:
		return true
	default:
		return false
	}
}

func (obk OrderByKey) GetMirrorKey() OrderByKey {
	switch obk {
	case OrderByCreated:
		return OrderByCreatedDesc
	case OrderByCreatedDesc:
		return OrderByCreated
	case OrderByUpdated:
		return OrderByCreatedDesc
	case OrderByUpdatedDesc:
		return OrderByUpdated
	case OrderByName:
		return OrderByNameDesc
	case OrderByNameDesc:
		return OrderByName
	case OrderByPrice:
		return OrderByPriceDesc
	case OrderByPriceDesc:
		return OrderByPrice
	default:
		return ""
	}
}

type OrderBy struct {
	keys []OrderByKey
}

func orderByContains(keys []OrderByKey, key OrderByKey) bool {
	for _, k := range keys {
		if key == k {
			return true
		}
	}
	return false
}

func (ob *OrderBy) Add(key OrderByKey) error {
	if !key.Supported() {
		return errors.Errorf("Attempted to add unsupported key %s", key)
	}

	if orderByContains(ob.keys, key) {
		return errors.Errorf("Attempted to add duplicate key %s", key)
	}

	if key.GetMirrorKey() != "" && orderByContains(ob.keys, key.GetMirrorKey()) {
		return errors.Errorf("Attempted to add conflicting key %s", key)
	}

	ob.keys = append(ob.keys, key)
	return nil
}

func (ob *OrderBy) Order() []OrderByKey {
	if len(ob.keys) == 0 {
		ob.Add(OrderByUpdatedDesc)
		ob.Add(OrderByCreatedDesc)
	}

	result := make([]OrderByKey, len(ob.keys))
	copy(result, ob.keys)
	return result
}
