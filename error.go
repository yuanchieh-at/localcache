package localcache

import (
	"errors"
	"fmt"
)

const (
	KeyNotFound = "KeyNotFound"
)

type Error struct {
	err error
	k string
	code string
}

func (l *Error) Error() string {
	return l.err.Error()
}

func NewKeyNotFound(k string) error {
	err := errors.New(fmt.Sprintf("key %s is not found", k))
	return &Error{
		k: k,
		err: err,
		code: KeyNotFound,
	}
}