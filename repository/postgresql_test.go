package repository_test

import (
	"context"
	"database/sql"
	"errors"
	"github.com/stone1549/product-service/repository"
	"gopkg.in/DATA-DOG/go-sqlmock.v2"
	"testing"
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

func TestMakePostgresqlProductRespository_DsSmall(t *testing.T) {
	db, mock, _, err := makeAndTestPgSmallRepo()
	defer db.Close()
	ok(t, err)
	ok(t, mock.ExpectationsWereMet())
}

func TestMakePostgresqlProductRespository_DsEmpty(t *testing.T) {
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
	return columns
}
func addExpectedProductId1Row(rows *sqlmock.Rows) *sqlmock.Rows {
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
	)
}

func newProductRows() *sqlmock.Rows {
	return sqlmock.NewRows(getProductColumns())
}

func TestGetProduct_SuccessWithResult(t *testing.T) {
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

func TestGetProduct_SuccessWithNoResult(t *testing.T) {
	db, mock, repo, err := makeAndTestPgSmallRepo()
	defer db.Close()
	ok(t, err)

	mock.ExpectQuery("SELECT .* FROM product WHERE id=\\$1").WithArgs("1").
		WillReturnRows(newProductRows())
	product, err := repo.GetProduct(context.Background(), "1")

	ok(t, err)
	assert(t, product == nil, "Expected product to be nil")
	ok(t, mock.ExpectationsWereMet())
}

func TestGetProduct_Error(t *testing.T) {
	db, mock, repo, err := makeAndTestPgSmallRepo()
	defer db.Close()
	ok(t, err)

	mock.ExpectQuery("SELECT .* FROM product WHERE id=\\$1").WithArgs("1").
		WillReturnError(errors.New("test mock error"))
	_, err = repo.GetProduct(context.Background(), "1")

	notOk(t, err)
	ok(t, mock.ExpectationsWereMet())
}
