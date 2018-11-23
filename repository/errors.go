package repository

import "github.com/pkg/errors"

type errRepository struct {
	err error
}

func (er errRepository) Error() string {
	return er.err.Error()
}

func ErrRepository(msg string) error {
	return errRepository{errors.New(msg)}
}
