package repository

import (
	"context"
	"github.com/stone1549/product-service/common"
)

type ProductList struct {
	Products []common.Product
	Cursor   string
}

type ProductRepository interface {
	ProductsFromRepo(ctx context.Context, first int, cursor string) (ProductList, error)
	ProductFromRepo(ctx context.Context, id string) (*common.Product, error)
	SearchProducts(ctx context.Context, searchTxt string, first int, cursor string) (ProductList, error)
}

var repo ProductRepository

func GetProductRepository() (ProductRepository, error) {
	if repo == nil {
		return nil, newErrRepository("ConfigureProductRepository must be called first")
	}

	return repo, nil
}

func ConfigureProductRepository(config common.Configuration) error {
	var err error
	if repo == nil {
		switch config.GetRepoType() {
		case common.InMemory:
			repo, err = makeInMemoryRepository(config)
		case common.PostgreSQL:
			repo, err = makePostgresqlProductRespository(config)
		default:
			err = newErrRepository("repository type unimplemented")
		}
	} else {
		err = newErrRepository("ConfigureProductRepository called twice")
	}

	return err
}
