package repository_test

import (
	"fmt"
	"github.com/stone1549/product-service/common"
	"github.com/stone1549/product-service/repository"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"
	"time"
)

// assert fails the test if the condition is false.
func assert(tb testing.TB, condition bool, msg string, v ...interface{}) {
	if !condition {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d: "+msg+"\033[39m\n\n", append([]interface{}{filepath.Base(file), line}, v...)...)
		tb.FailNow()
	}
}

// ok fails the test if an err is not nil.
func ok(tb testing.TB, err error) {
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d: unexpected error: %s\033[39m\n\n", filepath.Base(file), line, err.Error())
		tb.FailNow()
	}
}

// notOk fails the test if an err is nil.
func notOk(tb testing.TB, err error) {
	if err == nil {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d: unexpected lack of error: \033[39m\n\n", filepath.Base(file), line)
		tb.FailNow()
	}
}

// equals fails the test if exp is not equal to act.
func equals(tb testing.TB, exp, act interface{}) {
	if !reflect.DeepEqual(exp, act) {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d:\n\n\texp: %#v\n\n\tgot: %#v\033[39m\n\n", filepath.Base(file), line, exp, act)
		tb.FailNow()
	}
}

type configuration int

const (
	inMemoryEmpty configuration = 0
	inMemorySmall configuration = iota
	pgEmpty       configuration = iota
	pgSmall       configuration = iota
)

func (c configuration) GetLifeCycle() common.LifeCycle {
	return common.DevLifeCycle
}

func (c configuration) GetRepoType() common.ProductRepositoryType {
	switch c {
	case pgSmall:
		fallthrough
	case pgEmpty:
		return common.PostgreSqlRepo
	case inMemorySmall:
		fallthrough
	case inMemoryEmpty:
		return common.InMemoryRepo
	default:
		return common.InMemoryRepo
	}
}

func (c configuration) GetTimeout() time.Duration {
	return 60 * time.Second
}

func (c configuration) GetPort() int {
	return 3333
}

func (c configuration) GetInitDataSet() string {
	switch c {
	case inMemoryEmpty:
		fallthrough
	case pgEmpty:
		return ""
	case inMemorySmall:
		fallthrough
	case pgSmall:
		return "../data/small_set.json"
	default:
		return ""
	}
}

func (c configuration) GetPgUrl() string {
	switch c {
	case inMemoryEmpty:
		fallthrough
	case inMemorySmall:
		return ""
	case pgEmpty:
		fallthrough
	case pgSmall:
		return "postgres://test:test@localhost:5432/postgres?sslmode=disable"
	default:
		return ""
	}
}

// TestNewProductRepository_ImSuccessEmpty ensures an empty in memory repo can be constructed
func TestNewProductRepository_ImSuccessEmpty(t *testing.T) {
	_, err := repository.NewProductRepository(inMemoryEmpty)
	ok(t, err)
}

// TestNewProductRepository_ImSuccessSmall ensures a prepopulated memory repo can be constructed
func TestNewProductRepository_ImSuccessSmall(t *testing.T) {
	_, err := repository.NewProductRepository(inMemorySmall)
	ok(t, err)
}

// TestNewProductRepository_PgSuccessEmpty ensures an empty PG repo can be constructed
func TestNewProductRepository_PgSuccessEmpty(t *testing.T) {
	_, err := repository.NewProductRepository(pgEmpty)
	ok(t, err)
}
