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

// LifeCycle represents a particular application life cycle.
type LifeCycle int

const (
	// DevLifeCycle represents the development environment.
	DevLifeCycle LifeCycle = 0
	// PreProdLifeCycle represents the pre production environment.
	PreProdLifeCycle LifeCycle = iota
	// ProdLifeCycle represents the production environment.
	ProdLifeCycle LifeCycle = iota
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

// ProductRepositoryType represents a type of ProductRepository
type ProductRepositoryType int

const (
	// InMemoryRepo represents a ProductRepository that is stored entirely in memory.
	InMemoryRepo ProductRepositoryType = 0
	// PostgreSqlRepo represents a ProductRepository that utilizes a PostgreSQL database.
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

// Configuration provides methods for retrieving aspects of the applications configuration.
type Configuration interface {
	// GetLifeCycle retrieves the configured life cycle.
	GetLifeCycle() LifeCycle
	// GetRepoType retrieves the configured repo type.
	GetRepoType() ProductRepositoryType
	// GetTimeout retrieves the configured request timeout.
	GetTimeout() time.Duration
	// GetPort retrieves the configured port.
	GetPort() int

	// GetInitDataSet retrieves the path to an initial dataset to load on app launch, mostly for testing and dev use.
	GetInitDataSet() string

	// GetPgUrl retrieves the configured url string for connecting to PostgreSQL.
	GetPgUrl() string
}

type configuration struct {
	lifeCycle   LifeCycle
	repoType    ProductRepositoryType
	timeout     time.Duration
	port        int
	pgUrl       string
	initDataset string
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

func (conf *configuration) GetInitDataSet() string {
	return conf.initDataset
}

// GetConfiguration constucts a Configuration based on environment variables.
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

	config.initDataset = os.Getenv(initDatasetKey)

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
