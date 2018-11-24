package repository

import (
	"context"
	"github.com/stone1549/product-service/common"
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
