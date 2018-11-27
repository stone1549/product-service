package repository

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/shopspring/decimal"
	"github.com/stone1549/product-service/common"
	"strconv"
	"strings"
)

const (
	listProductsQuery = "SELECT id, name, description, short_description, display_image, thumbnail, price, qty_in_stock, created_at, updated_at FROM product ORDER BY %s LIMIT $1 OFFSET $2"
	getProductQuery   = `SELECT id, name, description, short_description, display_image, thumbnail, price, qty_in_stock, 
						created_at, updated_at FROM product WHERE id=$1`
	insertProductQuery = `INSERT INTO product (id, name, description, short_description, display_image, thumbnail, 
							price, qty_in_stock) 
						  	VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	searchProductQuery = `SELECT id, name, description, short_description, display_image, thumbnail, price, 
							qty_in_stock, created_at, updated_at FROM product WHERE 
								textsearchable_index_col @@ to_tsquery($1) 
							ORDER BY textsearchable_index_col 
							LIMIT $2 OFFSET $3`
)

type postgresqlProductRepository struct {
	db *sql.DB
}

func scanProductFromRow(row *sql.Row) (*common.Product, error) {
	var result common.Product

	var priceStr string
	err := row.Scan(&result.Id, &result.Name, &result.Description, &result.ShortDescription, &result.DisplayImage,
		&result.Thumbnail, &priceStr, &result.QtyInStock, &result.CreatedAt, &result.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
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
		&result.Thumbnail, &priceStr, &result.QtyInStock, &result.CreatedAt, &result.UpdatedAt)

	if priceStr != "" {
		price, err := decimal.NewFromString(priceStr)

		if err != nil {
			return nil, err
		}

		result.Price = &price
	}

	return &result, err
}

func orderByFields(orderBy common.OrderBy) (string, error) {
	var keys []string
	for _, key := range orderBy.Order() {
		switch key {
		case common.OrderByCreated:
			keys = append(keys, "created_at")
		case common.OrderByCreatedDesc:
			keys = append(keys, "created_at DESC")
		case common.OrderByUpdated:
			keys = append(keys, "updated_at")
		case common.OrderByUpdatedDesc:
			keys = append(keys, "updated_at DESC")
		case common.OrderByName:
			keys = append(keys, "name")
		case common.OrderByNameDesc:
			keys = append(keys, "name DESC")
		case common.OrderByPrice:
			keys = append(keys, "price")
		case common.OrderByPriceDesc:
			keys = append(keys, "price DESC")
		default:
			return "", newErrRepository(fmt.Sprintf("Unsupported order by field %s", key))
		}
	}

	return strings.Join(keys, ", "), nil
}

func (ppr postgresqlProductRepository) GetProducts(ctx context.Context, first int, cursor string,
	orderBy common.OrderBy) (ProductList, error) {
	var result ProductList
	var offset int
	var err error

	if strings.TrimSpace(cursor) != "" {
		offset, err = strconv.Atoi(cursor)

		if err != nil {
			return result, newErrRepository("Invalid cursor")
		}
	}

	orderByStr, err := orderByFields(orderBy)

	if err != nil {
		return ProductList{}, err
	}

	rows, err := ppr.db.QueryContext(ctx, fmt.Sprintf(listProductsQuery, orderByStr), first, offset)
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

func (ppr postgresqlProductRepository) GetProduct(ctx context.Context, id string) (*common.Product, error) {
	row := ppr.db.QueryRowContext(ctx, getProductQuery, id)

	if row == nil {
		return nil, nil
	}

	return scanProductFromRow(row)
}

func (ppr *postgresqlProductRepository) SearchProducts(ctx context.Context, searchTxt string, first int,
	cursor string) (ProductList, error) {
	var result ProductList
	var offset int
	var err error

	if strings.TrimSpace(cursor) != "" {
		offset, err = strconv.Atoi(cursor)

		if err != nil {
			return result, newErrRepository("Invalid cursor")
		}
	}

	// TODO: handle tokenizing searchTxt or require clients to use PG syntax?
	rows, err := ppr.db.QueryContext(ctx, searchProductQuery, searchTxt, first, offset)
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

func loadInitPostgresqlData(db *sql.DB, dataset string) error {
	products, err := loadInitInMemoryDataset(dataset)

	if err != nil {
		return err
	}

	txn, err := db.Begin()

	if err != nil {
		return err
	}

	for _, product := range products {
		_, err = txn.Exec(insertProductQuery, product.Id, product.Name, product.Description, product.ShortDescription,
			product.DisplayImage, product.Thumbnail, product.Price.StringFixed(6), product.QtyInStock)

		if err != nil {
			return err
		}
	}

	return txn.Commit()
}

func MakePostgresqlProductRespository(config common.Configuration, db *sql.DB) (ProductRepository, error) {
	var err error
	if config.GetInitDataSet() != "" {
		err = loadInitPostgresqlData(db, config.GetInitDataSet())
	}

	if err != nil {
		return nil, err
	}

	return &postgresqlProductRepository{db}, nil
}
