package repository

import (
	"github.com/stone1549/product-service/common"
)

type productList struct {
	Products []*common.Product
	Cursor   string
}

type ProductRepository interface {
	ProductsFromRepo(first int, cursor string) (productList, error)
	ProductFromRepo(id string) (common.Product, error)
}

var repo ProductRepository

func GetProductRepository() (ProductRepository, error) {
	if repo == nil {
		return nil, ErrRepository("ConfigureProductRepository must be called first")
	}

	return repo, nil
}

func ConfigureProductRepository(repoType common.ProductRepositoryType) error {
	var err error
	if repo == nil {
		switch repoType {
		case common.InMemory:
			repo = &InMemoryProductRepository{make([]common.Product, 0)}
		case common.PostgreSQL:
			err = ErrRepository("PostgreSQL repository type unimplemented")
		default:
			err = ErrRepository("repository type unimplemented")
		}
	} else {
		err = ErrRepository("ConfigureProductRepository called twice")
	}

	return err
}
