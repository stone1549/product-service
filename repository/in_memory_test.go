package repository_test

import (
	"context"
	"github.com/stone1549/product-service/repository"
	"testing"
)

func makeNewImRepo(t *testing.T) repository.ProductRepository {
	repo, err := repository.MakeInMemoryRepository(inMemorySmall)

	ok(t, err)
	return repo
}

func TestGetProduct_ImSuccessWithResult(t *testing.T) {
	repo := makeNewImRepo(t)
	product, err := repo.GetProduct(context.Background(), "1")

	ok(t, err)
	assert(t, product != nil, "Expected product to not be nil")
}

func TestGetProduct_ImSuccessWithNoResult(t *testing.T) {
	repo := makeNewImRepo(t)
	product, err := repo.GetProduct(context.Background(), "A")

	ok(t, err)
	assert(t, product == nil, "expected product to be nil")
}

func TestGetProducts_ImSuccessWithPartialResults(t *testing.T) {
	repo := makeNewImRepo(t)
	products, err := repo.GetProducts(context.Background(), 5, "18")

	ok(t, err)
	equals(t, 2, len(products.Products))
	equals(t, "20", products.Cursor)
}

func TestGetProduct_ImSuccessWithFullResults(t *testing.T) {
	repo := makeNewImRepo(t)

	products, err := repo.GetProducts(context.Background(), 5, "")

	ok(t, err)
	equals(t, 5, len(products.Products))
	equals(t, "5", products.Cursor)
}

func TestGetProduct_ImSuccessWithFullResultsEmptyPageTwo(t *testing.T) {
	repo := makeNewImRepo(t)
	products, err := repo.GetProducts(context.Background(), 5, "21")

	ok(t, err)
	equals(t, 0, len(products.Products))
	equals(t, "21", products.Cursor)
}

func TestSearchProducts_ImSuccessWithPartialResults(t *testing.T) {
	repo := makeNewImRepo(t)
	products, err := repo.SearchProducts(context.Background(), "portal OR shrink", 5, "")

	ok(t, err)
	equals(t, 2, len(products.Products))
	equals(t, "14", products.Cursor)
}

func TestSearchProducts_ImSuccessWithFullResults(t *testing.T) {
	repo := makeNewImRepo(t)
	products, err := repo.SearchProducts(context.Background(), "portal OR time OR ray", 5, "")

	ok(t, err)
	equals(t, 5, len(products.Products))
	equals(t, "14", products.Cursor)
}

func TestSearchProducts_ImSuccessWithFullResultsEmptyPageTwo(t *testing.T) {
	repo := makeNewImRepo(t)
	products, err := repo.SearchProducts(context.Background(), "portal", 5, "20")

	ok(t, err)
	equals(t, 0, len(products.Products))
	equals(t, "20", products.Cursor)
}
