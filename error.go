package main

import (
	"errors"
	"fmt"

	"github.com/ipld/go-ipld-prime"
	"github.com/willscott/wasm-adl/common"
)

type errno = common.Errno

const (
	ErrUnknown       = 1
	ErrAlreadyExists = 2
	ErrDoesNotExist  = 3
	ErrWrongKind     = 4
)

// ErrCode provides a golang error wrapper for numeric errors
type ErrCode struct {
	errno
}

func (e ErrCode) Error() string {
	return fmt.Sprintf("error: %d", e.errno)
}

func toErrno(e error) errno {
	if ec, ok := e.(ErrCode); ok {
		return ec.errno
	}
	if errors.Is(e, ipld.ErrWrongKind{}) {
		return ErrWrongKind
	}
	if errors.Is(e, ipld.ErrNotExists{}) {
		return ErrDoesNotExist
	}
	return ErrUnknown
}
