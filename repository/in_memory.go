package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/blevesearch/bleve"
	"github.com/stone1549/product-service/common"
	"io/ioutil"
)

type inMemoryProductRepository struct {
	products []common.Product
	index    bleve.Index
}

func (impr *inMemoryProductRepository) GetProducts(_ context.Context, first int, cursor string) (ProductList, error) {
	products := make([]common.Product, 0)

	for _, product := range impr.products {
		productCopy := product
		products = append(products, productCopy)
	}
	return ProductList{products, ""}, nil
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
	var newCursor string
	for _, hit := range searchResults.Hits {
		if len(products) == first {
			break
		} else if hit.ID == cursor || cursor == "" || len(products) > 1 {
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

func makeInMemoryRepository(config common.Configuration) (ProductRepository, error) {
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

	switch config.GetInitDataSet() {
	case common.NoDataset:
		products = make([]common.Product, 0)
	case common.SmallDataset:
		products, err = loadInitInMemoryDataset(config.GetInitDataSet())
	default:
		err = newErrRepository("Unsupported dataset %s for repo type PostgreSqlRepo")
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

func loadInitInMemoryDataset(dataset common.InitDataset) ([]common.Product, error) {
	var err error
	products := make([]common.Product, 0)

	var filename string
	switch dataset {
	case common.SmallDataset:
		filename = "small_set.json"
	default:
		err = newErrRepository("Unsupported dataset %s for repo type PostgreSqlRepo")
	}

	if err != nil {
		return products, err
	}

	jsonBytes, err := ioutil.ReadFile("../data/" + filename)
	if err != nil {
		return products, err
	}

	err = json.Unmarshal(jsonBytes, &products)

	return products, err
}
