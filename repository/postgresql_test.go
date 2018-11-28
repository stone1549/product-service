package repository_test

import (
	"context"
	"database/sql"
	"errors"
	"github.com/stone1549/product-service/common"
	"github.com/stone1549/product-service/repository"
	"gopkg.in/DATA-DOG/go-sqlmock.v2"
	"testing"
	"time"
)

func mockExpectExecTimes(mock sqlmock.Sqlmock, sqlRegexStr string, times int) {
	for i := 0; i < times; i++ {
		mock.ExpectExec(sqlRegexStr).WillReturnResult(sqlmock.NewResult(int64(i), 1))
	}
}

func makeAndTestPgSmallRepo() (*sql.DB, sqlmock.Sqlmock, repository.ProductRepository, error) {
	var err error
	db, mock, err := sqlmock.New()
	if err != nil {
		return nil, nil, nil, err
	}

	mock.ExpectBegin()
	mockExpectExecTimes(mock, "INSERT INTO product", 20)
	mock.ExpectCommit()
	repo, err := repository.MakePostgresqlProductRespository(pgSmall, db)

	return db, mock, repo, err
}

// TestMakePostgresqlProductRespository_Ds ensures that a dataset can be loaded when a pg repo is constructed.
func TestMakePostgresqlProductRespository_Ds(t *testing.T) {
	db, mock, _, err := makeAndTestPgSmallRepo()
	defer db.Close()
	ok(t, err)
	ok(t, mock.ExpectationsWereMet())
}

// TestMakePostgresqlProductRespository ensures that an empty pg repo can be constructed.
func TestMakePostgresqlProductRespository(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	_, err = repository.MakePostgresqlProductRespository(pgEmpty, db)
	ok(t, err)
	ok(t, mock.ExpectationsWereMet())
}

func getProductColumns() []string {
	columns := make([]string, 0)
	columns = append(columns, "id")
	columns = append(columns, "name")
	columns = append(columns, "description")
	columns = append(columns, "short_description")
	columns = append(columns, "display_image")
	columns = append(columns, "thumbnail")
	columns = append(columns, "price")
	columns = append(columns, "qty_in_stock")
	columns = append(columns, "created_at")
	columns = append(columns, "updated_at")
	return columns
}
func addExpectedProductId1Row(rows *sqlmock.Rows) *sqlmock.Rows {
	createdAt, _ := time.Parse("2006-01-15T15:20:59", "2017-01-01T00:00:00Z")
	updatedAt, _ := time.Parse("2006-01-15T15:20:59", "2018-01-01T00:00:20Z")
	return rows.AddRow(
		"1",
		"Portal Gun",
		"The Portal Gun is a gadget that allows the user(s) to travel between different universes/dimensions/"+
			"realities.\n\nThe Gun was likely created by a Rick, although it is unknown which one; if there is any "+
			"truth to C-137's fabricated origin story, then he may not be the original inventor.",
		"Travel between different dimensions!",
		"https://images-na.ssl-images-amazon.com/images/I/31s7nNMzMUL.jpg",
		"https://images-na.ssl-images-amazon.com/images/I/31s7nNMzMUL.jpg",
		"2499.990000",
		1,
		createdAt,
		updatedAt,
	)
}

func addExpectedProductId2Row(rows *sqlmock.Rows) *sqlmock.Rows {
	createdAt, _ := time.Parse("2006-01-15T15:20:59", "2017-01-01T00:00:01Z")
	updatedAt, _ := time.Parse("2006-01-15T15:20:59", "2018-01-01T00:00:19Z")
	return rows.AddRow(
		"2",
		"Portal Gun",
		"The Portal Gun is a gadget that allows the user(s) to travel between different universes/dimensions/"+
			"realities.\n\nThe Gun was likely created by a Rick, although it is unknown which one; if there is any "+
			"truth to C-137's fabricated origin story, then he may not be the original inventor.",
		"Travel between different dimensions!",
		"https://images-na.ssl-images-amazon.com/images/I/31s7nNMzMUL.jpg",
		"https://images-na.ssl-images-amazon.com/images/I/31s7nNMzMUL.jpg",
		"2499.990000",
		1,
		createdAt,
		updatedAt,
	)
}

func addExpectedProductId3Row(rows *sqlmock.Rows) *sqlmock.Rows {
	createdAt, _ := time.Parse("2006-01-15T15:20:59", "2017-01-01T00:00:02Z")
	updatedAt, _ := time.Parse("2006-01-15T15:20:59", "2018-01-01T00:00:18Z")
	return rows.AddRow(
		"3",
		"Portal Gun",
		"The Portal Gun is a gadget that allows the user(s) to travel between different universes/dimensions/"+
			"realities.\n\nThe Gun was likely created by a Rick, although it is unknown which one; if there is any "+
			"truth to C-137's fabricated origin story, then he may not be the original inventor.",
		"Travel between different dimensions!",
		"https://images-na.ssl-images-amazon.com/images/I/31s7nNMzMUL.jpg",
		"https://images-na.ssl-images-amazon.com/images/I/31s7nNMzMUL.jpg",
		"2499.990000",
		10,
		createdAt,
		updatedAt,
	)
}

func addExpectedProductId4Row(rows *sqlmock.Rows) *sqlmock.Rows {
	createdAt, _ := time.Parse("2006-01-15T15:20:59", "2017-01-01T00:00:03Z")
	updatedAt, _ := time.Parse("2006-01-15T15:20:59", "2018-01-01T00:00:17Z")
	return rows.AddRow(
		"4",
		"Portal Gun",
		"The Portal Gun is a gadget that allows the user(s) to travel between different universes/dimensions/"+
			"realities.\n\nThe Gun was likely created by a Rick, although it is unknown which one; if there is any "+
			"truth to C-137's fabricated origin story, then he may not be the original inventor.",
		"Travel between different dimensions!",
		"https://images-na.ssl-images-amazon.com/images/I/31s7nNMzMUL.jpg",
		"https://images-na.ssl-images-amazon.com/images/I/31s7nNMzMUL.jpg",
		"2499.990000",
		1,
		createdAt,
		updatedAt,
	)
}

func addExpectedProductId5Row(rows *sqlmock.Rows) *sqlmock.Rows {
	createdAt, _ := time.Parse("2006-01-15T15:20:59", "2017-01-01T00:00:04Z")
	updatedAt, _ := time.Parse("2006-01-15T15:20:59", "2018-01-01T00:00:16Z")
	return rows.AddRow(
		"5",
		"Portal Gun",
		"The Portal Gun is a gadget that allows the user(s) to travel between different universes/dimensions/"+
			"realities.\n\nThe Gun was likely created by a Rick, although it is unknown which one; if there is any "+
			"truth to C-137's fabricated origin story, then he may not be the original inventor.",
		"Travel between different dimensions!",
		"https://images-na.ssl-images-amazon.com/images/I/31s7nNMzMUL.jpg",
		"https://images-na.ssl-images-amazon.com/images/I/31s7nNMzMUL.jpg",
		"2499.990000",
		1,
		createdAt,
		updatedAt,
	)
}

func newProductRows() *sqlmock.Rows {
	return sqlmock.NewRows(getProductColumns())
}

// TestGetProduct_PgSuccessWithResult ensures that a product can be retrieved by its id.
func TestGetProduct_PgSuccessWithResult(t *testing.T) {
	db, mock, repo, err := makeAndTestPgSmallRepo()
	defer db.Close()
	ok(t, err)

	mock.ExpectQuery("SELECT .* FROM product WHERE id=\\$1").
		WithArgs("1").
		WillReturnRows(addExpectedProductId1Row(newProductRows()))
	product, err := repo.GetProduct(context.Background(), "1")

	ok(t, err)
	assert(t, product != nil, "Expected product to not be nil")
	ok(t, mock.ExpectationsWereMet())
}

// TestGetProduct_PgSuccessWithNoResult ensures that attempting to retrieve a product that does not exist will return
// nil.
func TestGetProduct_PgSuccessWithNoResult(t *testing.T) {
	db, mock, repo, err := makeAndTestPgSmallRepo()
	defer db.Close()
	ok(t, err)

	mock.ExpectQuery("SELECT .* FROM product WHERE id=\\$1").WithArgs("1").
		WillReturnRows(newProductRows())
	product, err := repo.GetProduct(context.Background(), "1")

	ok(t, err)
	assert(t, product == nil, "expected product to be nil")
	ok(t, mock.ExpectationsWereMet())
}

// TestGetProduct_PgError ensures that an error will be returned if a query to PG fails.
func TestGetProduct_PgError(t *testing.T) {
	db, mock, repo, err := makeAndTestPgSmallRepo()
	defer db.Close()
	ok(t, err)

	mock.ExpectQuery("SELECT .* FROM product WHERE id=\\$1").WithArgs("1").
		WillReturnError(errors.New("test mock error"))
	_, err = repo.GetProduct(context.Background(), "1")

	notOk(t, err)
	ok(t, mock.ExpectationsWereMet())
}

const getProductsRegexStr = "SELECT .* FROM product .* LIMIT \\$1 OFFSET \\$2"

// TestGetProducts_PgSuccessWithPartialResults ensures that a partial set of products will be returned when appropriate.
func TestGetProducts_PgSuccessWithPartialResults(t *testing.T) {
	db, mock, repo, err := makeAndTestPgSmallRepo()
	defer db.Close()
	ok(t, err)

	mock.ExpectQuery(getProductsRegexStr).
		WithArgs(5, 0).
		WillReturnRows(addExpectedProductId2Row(addExpectedProductId1Row(newProductRows())))
	products, err := repo.GetProducts(context.Background(), 5, "", common.OrderBy{})

	ok(t, err)
	equals(t, 2, len(products.Products))
	equals(t, "2", products.Cursor)
	ok(t, mock.ExpectationsWereMet())
}

// TestGetProducts_PgSuccessWithFullResults ensures that a full set of products will be returned when appropriate.
func TestGetProducts_PgSuccessWithFullResults(t *testing.T) {
	db, mock, repo, err := makeAndTestPgSmallRepo()
	defer db.Close()
	ok(t, err)

	expRows := addExpectedProductId5Row(addExpectedProductId4Row(addExpectedProductId3Row(
		addExpectedProductId2Row(addExpectedProductId1Row(newProductRows())))))
	mock.ExpectQuery(getProductsRegexStr).
		WithArgs(5, 0).
		WillReturnRows(expRows)
	products, err := repo.GetProducts(context.Background(), 5, "", common.OrderBy{})

	ok(t, err)
	equals(t, 5, len(products.Products))
	equals(t, "5", products.Cursor)
	ok(t, mock.ExpectationsWereMet())
}

// TestGetProducts_PgSuccessWithFullResultsEmptyPageTwo ensures that an empty set of products will be returned when
// appropriate.
func TestGetProducts_PgSuccessWithFullResultsEmptyPageTwo(t *testing.T) {
	db, mock, repo, err := makeAndTestPgSmallRepo()
	defer db.Close()
	ok(t, err)

	mock.ExpectQuery(getProductsRegexStr).
		WithArgs(5, 5).
		WillReturnRows(newProductRows())
	products, err := repo.GetProducts(context.Background(), 5, "5", common.OrderBy{})

	ok(t, err)
	equals(t, 0, len(products.Products))
	equals(t, "5", products.Cursor)
	ok(t, mock.ExpectationsWereMet())
}

// TestGetProducts_PgError ensures that an error will be returned if there is a problem querying PG.
func TestGetProducts_PgError(t *testing.T) {
	db, mock, repo, err := makeAndTestPgSmallRepo()
	defer db.Close()
	ok(t, err)

	mock.ExpectQuery(getProductsRegexStr).
		WithArgs(5, 5).
		WillReturnError(errors.New("test error"))
	_, err = repo.GetProducts(context.Background(), 5, "5", common.OrderBy{})

	notOk(t, err)
	ok(t, mock.ExpectationsWereMet())
}

// TestSearchProducts_PgSuccessWithPartialResults ensures that a partial set will be returned when appropriate.
func TestSearchProducts_PgSuccessWithPartialResults(t *testing.T) {
	db, mock, repo, err := makeAndTestPgSmallRepo()
	defer db.Close()
	ok(t, err)

	mock.ExpectQuery("SELECT .* FROM product").
		WithArgs("portal", 5, 0).
		WillReturnRows(addExpectedProductId2Row(addExpectedProductId1Row(newProductRows())))
	products, err := repo.SearchProducts(context.Background(), "portal", 5, "")

	ok(t, err)
	equals(t, 2, len(products.Products))
	equals(t, "2", products.Cursor)
	ok(t, mock.ExpectationsWereMet())
}

// TestSearchProducts_PgSuccessWithFullResults ensures that a full set will be returned when appropriate.
func TestSearchProducts_PgSuccessWithFullResults(t *testing.T) {
	db, mock, repo, err := makeAndTestPgSmallRepo()
	defer db.Close()
	ok(t, err)

	expRows := addExpectedProductId5Row(addExpectedProductId4Row(addExpectedProductId3Row(
		addExpectedProductId2Row(addExpectedProductId1Row(newProductRows())))))
	mock.ExpectQuery("SELECT .* FROM product").
		WithArgs("portal", 5, 0).
		WillReturnRows(expRows)
	products, err := repo.SearchProducts(context.Background(), "portal", 5, "")

	ok(t, err)
	equals(t, 5, len(products.Products))
	equals(t, "5", products.Cursor)
	ok(t, mock.ExpectationsWereMet())
}

// TestSearchProducts_PgSuccessWithFullResultsEmptyPageTwo ensures that an empty set will be returned when appropriate.
func TestSearchProducts_PgSuccessWithFullResultsEmptyPageTwo(t *testing.T) {
	db, mock, repo, err := makeAndTestPgSmallRepo()
	defer db.Close()
	ok(t, err)

	mock.ExpectQuery("SELECT .* FROM product").
		WithArgs("portal", 5, 5).
		WillReturnRows(newProductRows())
	products, err := repo.SearchProducts(context.Background(), "portal", 5, "5")

	ok(t, err)
	equals(t, 0, len(products.Products))
	equals(t, "5", products.Cursor)
	ok(t, mock.ExpectationsWereMet())
}

// TestSearchProducts_PgSuccessWithFullResultsEmptyPageTwo ensures that an error will be returned if there is a problem
// querying PG.
func TestSearchProducts_PgError(t *testing.T) {
	db, mock, repo, err := makeAndTestPgSmallRepo()
	defer db.Close()
	ok(t, err)

	mock.ExpectQuery("SELECT .* FROM product").
		WithArgs("portal", 5, 5).
		WillReturnError(errors.New("test error"))
	_, err = repo.SearchProducts(context.Background(), "portal", 5, "5")

	notOk(t, err)
	ok(t, mock.ExpectationsWereMet())
}
