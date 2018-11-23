package repository

import "github.com/stone1549/product-service/common"

type InMemoryProductRepository struct {
	products []common.Product
}

func (impr *InMemoryProductRepository) ProductsFromRepo(first int, cursor string) (productList, error) {
	productPointers := make([]*common.Product, 0)

	for _, product := range impr.products {
		productCopy := product
		productPointers = append(productPointers, &productCopy)
	}
	return productList{productPointers, ""}, nil
}

func findProductById(products []common.Product, id string) (common.Product, error) {
	for _, product := range products {
		if id == product.Id {
			return product, nil
		}
	}

	return common.Product{}, nil
}

func (impr *InMemoryProductRepository) ProductFromRepo(id string) (common.Product, error) {
	return findProductById(impr.products, id)
}
