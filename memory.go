package main

import (
	"unsafe"
)

func memFor(buf []byte) (ptr, size) {
	return ptr(uintptr(unsafe.Pointer(&buf[0]))), size(len(buf))
}

func getBytes(nr node_ref) ([]byte, error) {
	// first try with a default buffer
	var smallBuf [2048]byte
	sp, ss := memFor(smallBuf[:])

	sz, err := as_bytes(nr, sp, ss)
	if err != 0 && sz > ss {
		// couldn't fit.
		alloc := make([]byte, sz)
		sp, ss = memFor(alloc)
		_, err = as_bytes(nr, sp, ss)
		return alloc, ErrCode{err}
	} else if err != 0 {
		return nil, ErrCode{err}
	}
	return smallBuf[0:sz], nil
}

func getString(nr node_ref) (string, error) {
	// first try with a default buffer
	var smallBuf [2048]byte
	sp, ss := memFor(smallBuf[:])

	sz, err := as_string(nr, sp, ss)
	if err != 0 && sz > ss {
		// couldn't fit.
		alloc := make([]byte, sz)
		sp, ss = memFor(alloc)
		_, err = as_string(nr, sp, ss)
		return string(alloc), ErrCode{err}
	} else if err != 0 {
		return "", ErrCode{err}
	}
	return string(smallBuf[0:sz]), nil
}

func getLink(nr node_ref) ([]byte, error) {
	// first try with a default buffer
	var smallBuf [2048]byte
	sp, ss := memFor(smallBuf[:])

	sz, err := as_link(nr, sp, ss)
	if err != 0 && sz > ss {
		// couldn't fit.
		alloc := make([]byte, sz)
		sp, ss = memFor(alloc)
		_, err = as_link(nr, sp, ss)
		return alloc, ErrCode{err}
	} else if err != 0 {
		return nil, ErrCode{err}
	}
	return smallBuf[0:sz], nil
}
