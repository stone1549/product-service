package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/blevesearch/bleve"
	"github.com/shopspring/decimal"
	"github.com/stone1549/product-service/common"
	"io/ioutil"
	"sort"
	"time"
)

type inMemoryProductRepository struct {
	products []common.Product
	index    bleve.Index
}

type orderBySort struct {
	Products []common.Product
	Order    common.OrderBy
}

func (obs *orderBySort) Len() int {
	return len(obs.Products)
}

func (obs *orderBySort) Swap(i, j int) {
	obs.Products[i], obs.Products[j] = obs.Products[j], obs.Products[i]
}

func compareTimePtr(a, b *time.Time) int {
	if a == nil && b == nil {
		return 0
	} else if a == nil {
		return 1
	} else if b == nil {
		return -1
	}

	if a.Unix() == b.Unix() {
		return 0
	} else if a.Unix() < b.Unix() {
		return -1
	} else {
		return 1
	}
}

func compareStrPtr(a, b *string) int {
	if a == nil && b == nil {
		return 0
	} else if a == nil {
		return 1
	} else if b == nil {
		return -1
	}

	if *a == *b {
		return 0
	} else if *a < *b {
		return -1
	} else {
		return 1
	}
}

func compareDecimalPtr(a, b *decimal.Decimal) int {
	if a == nil && b == nil {
		return 0
	} else if a == nil {
		return 1
	} else if b == nil {
		return -1
	}

	if *a == *b {
		return 0
	} else if a.LessThan(*b) {
		return -1
	} else {
		return 1
	}
}

func (obs *orderBySort) Less(i, j int) bool {
	for _, key := range obs.Order.Order() {
		switch key {
		case common.OrderByCreated:
			value := compareTimePtr(obs.Products[i].CreatedAt, obs.Products[j].CreatedAt)
			if value == 0 {
				continue
			} else {
				return value < 0
			}
		case common.OrderByCreatedDesc:
			value := compareTimePtr(obs.Products[i].CreatedAt, obs.Products[j].CreatedAt) * -1
			if value == 0 {
				continue
			} else {
				return value < 0
			}
		case common.OrderByUpdated:
			value := compareTimePtr(obs.Products[i].UpdatedAt, obs.Products[j].UpdatedAt)
			if value == 0 {
				continue
			} else {
				return value < 0
			}
		case common.OrderByUpdatedDesc:
			value := compareTimePtr(obs.Products[i].UpdatedAt, obs.Products[j].UpdatedAt) * -1
			if value == 0 {
				continue
			} else {
				return value < 0
			}
		case common.OrderByName:
			value := compareStrPtr(&obs.Products[i].Name, &obs.Products[j].Name)
			if value == 0 {
				continue
			} else {
				return value < 0
			}
		case common.OrderByNameDesc:
			value := compareStrPtr(&obs.Products[i].Name, &obs.Products[j].Name) * -1
			if value == 0 {
				continue
			} else {
				return value < 0
			}
		case common.OrderByPrice:
			value := compareDecimalPtr(obs.Products[i].Price, obs.Products[j].Price)
			if value == 0 {
				continue
			} else {
				return value < 0
			}
		case common.OrderByPriceDesc:
			value := compareDecimalPtr(obs.Products[i].Price, obs.Products[j].Price) * -1
			if value == 0 {
				continue
			} else {
				return value > 0
			}
		default:
			continue
		}
	}

	return false
}

// GetProducts retrieves a list of the first X products starting from the given cursor.
func (impr *inMemoryProductRepository) GetProducts(_ context.Context, first int, cursor string,
	orderBy common.OrderBy) (ProductList, error) {
	products := make([]common.Product, 0)

	newCursor := cursor
	reachedCursor := false
	if cursor == "" {
		reachedCursor = true
	}

	sortedProducts := make([]common.Product, len(impr.products))
	copy(sortedProducts, impr.products)
	sortByOrderBy := orderBySort{sortedProducts, orderBy}
	sort.Sort(&sortByOrderBy)

	for _, product := range sortByOrderBy.Products {
		if len(products) == first {
			break
		} else if product.Id == cursor {
			reachedCursor = true
			continue
		} else if reachedCursor {
			productCopy := product
			products = append(products, productCopy)
			newCursor = productCopy.Id
		}
	}
	return ProductList{products, newCursor}, nil
}

func findProductById(products []common.Product, id string) (*common.Product, error) {
	for _, product := range products {
		if id == product.Id {
			return &product, nil
		}
	}

	return nil, nil
}

// GetProduct retrieves a product from the given id.
func (impr *inMemoryProductRepository) GetProduct(_ context.Context, id string) (*common.Product, error) {
	return findProductById(impr.products, id)
}

// SearchProducts retrieves the first X matches starting from the given cursor.
func (impr *inMemoryProductRepository) SearchProducts(ctx context.Context, searchTxt string, first int,
	cursor string) (ProductList, error) {
	query := bleve.NewMatchQuery(searchTxt)
	search := bleve.NewSearchRequest(query)
	searchResults, err := impr.index.Search(search)

	if err != nil {
		return ProductList{}, err
	}

	products := make([]common.Product, 0)
	newCursor := cursor
	reachedCursor := false
	if cursor == "" {
		reachedCursor = true
	}
	for _, hit := range searchResults.Hits {
		if len(products) == first {
			break
		} else if hit.ID == cursor {
			reachedCursor = true
			continue
		} else if reachedCursor {
			product, err := findProductById(impr.products, hit.ID)

			if err != nil {
				return ProductList{}, err
			}
			products = append(products, *product)
			newCursor = product.Id
		}
	}

	return ProductList{products, newCursor}, nil
}

// MakeInMemoryRepository constructs an in memory backed ProductRepository from the given configuration.
func MakeInMemoryRepository(config common.Configuration) (ProductRepository, error) {
	var products []common.Product
	var err error

	// open a new index
	mapping := bleve.NewIndexMapping()
	idx, err := bleve.NewMemOnly(mapping)

	if err == bleve.ErrorIndexPathExists {
		idx, err = bleve.Open("document")
	}

	if err != nil {
		return nil, err
	}

	if config.GetInitDataSet() == "" {
		products = make([]common.Product, 0)
	} else {
		products, err = loadInitInMemoryDataset(config.GetInitDataSet())
	}

	for _, product := range products {
		var shortDescription string
		if product.ShortDescription != nil {
			shortDescription = *product.ShortDescription
		}
		var description string
		if product.Description != nil {
			description = *product.Description
		}
		// index name, id, short description, and full description
		idxData := fmt.Sprintf("%s %s %s %s", product.Name, product.Id, shortDescription,
			description)
		err = idx.Index(product.Id, idxData)

		if err != nil {
			return nil, err
		}
	}

	return &inMemoryProductRepository{products, idx}, err
}

func loadInitInMemoryDataset(dataset string) ([]common.Product, error) {
	var err error
	products := make([]common.Product, 0)

	if err != nil {
		return products, err
	}

	jsonBytes, err := ioutil.ReadFile(dataset)
	if err != nil {
		return products, err
	}

	err = json.Unmarshal(jsonBytes, &products)

	return products, err
}
