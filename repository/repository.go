package repository

import (
	"context"
	"database/sql"
	"github.com/stone1549/product-service/common"
)

// ProductList holds a slice of products and a cursor that can be used to retrieve more results
type ProductList struct {
	Products []common.Product
	Cursor   string
}

// ProductRepository represents a data source through which products can be retrieved.
type ProductRepository interface {
	// GetProducts retrieves a list of the first X products starting from the given cursor.
	GetProducts(ctx context.Context, first int, cursor string, orderBy common.OrderBy) (ProductList, error)
	// GetProduct retrieves a product from the given id.
	GetProduct(ctx context.Context, id string) (*common.Product, error)
	// SearchProducts retrieves the first X matches starting from the given cursor.
	SearchProducts(ctx context.Context, searchTxt string, first int, cursor string) (ProductList, error)
}

// NewProductRepository constructs a ProductRepository from the given configuration.
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
