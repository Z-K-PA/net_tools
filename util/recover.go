package util

import (
	"github.com/go-errors/errors"
)

func Recover(err interface{}) *errors.Error {
	if err != nil {
		return errors.Wrap(err, 3)
	}
	return nil
}
