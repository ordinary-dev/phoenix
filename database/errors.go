package database

import (
	"errors"
)

var (
	ErrWrongNumberOfAffectedRows = errors.New("wrong number of affected rows")
)
