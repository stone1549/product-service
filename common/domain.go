package common

import (
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"time"
)

// Product holds information on an item for sale.
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

// OrderByKey represents a particular field that products can be sorted by.
type OrderByKey string

const (
	// OrderByCreated order from oldest to newest.
	OrderByCreated OrderByKey = "created"
	// OrderByCreatedDesc order from newest to oldest.
	OrderByCreatedDesc OrderByKey = "createdDesc"
	// OrderByUpdated order from least recently updated to most recently updated.
	OrderByUpdated OrderByKey = "updated"
	// OrderByUpdatedDesc order from most recently updated to least recently updated.
	OrderByUpdatedDesc OrderByKey = "updatedDesc"
	// OrderByName order alphabetically by name.
	OrderByName OrderByKey = "name"
	// OrderByNameDesc order reverse alphabetically by name.
	OrderByNameDesc OrderByKey = "nameDesc"
	// OrderByPrice order from least expensive to most expensive.
	OrderByPrice OrderByKey = "price"
	// OrderByPriceDesc order from most expensive to least expensive.
	OrderByPriceDesc OrderByKey = "priceDesc"
)

// Supported returns true if sorting by the given key is currently implemented.
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

// GetMirrorKey returns the opposite or mirror of the receiver key, for instance alphabetically by name ascending and
// alphabetically by name descending.
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

// OrderBy represents of and list of keys to sort by.
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

// Add will add the given key to the underlying list of keys or returns an error if doing so is not valid.
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

// Order retrieves an ordered slice of keys to sort by.
func (ob *OrderBy) Order() []OrderByKey {
	if len(ob.keys) == 0 {
		ob.Add(OrderByUpdatedDesc)
		ob.Add(OrderByCreatedDesc)
	}

	result := make([]OrderByKey, len(ob.keys))
	copy(result, ob.keys)
	return result
}
