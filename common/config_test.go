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

func TestGetConfiguration_Defaults(t *testing.T) {
	clearEnv()
	_, err := common.GetConfiguration()
	ok(t, err)
}

func TestGetConfiguration_ImSuccess(t *testing.T) {
	setEnv("DEV", "IN_MEMORY", "60", "3333", "", "NONE")
	_, err := common.GetConfiguration()
	ok(t, err)
}

func TestGetConfiguration_ImSuccessSmallDataset(t *testing.T) {
	setEnv("DEV", "IN_MEMORY", "60", "3333", "", "SMALL")
	_, err := common.GetConfiguration()
	ok(t, err)
}

func TestGetConfiguration_ImSuccessNoneDataset(t *testing.T) {
	setEnv("DEV", "IN_MEMORY", "60", "3333", "", "NONE")
	_, err := common.GetConfiguration()
	ok(t, err)
}

func TestGetConfiguration_FailRepo(t *testing.T) {
	setEnv("PROD", "", "60", "3333", "", "NONE")
	_, err := common.GetConfiguration()
	notOk(t, err)
}

func TestGetConfiguration_FailTimeout(t *testing.T) {
	setEnv("PROD", "IN_MEMORY", "", "3333", "", "NONE")
	_, err := common.GetConfiguration()
	notOk(t, err)
}

func TestGetConfiguration_FailPort(t *testing.T) {
	setEnv("PROD", "IN_MEMORY", "60", "", "", "NONE")
	_, err := common.GetConfiguration()
	notOk(t, err)
}

func TestGetConfiguration_FailDatAset(t *testing.T) {
	setEnv("PROD", "IN_MEMORY", "60", "3333", "", "HUGE")
	_, err := common.GetConfiguration()
	notOk(t, err)
}

func TestGetConfiguration_PgSuccess(t *testing.T) {
	setEnv("PROD", "POSTGRESQL", "60", "3333",
		"postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable", "NONE")
	_, err := common.GetConfiguration()
	ok(t, err)
}

func TestGetConfiguration_PgFailPgUrl(t *testing.T) {
	setEnv("PROD", "POSTGRESQL", "60", "3333", "", "NONE")
	_, err := common.GetConfiguration()
	notOk(t, err)
}
