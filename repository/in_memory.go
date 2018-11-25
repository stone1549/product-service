package repository

import (
	"context"
	"encoding/json"
	"github.com/stone1549/product-service/common"
	"io/ioutil"
)

type inMemoryProductRepository struct {
	products []common.Product
}

func (impr *inMemoryProductRepository) ProductsFromRepo(_ context.Context, first int, cursor string) (ProductList, error) {
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

func (impr *inMemoryProductRepository) ProductFromRepo(_ context.Context, id string) (*common.Product, error) {
	return findProductById(impr.products, id)
}

func (impr *inMemoryProductRepository) SearchProducts(ctx context.Context, searchTxt string, first int,
	cursor string) (ProductList, error) {
	return ProductList{}, newErrRepository("Search not implemented for in memory repo")
}

func makeInMemoryRepository(config common.Configuration) (ProductRepository, error) {
	var products []common.Product
	var err error

	switch config.GetInitDataSet() {
	case common.NoDataset:
		products = make([]common.Product, 0)
	case common.SmallDataset:
		products, err = loadInitInMemoryDataset(config.GetInitDataSet())
	default:
		err = newErrRepository("Unsupported dataset %s for repo type PostgreSQL")
	}

	return &inMemoryProductRepository{products}, err
}

func loadInitInMemoryDataset(dataset common.InitDataset) ([]common.Product, error) {
	var err error
	products := make([]common.Product, 0)

	var filename string
	switch dataset {
	case common.SmallDataset:
		filename = "small_set.json"
	default:
		err = newErrRepository("Unsupported dataset %s for repo type PostgreSQL")
	}

	if err != nil {
		return products, err
	}

	jsonBytes, err := ioutil.ReadFile("data/" + filename)
	if err != nil {
		return products, err
	}

	err = json.Unmarshal(jsonBytes, &products)

	return products, err
}
