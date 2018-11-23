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
	LifeCycleKey      string = "PRODUCT_SERVICE_ENVIRONMENT"
	RepoTypeKey       string = "PRODUCT_SERVICE_REPO_TYPE"
	TimeoutSecondsKey string = "PRODUCT_SERVICE_TIMEOUT"
	PortKey           string = "PRODUCT_SERVICE_PORT"
	PgUrlKey          string = "PRODUCT_SERVICE_PG_URL"
)

type LifeCycle int

const (
	Dev     LifeCycle = 0
	PreProd LifeCycle = iota
	Prod    LifeCycle = iota
)

type ProductRepositoryType int

const (
	InMemory   ProductRepositoryType = 0
	PostgreSQL ProductRepositoryType = iota
)

type Configuration interface {
	// Required config
	GetLifeCycle() LifeCycle
	GetRepoType() ProductRepositoryType
	GetTimeout() time.Duration
	GetPort() int

	// PostgreSQL config
	GetPgUrl() string
}

type configuration struct {
	lifeCycle LifeCycle
	repoType  ProductRepositoryType
	timeout   time.Duration
	port      int
	pgUrl     string
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

func GetConfiguration() (Configuration, error) {
	var err error
	config := configuration{}

	lcStr := os.Getenv(LifeCycleKey)

	switch lcStr {
	case "DEV":
		config.lifeCycle = Dev
	case "PRE_PROD":
		config.lifeCycle = PreProd
	case "PROD":
		config.lifeCycle = Prod
	default:
		config.lifeCycle = Dev
	}

	if err != nil {
		return nil, err
	}

	repoTypeStr := os.Getenv(RepoTypeKey)

	switch repoTypeStr {
	case "IN_MEMORY":
		config.repoType = InMemory
	case "POSTGRESQL":
		config.repoType = PostgreSQL
	default:
		if config.lifeCycle == Dev {
			config.repoType = InMemory
		} else {
			err = errors.New(fmt.Sprintf("No repo type configured, set %s environment variable", RepoTypeKey))
		}
	}

	if err != nil {
		return nil, err
	}

	timeoutStr := os.Getenv(TimeoutSecondsKey)

	timeoutInt, err := strconv.Atoi(timeoutStr)

	if config.lifeCycle == Dev && err != nil {
		timeoutInt = 60
	} else if err != nil {
		err = errors.New(fmt.Sprintf("No timeout configured, set %s environment variable", TimeoutSecondsKey))
		return nil, err
	}

	config.timeout = time.Duration(timeoutInt) * time.Second

	portStr := os.Getenv(PortKey)
	port, err := strconv.Atoi(portStr)

	if config.lifeCycle == Dev && err != nil {
		config.port = 3333
	} else if err != nil {
		err = errors.New(fmt.Sprintf("No port configured, set %s environment variable", PortKey))
		return nil, err
	}

	config.port = port

	if config.repoType == PostgreSQL {
		setPostgresqlConfig(&config)
	}

	return &config, nil
}

func setPostgresqlConfig(config *configuration) error {
	var err error

	config.pgUrl = os.Getenv(PgUrlKey)

	if strings.TrimSpace(config.pgUrl) == "" {
		err = errors.New(fmt.Sprintf("No PostgreSQL url configured, set %s environment variable", PgUrlKey))
	}

	return err
}
