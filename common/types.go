package common

import "fmt"

type Ptr uint64
type Size uint64
type Node_ref uint32
type Map_ref uint32
type List_ref uint32
type Map_builder_ref uint32
type List_builder_ref uint32
type Errno uint8

const (
	ErrUnknown       = 1
	ErrAlreadyExists = 2
	ErrDoesNotExist  = 3
	ErrWrongKind     = 4
	ErrTooBig        = 5
)

type ErrCode struct {
	Errno
}

func (e ErrCode) Error() string {
	return fmt.Sprintf("error: %d", e.Errno)
}
