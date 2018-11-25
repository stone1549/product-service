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
	initDatasetKey    string = "PRODUCT_SERVICE_INIT_DATASET"
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
		return ""
	}
}

type ProductRepositoryType int

const (
	InMemoryRepo   ProductRepositoryType = 0
	PostgreSqlRepo ProductRepositoryType = iota
)

func (prt ProductRepositoryType) String() string {
	switch prt {
	case PostgreSqlRepo:
		return "POSTGRESQL"
	case InMemoryRepo:
		return "IN_MEMORY"
	default:
		return ""
	}
}

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

	// PostgreSqlRepo config
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
	case InMemoryRepo.String():
		config.repoType = InMemoryRepo
	case PostgreSqlRepo.String():
		config.repoType = PostgreSqlRepo
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

	if timeoutStr == "" && config.lifeCycle == DevLifeCycle {
		timeoutStr = "60"
	}

	timeoutInt, err := strconv.Atoi(timeoutStr)

	if err != nil {
		err = errors.New(fmt.Sprintf("No timeout configured, set %s environment variable", timeoutSecondsKey))
		return nil, err
	}

	config.timeout = time.Duration(timeoutInt) * time.Second

	portStr := os.Getenv(portKey)

	if portStr == "" && config.lifeCycle == DevLifeCycle {
		portStr = "3333"
	}
	port, err := strconv.Atoi(portStr)

	if err != nil {
		err = errors.New(fmt.Sprintf("No port configured, set %s environment variable", portKey))
		return nil, err
	}

	config.port = port

	if config.repoType == PostgreSqlRepo {
		err = setPostgresqlConfig(&config)
	}

	if err != nil {
		return nil, err
	}

	initDatasetStr := os.Getenv(initDatasetKey)
	switch initDatasetStr {
	case NoDataset.String():
		config.initDataset = NoDataset
	case SmallDataset.String():
		config.initDataset = SmallDataset
	default:
		if initDatasetStr == "" {
			config.initDataset = NoDataset
		} else {
			err = errors.New(fmt.Sprintf("Invalid dataset, set %s environment variable properly or omit it", initDatasetKey))
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
		err = errors.New(fmt.Sprintf("No PostgreSqlRepo url configured, set %s environment variable", pgUrlKey))
	}

	return err
}
