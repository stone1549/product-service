package repository_test

import (
	"context"
	"github.com/stone1549/product-service/common"
	"github.com/stone1549/product-service/repository"
	"testing"
)

func makeNewImRepo(t *testing.T) repository.ProductRepository {
	repo, err := repository.MakeInMemoryRepository(inMemorySmall)

	ok(t, err)
	return repo
}

// TestGetProduct_ImSuccessWithResult ensures a product can be retrieved from the repo by id.
func TestGetProduct_ImSuccessWithResult(t *testing.T) {
	repo := makeNewImRepo(t)
	product, err := repo.GetProduct(context.Background(), "1")

	ok(t, err)
	assert(t, product != nil, "Expected product to not be nil")
	equals(t, "1", product.Id)
}

// TestGetProduct_ImSuccessWithResult ensures that nil is returned if the requested product does not exist.
func TestGetProduct_ImSuccessWithNoResult(t *testing.T) {
	repo := makeNewImRepo(t)
	product, err := repo.GetProduct(context.Background(), "A")

	ok(t, err)
	assert(t, product == nil, "expected product to be nil")
}

// TestGetProducts_ImSuccessWithPartialResults ensures that a partial list of results will be returned when necessary.
func TestGetProducts_ImSuccessWithPartialResults(t *testing.T) {
	repo := makeNewImRepo(t)
	products, err := repo.GetProducts(context.Background(), 5, "18", common.OrderBy{})

	ok(t, err)
	equals(t, 2, len(products.Products))
	equals(t, "20", products.Cursor)
}

// TestGetProducts_ImSuccessOrderByCreated ensures that products can be retrieved properly sorted by a single key.
func TestGetProducts_ImSuccessOrderByCreated(t *testing.T) {
	repo := makeNewImRepo(t)

	orderBy := common.OrderBy{}
	orderBy.Add(common.OrderByCreatedDesc)
	products, err := repo.GetProducts(context.Background(), 5, "", orderBy)

	ok(t, err)
	equals(t, 5, len(products.Products))
	equals(t, "16", products.Cursor)
}

// TestGetProducts_ImSuccessOrderByCreatedAndName ensures that products can be retrieved properly sorted with multiple
// keys.
func TestGetProducts_ImSuccessOrderByCreatedAndName(t *testing.T) {
	repo := makeNewImRepo(t)

	orderBy := common.OrderBy{}
	orderBy.Add(common.OrderByCreatedDesc)
	orderBy.Add(common.OrderByName)
	products, err := repo.GetProducts(context.Background(), 2, "", orderBy)

	ok(t, err)
	equals(t, 2, len(products.Products))
	equals(t, "20", products.Cursor)
}

// TestGetProducts_ImSuccessOrderByName ensures that products can be retrieved properly sorted by a text field.
func TestGetProducts_ImSuccessOrderByName(t *testing.T) {
	repo := makeNewImRepo(t)

	orderBy := common.OrderBy{}
	orderBy.Add(common.OrderByName)
	products, err := repo.GetProducts(context.Background(), 5, "", orderBy)

	ok(t, err)
	equals(t, 5, len(products.Products))
	equals(t, "11", products.Cursor)
}

// TestGetProduct_ImSuccessWithFullResults ensures that a full set of products will be returned where appropriate.
func TestGetProduct_ImSuccessWithFullResults(t *testing.T) {
	repo := makeNewImRepo(t)

	products, err := repo.GetProducts(context.Background(), 5, "", common.OrderBy{})

	ok(t, err)
	equals(t, 5, len(products.Products))
	equals(t, "5", products.Cursor)
}

// TestGetProduct_ImSuccessWithFullResults ensures that an empty product set will be returned when appropriate.
func TestGetProduct_ImSuccessWithFullResultsEmptyPageTwo(t *testing.T) {
	repo := makeNewImRepo(t)
	products, err := repo.GetProducts(context.Background(), 5, "21", common.OrderBy{})

	ok(t, err)
	equals(t, 0, len(products.Products))
	equals(t, "21", products.Cursor)
}

// TestSearchProducts_ImSuccessWithPartialResults ensures that a partial product set will be returned when appropriate.
func TestSearchProducts_ImSuccessWithPartialResults(t *testing.T) {
	repo := makeNewImRepo(t)
	products, err := repo.SearchProducts(context.Background(), "portal OR shrink", 5, "")

	ok(t, err)
	equals(t, 2, len(products.Products))
	equals(t, "14", products.Cursor)
}

// TestSearchProducts_ImSuccessWithFullResults ensures that a full product set will be returned when appropriate.
func TestSearchProducts_ImSuccessWithFullResults(t *testing.T) {
	repo := makeNewImRepo(t)
	products, err := repo.SearchProducts(context.Background(), "portal OR time OR ray", 5, "")

	ok(t, err)
	equals(t, 5, len(products.Products))
	equals(t, "14", products.Cursor)
}

// TestSearchProducts_ImSuccessWithFullResultsEmptyPageTwo ensures that an empty product set will be returned when
// appropriate.
func TestSearchProducts_ImSuccessWithFullResultsEmptyPageTwo(t *testing.T) {
	repo := makeNewImRepo(t)
	products, err := repo.SearchProducts(context.Background(), "portal", 5, "20")

	ok(t, err)
	equals(t, 0, len(products.Products))
	equals(t, "20", products.Cursor)
}
