package repository

import (
	"context"
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"github.com/stone1549/product-service/common"
	"strconv"
	"strings"
)

const (
	listProductsQuery = "SELECT id, name, description, short_description, display_image, thumbnail, price, qty_in_stock FROM product LIMIT $1 OFFSET $2"
	getProductQuery   = "SELECT id, name, description, short_description, display_image, thumbnail, price, qty_in_stock FROM product WHERE id=$1"
)

type postgresqlProductRepository struct {
	db *sql.DB
}

func scanProductFromRow(row *sql.Row) (*common.Product, error) {
	var result common.Product

	var priceStr string
	err := row.Scan(&result.Id, &result.Name, &result.Description, &result.ShortDescription, &result.DisplayImage,
		&result.Thumbnail, &priceStr, &result.Quantity)

	if priceStr != "" {
		price, err := decimal.NewFromString(priceStr)

		if err != nil {
			return nil, err
		}

		result.Price = &price
	}

	return &result, err
}

func scanProductFromRows(rows *sql.Rows) (*common.Product, error) {
	var result common.Product

	var priceStr string
	err := rows.Scan(&result.Id, &result.Name, &result.Description, &result.ShortDescription, &result.DisplayImage,
		&result.Thumbnail, &priceStr, &result.Quantity)

	if priceStr != "" {
		price, err := decimal.NewFromString(priceStr)

		if err != nil {
			return nil, err
		}

		result.Price = &price
	}

	return &result, err
}

func (ppr postgresqlProductRepository) ProductsFromRepo(ctx context.Context, first int, cursor string) (productList, error) {
	var result productList
	var offset int
	var err error

	if strings.TrimSpace(cursor) != "" {
		offset, err = strconv.Atoi(cursor)

		if err != nil {
			errors.New("Invalid cursor")
		}
	}

	rows, err := ppr.db.QueryContext(ctx, listProductsQuery, first, offset)
	if err != nil {
		return result, err
	}
	defer rows.Close()

	result.Products = make([]common.Product, 0)
	for rows.Next() {
		product, err := scanProductFromRows(rows)

		if err != nil {
			return result, err
		}

		result.Products = append(result.Products, *product)
	}
	if err := rows.Err(); err != nil {
		return result, err
	}

	result.Cursor = strconv.Itoa(offset + len(result.Products))
	return result, nil
}

func (ppr postgresqlProductRepository) ProductFromRepo(ctx context.Context, id string) (common.Product, error) {
	result, err := scanProductFromRow(ppr.db.QueryRowContext(ctx, getProductQuery, id))

	if err != nil {
		return common.Product{}, err
	}

	return *result, err
}

func makePostgresqlProductRespository(config common.Configuration) (ProductRepository, error) {
	db, err := sql.Open("postgres", config.GetPgUrl())

	if err != nil {
		return nil, err
	}

	return &postgresqlProductRepository{db}, nil
}
