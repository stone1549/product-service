package repository

import (
	"context"
	"database/sql"
	"github.com/stone1549/product-service/common"
)

type ProductList struct {
	Products []common.Product
	Cursor   string
}

type ProductRepository interface {
	GetProducts(ctx context.Context, first int, cursor string, orderBy common.OrderBy) (ProductList, error)
	GetProduct(ctx context.Context, id string) (*common.Product, error)
	SearchProducts(ctx context.Context, searchTxt string, first int, cursor string) (ProductList, error)
}

func NewProductRepository(config common.Configuration) (ProductRepository, error) {
	var err error
	var repo ProductRepository
	var db *sql.DB
	switch config.GetRepoType() {
	case common.InMemoryRepo:
		repo, err = MakeInMemoryRepository(config)
	case common.PostgreSqlRepo:
		db, err = sql.Open("postgres", config.GetPgUrl())

		if err != nil {
			return nil, err
		}
		repo, err = MakePostgresqlProductRespository(config, db)
	default:
		err = newErrRepository("repository type unimplemented")
	}

	return repo, err
}
