package common_test

import (
	"github.com/stone1549/product-service/common"
	"os"
	"testing"
)

const (
	lifeCycleKey      string = "PRODUCT_SERVICE_ENVIRONMENT"
	repoTypeKey       string = "PRODUCT_SERVICE_REPO_TYPE"
	timeoutSecondsKey string = "PRODUCT_SERVICE_TIMEOUT"
	portKey           string = "PRODUCT_SERVICE_PORT"
	pgUrlKey          string = "PRODUCT_SERVICE_PG_URL"
	pgInitDatasetKey  string = "PRODUCT_SERVICE_INIT_DATASET"
)

func clearEnv() {
	os.Setenv(lifeCycleKey, "")
	os.Setenv(repoTypeKey, "")
	os.Setenv(timeoutSecondsKey, "")
	os.Setenv(portKey, "")
	os.Setenv(pgUrlKey, "")
	os.Setenv(pgInitDatasetKey, "")
}

func setEnv(lifeCycle, repoType, timeoutSeconds, port, pgUrl, pgInitDataset string) {
	os.Setenv(lifeCycleKey, lifeCycle)
	os.Setenv(repoTypeKey, repoType)
	os.Setenv(timeoutSecondsKey, timeoutSeconds)
	os.Setenv(portKey, port)
	os.Setenv(pgUrlKey, pgUrl)
	os.Setenv(pgInitDatasetKey, pgInitDataset)
}

// TestGetConfiguration_Defaults ensures that a default configuration is returned if no configuration is provided in
// the environment.
func TestGetConfiguration_Defaults(t *testing.T) {
	clearEnv()
	_, err := common.GetConfiguration()
	ok(t, err)
}

// TestGetConfiguration_Defaults ensures that a default configuration is returned if no configuration is provided in
// the environment.
func TestGetConfiguration_ImSuccess(t *testing.T) {
	setEnv("DEV", "IN_MEMORY", "60", "3333", "", "")
	_, err := common.GetConfiguration()
	ok(t, err)
}

// TestGetConfiguration_ImSuccessSmallDataset ensures that a configuration is returned when specifying an in memory
// repo with an initial dataset.
func TestGetConfiguration_ImSuccessSmallDataset(t *testing.T) {
	setEnv("DEV", "IN_MEMORY", "60", "3333", "", "../data/small_set.json")
	_, err := common.GetConfiguration()
	ok(t, err)
}

// TestGetConfiguration_ImSuccessNoneDataset ensures that a configuration is returned when specifying an in memory
// repo without an initial dataset.
func TestGetConfiguration_ImSuccessNoneDataset(t *testing.T) {
	setEnv("DEV", "IN_MEMORY", "60", "3333", "", "")
	_, err := common.GetConfiguration()
	ok(t, err)
}

// TestGetConfiguration_FailRepo ensures that an error is returned when specifying an invalid repo type.
func TestGetConfiguration_FailRepo(t *testing.T) {
	setEnv("PROD", "", "60", "3333", "", "")
	_, err := common.GetConfiguration()
	notOk(t, err)
}

// TestGetConfiguration_FailTimeout ensures that an error is returned when specifying an invalid timeout.
func TestGetConfiguration_FailTimeout(t *testing.T) {
	setEnv("PROD", "IN_MEMORY", "", "3333", "", "")
	_, err := common.GetConfiguration()
	notOk(t, err)
}

// TestGetConfiguration_FailPort ensures that an error is returned when specifying an invalid port.
func TestGetConfiguration_FailPort(t *testing.T) {
	setEnv("PROD", "IN_MEMORY", "60", "", "", "")
	_, err := common.GetConfiguration()
	notOk(t, err)
}

// TestGetConfiguration_PgSuccess ensures that a configuration is returned when specifying a PostgreSQL repo type.
func TestGetConfiguration_PgSuccess(t *testing.T) {
	setEnv("PROD", "POSTGRESQL", "60", "3333",
		"postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable", "")
	_, err := common.GetConfiguration()
	ok(t, err)
}

// TestGetConfiguration_PgFailPgUrl ensures that an error is returned when specifying a PostgreSQL repo type without a
// connection url.
func TestGetConfiguration_PgFailPgUrl(t *testing.T) {
	setEnv("PROD", "POSTGRESQL", "60", "3333", "", "")
	_, err := common.GetConfiguration()
	notOk(t, err)
}
