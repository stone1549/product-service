package common

import (
	"fmt"
	"github.com/pkg/errors"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	lifeCycleKey      string = "PRODUCT_SERVICE_ENVIRONMENT"
	repoTypeKey       string = "PRODUCT_SERVICE_REPO_TYPE"
	timeoutSecondsKey string = "PRODUCT_SERVICE_TIMEOUT"
	portKey           string = "PRODUCT_SERVICE_PORT"
	pgUrlKey          string = "PRODUCT_SERVICE_PG_URL"
	pgInitDatasetKey  string = "PRODUCT_SERVICE_INIT_DATASET"
)

type LifeCycle int

const (
	DevLifeCycle     LifeCycle = 0
	PreProdLifeCycle LifeCycle = iota
	ProdLifeCycle    LifeCycle = iota
)

func (lc LifeCycle) String() string {
	switch lc {
	case DevLifeCycle:
		return "DEV"
	case PreProdLifeCycle:
		return "PRE_PROD"
	case ProdLifeCycle:
		return "PROD"
	default:
		return "DEV"
	}
}

type ProductRepositoryType int

const (
	InMemoryRepo ProductRepositoryType = 0
	PostgreSQL   ProductRepositoryType = iota
)

type InitDataset int

const (
	NoDataset    InitDataset = 0
	SmallDataset InitDataset = iota
)

func (id InitDataset) String() string {
	switch id {
	case SmallDataset:
		return "SMALL"
	case NoDataset:
		return "NONE"
	default:
		return ""
	}
}

type Configuration interface {
	// Required config
	GetLifeCycle() LifeCycle
	GetRepoType() ProductRepositoryType
	GetTimeout() time.Duration
	GetPort() int

	// Optional config
	GetInitDataSet() InitDataset

	// PostgreSQL config
	GetPgUrl() string
}

type configuration struct {
	lifeCycle   LifeCycle
	repoType    ProductRepositoryType
	timeout     time.Duration
	port        int
	pgUrl       string
	initDataset InitDataset
}

func (conf *configuration) GetLifeCycle() LifeCycle {
	return conf.lifeCycle
}

func (conf *configuration) GetRepoType() ProductRepositoryType {
	return conf.repoType
}

func (conf *configuration) GetTimeout() time.Duration {
	return conf.timeout
}

func (conf *configuration) GetPort() int {
	return conf.port
}

func (conf *configuration) GetPgUrl() string {
	return conf.pgUrl
}

func (conf *configuration) GetInitDataSet() InitDataset {
	return conf.initDataset
}

func GetConfiguration() (Configuration, error) {
	var err error
	config := configuration{}

	lcStr := os.Getenv(lifeCycleKey)

	switch lcStr {
	case DevLifeCycle.String():
		config.lifeCycle = DevLifeCycle
	case PreProdLifeCycle.String():
		config.lifeCycle = PreProdLifeCycle
	case ProdLifeCycle.String():
		config.lifeCycle = ProdLifeCycle
	default:
		config.lifeCycle = DevLifeCycle
	}

	if err != nil {
		return nil, err
	}

	repoTypeStr := os.Getenv(repoTypeKey)

	switch repoTypeStr {
	case "IN_MEMORY":
		config.repoType = InMemoryRepo
	case "POSTGRESQL":
		config.repoType = PostgreSQL
	default:
		if config.lifeCycle == DevLifeCycle {
			config.repoType = InMemoryRepo
		} else {
			err = errors.New(fmt.Sprintf("No repo type configured, set %s environment variable", repoTypeKey))
		}
	}

	if err != nil {
		return nil, err
	}

	timeoutStr := os.Getenv(timeoutSecondsKey)

	timeoutInt, err := strconv.Atoi(timeoutStr)

	if config.lifeCycle == DevLifeCycle && err != nil {
		timeoutInt = 60
	} else if err != nil {
		err = errors.New(fmt.Sprintf("No timeout configured, set %s environment variable", timeoutSecondsKey))
		return nil, err
	}

	config.timeout = time.Duration(timeoutInt) * time.Second

	portStr := os.Getenv(portKey)
	port, err := strconv.Atoi(portStr)

	if config.lifeCycle == DevLifeCycle && err != nil {
		config.port = 3333
	} else if err != nil {
		err = errors.New(fmt.Sprintf("No port configured, set %s environment variable", portKey))
		return nil, err
	}

	config.port = port

	if config.repoType == PostgreSQL {
		setPostgresqlConfig(&config)
	}

	initDatasetStr := os.Getenv(pgInitDatasetKey)
	switch initDatasetStr {
	case NoDataset.String():
		config.initDataset = NoDataset
	case SmallDataset.String():
		config.initDataset = SmallDataset
	default:
		if initDatasetStr == "" {
			config.initDataset = NoDataset
		} else {
			err = errors.New(fmt.Sprintf("Invalid dataset, set %s environment variable properly or omit it", pgInitDatasetKey))
		}
	}

	if err != nil {
		return nil, err
	}

	return &config, nil
}

func setPostgresqlConfig(config *configuration) error {
	var err error

	config.pgUrl = os.Getenv(pgUrlKey)

	if strings.TrimSpace(config.pgUrl) == "" {
		err = errors.New(fmt.Sprintf("No PostgreSQL url configured, set %s environment variable", pgUrlKey))
	}

	return err
}
