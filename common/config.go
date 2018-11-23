package common

import (
	"fmt"
	"github.com/pkg/errors"
	"os"
	"strconv"
	"time"
)

const (
	LifeCycleKey      string = "PRODUCT_SERVICE_ENVIRONMENT"
	RepoTypeKey       string = "PRODUCT_REPO_TYPE"
	TimeoutSecondsKey string = "PRODUCT_SERVICE_TIMEOUT"
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
	GetLifeCycle() LifeCycle
	GetRepoType() ProductRepositoryType
	GetTimeout() time.Duration
}

type configuration struct {
	lifeCycle LifeCycle
	repoType  ProductRepositoryType
	timeout   time.Duration
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

func GetConfiguration() (Configuration, error) {
	var err error

	lcStr := os.Getenv(LifeCycleKey)
	var lifeCycle LifeCycle

	switch lcStr {
	case "DEV":
		lifeCycle = Dev
	case "PRE_PROD":
		lifeCycle = PreProd
	case "PROD":
		lifeCycle = Prod
	default:
		lifeCycle = Dev
	}

	if err != nil {
		return nil, err
	}

	repoTypeStr := os.Getenv(RepoTypeKey)
	var repoType ProductRepositoryType

	switch repoTypeStr {
	case "IN_MEMORY":
		repoType = InMemory
	case "POSTGRESQL":
		repoType = PostgreSQL
	default:
		if lifeCycle == Dev {
			repoType = InMemory
		} else {
			err = errors.New(fmt.Sprintf("No repo type configured, set %s environment variable", RepoTypeKey))
		}
	}

	if err != nil {
		return nil, err
	}

	timeoutStr := os.Getenv(TimeoutSecondsKey)

	timeoutInt, err := strconv.Atoi(timeoutStr)

	if lifeCycle == Dev && err != nil {
		timeoutInt = 60
	} else if err != nil {
		err = errors.New(fmt.Sprintf("No timeout configured, set %s environment variable", TimeoutSecondsKey))
		return nil, err
	}

	timeout := time.Duration(timeoutInt) * time.Second
	return &configuration{lifeCycle, repoType, timeout}, nil
}
