package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/blevesearch/bleve"
	"github.com/stone1549/product-service/common"
	"io/ioutil"
	"sort"
)

type inMemoryProductRepository struct {
	products []common.Product
	index    bleve.Index
}

type sortDefault []common.Product

func (s sortDefault) Len() int {
	return len(s)
}

func (s sortDefault) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s sortDefault) Less(i, j int) bool {
	if s[i].UpdatedAt.Unix() > s[j].UpdatedAt.Unix() {
		return true
	} else if s[i].UpdatedAt.Unix() == s[j].UpdatedAt.Unix() {
		return s[i].CreatedAt.Unix() > s[j].CreatedAt.Unix()
	} else {
		return false
	}
}

func (impr *inMemoryProductRepository) GetProducts(_ context.Context, first int, cursor string) (ProductList, error) {
	products := make([]common.Product, 0)

	newCursor := cursor
	reachedCursor := false
	if cursor == "" {
		reachedCursor = true
	}

	for _, product := range impr.products {
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

func (impr *inMemoryProductRepository) GetProduct(_ context.Context, id string) (*common.Product, error) {
	return findProductById(impr.products, id)
}

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

	if err != nil {
		sort.Sort(sortDefault(products))
	}
	return products, err
}
