package common_test

import (
	"github.com/stone1549/product-service/common"
	"testing"
)

// TestOrderBy_AddSuccess ensures that a valid key can be successfully added.
func TestOrderBy_AddSuccess(t *testing.T) {
	orderBy := common.OrderBy{}
	ok(t, orderBy.Add(common.OrderByCreated))
}

// TestOrderBy_AddMultipleSuccess ensures that multiple valid keys can be successfully added.
func TestOrderBy_AddMultipleSuccess(t *testing.T) {
	orderBy := common.OrderBy{}
	ok(t, orderBy.Add(common.OrderByCreated))
	ok(t, orderBy.Add(common.OrderByUpdated))
}

// TestOrderBy_AddDuplicateFail ensures that multiple identical keys can not be added.
func TestOrderBy_AddDuplicateFail(t *testing.T) {
	orderBy := common.OrderBy{}
	ok(t, orderBy.Add(common.OrderByCreated))
	notOk(t, orderBy.Add(common.OrderByCreated))
}

// TestOrderBy_AddConflictingFail ensures that conflicting keys can not be added.
func TestOrderBy_AddConflictingFail(t *testing.T) {
	orderBy := common.OrderBy{}
	ok(t, orderBy.Add(common.OrderByCreated))
	notOk(t, orderBy.Add(common.OrderByCreatedDesc))
}

// TestOrderBy_AddUnknownFail ensures that unknown keys can not be added.
func TestOrderBy_AddUnknownFail(t *testing.T) {
	orderBy := common.OrderBy{}
	notOk(t, orderBy.Add("Unknown"))
}

// TestOrderBy_OrderSingle ensures that output is correct for a single key.
func TestOrderBy_Order_Single(t *testing.T) {
	orderBy := common.OrderBy{}
	ok(t, orderBy.Add(common.OrderByCreated))
	equals(t, 1, len(orderBy.Order()))
	equals(t, common.OrderByCreated, orderBy.Order()[0])
}

// TestOrderBy_Order_Multiple ensures that output is correct for multiple keys.
func TestOrderBy_Order_Multiple(t *testing.T) {
	orderBy := common.OrderBy{}
	ok(t, orderBy.Add(common.OrderByCreated))
	ok(t, orderBy.Add(common.OrderByUpdated))
	equals(t, 2, len(orderBy.Order()))
	equals(t, common.OrderByCreated, orderBy.Order()[0])
	equals(t, common.OrderByUpdated, orderBy.Order()[1])
}

// TestOrderBy_OrderEmpty ensures that sensible defaults are returned when no keys are added.
func TestOrderBy_OrderEmpty(t *testing.T) {
	orderBy := common.OrderBy{}
	equals(t, 2, len(orderBy.Order()))
	equals(t, common.OrderByUpdatedDesc, orderBy.Order()[0])
	equals(t, common.OrderByCreatedDesc, orderBy.Order()[1])
}
