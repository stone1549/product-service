package common_test

import (
	"github.com/stone1549/product-service/common"
	"testing"
)

func TestOrderBy_AddSuccess(t *testing.T) {
	orderBy := common.OrderBy{}
	ok(t, orderBy.Add(common.OrderByCreated))
}

func TestOrderBy_AddMultipleSuccess(t *testing.T) {
	orderBy := common.OrderBy{}
	ok(t, orderBy.Add(common.OrderByCreated))
	ok(t, orderBy.Add(common.OrderByUpdated))
}

func TestOrderBy_AddDuplicateFail(t *testing.T) {
	orderBy := common.OrderBy{}
	ok(t, orderBy.Add(common.OrderByCreated))
	notOk(t, orderBy.Add(common.OrderByCreated))
}

func TestOrderBy_AddConflictingFail(t *testing.T) {
	orderBy := common.OrderBy{}
	ok(t, orderBy.Add(common.OrderByCreated))
	notOk(t, orderBy.Add(common.OrderByCreatedDesc))
}

func TestOrderBy_AddUnknownFail(t *testing.T) {
	orderBy := common.OrderBy{}
	notOk(t, orderBy.Add("Unknown"))
}

func TestOrderBy_OrderSingle(t *testing.T) {
	orderBy := common.OrderBy{}
	ok(t, orderBy.Add(common.OrderByCreated))
	equals(t, 1, len(orderBy.Order()))
	equals(t, common.OrderByCreated, orderBy.Order()[0])
}

func TestOrderBy_OrderMultiple(t *testing.T) {
	orderBy := common.OrderBy{}
	ok(t, orderBy.Add(common.OrderByCreated))
	ok(t, orderBy.Add(common.OrderByUpdated))
	equals(t, 2, len(orderBy.Order()))
	equals(t, common.OrderByCreated, orderBy.Order()[0])
	equals(t, common.OrderByUpdated, orderBy.Order()[1])
}

func TestOrderBy_OrderEmpty(t *testing.T) {
	orderBy := common.OrderBy{}
	equals(t, 2, len(orderBy.Order()))
	equals(t, common.OrderByUpdatedDesc, orderBy.Order()[0])
	equals(t, common.OrderByCreatedDesc, orderBy.Order()[1])
}
