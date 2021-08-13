package wasmadl

import (
	"errors"
	"reflect"

	"github.com/ipld/go-ipld-prime"
	cidlink "github.com/ipld/go-ipld-prime/linking/cid"
	wasm "github.com/willscott/gasm/wasm"
	"github.com/willscott/wasm-adl/common"
)

//(node_ref) uint32
func (w *wasmADL) kind(m *wasm.VirtualMachine) reflect.Value {
	body := func(nr common.Node_ref) uint32 {
		n, ok := w.functions[int32(nr)]
		if !ok {
			return uint8(ipld.Kind_Invalid)
		}
		return uint32(n.Kind())
	}
	return reflect.ValueOf(body)
}

//(node_ref, ptr, size) (node_ref, errno)
func (w *wasmADL) lookup_by_string(m *wasm.VirtualMachine) reflect.Value {
	body := func(nr common.Node_ref, p common.Ptr, s common.Size) (common.Node_ref, common.Errno) {
		n, ok := w.functions[int32(nr)]
		if !ok {
			return 0, common.ErrUnknown
		}
		buf := w.getMemPS(p, s)
		if buf == nil {
			return 0, common.ErrUnknown
		}
		val, err := n.LookupByString(string(buf))
		if err != nil {
			if errors.As(err, ipld.ErrNotExists{}) {
				return 0, common.ErrDoesNotExist
			}
			if errors.As(err, ipld.ErrWrongKind{}) {
				return 0, common.ErrWrongKind
			}
			return 0, common.ErrUnknown
		}
		nr = w.alloc(val)
		return nr, 0
	}
	return reflect.ValueOf(body)
}

//(node_ref, node_ref) (node_ref, errno)
func (w *wasmADL) lookup_by_node(m *wasm.VirtualMachine) reflect.Value {
	body := func(nr common.Node_ref, pathref common.Node_ref) (common.Node_ref, common.Errno) {
		n, ok := w.functions[int32(nr)]
		if !ok {
			return 0, common.ErrUnknown
		}
		key, ok := w.functions[int32(pathref)]
		if !ok {
			return 0, common.ErrUnknown
		}
		val, err := n.LookupByNode(key)
		if err != nil {
			if errors.As(err, ipld.ErrNotExists{}) {
				return 0, common.ErrDoesNotExist
			}
			if errors.As(err, ipld.ErrWrongKind{}) {
				return 0, common.ErrWrongKind
			}
			return 0, common.ErrUnknown
		}
		nr = w.alloc(val)
		return nr, 0
	}
	return reflect.ValueOf(body)
}

//(node_ref, uint64) (node_ref, errno)
func (w *wasmADL) lookup_by_index(m *wasm.VirtualMachine) reflect.Value {
	body := func(nr common.Node_ref, idx uint64) (common.Node_ref, common.Errno) {
		n, ok := w.functions[int32(nr)]
		if !ok {
			return 0, common.ErrUnknown
		}
		val, err := n.LookupByIndex(int64(idx))
		if err != nil {
			if errors.As(err, ipld.ErrNotExists{}) {
				return 0, common.ErrDoesNotExist
			}
			if errors.As(err, ipld.ErrWrongKind{}) {
				return 0, common.ErrWrongKind
			}
			return 0, common.ErrUnknown
		}
		nr = w.alloc(val)
		return nr, 0
	}
	return reflect.ValueOf(body)
}

//map_iterator(node_ref) map_ref
func (w *wasmADL) map_iterator(m *wasm.VirtualMachine) reflect.Value {
	return reflect.ValueOf(0)
}

//list_iterator(node_ref) list_ref
func (w *wasmADL) list_iterator(m *wasm.VirtualMachine) reflect.Value {
	return reflect.ValueOf(0)
}

//length(node_ref) uint64
func (w *wasmADL) length(m *wasm.VirtualMachine) reflect.Value {
	body := func(nr common.Node_ref) uint64 {
		n, ok := w.functions[int32(nr)]
		if !ok {
			return 0
		}

		return uint64(n.Length())
	}
	return reflect.ValueOf(body)
}

//is_absent(node_ref) uint8
func (w *wasmADL) is_absent(m *wasm.VirtualMachine) reflect.Value {
	body := func(nr common.Node_ref) uint32 {
		n, ok := w.functions[int32(nr)]
		if !ok {
			return 2
		}

		if n.IsAbsent() {
			return 1
		}
		return 0
	}
	return reflect.ValueOf(body)
}

//is_null(node_ref) uint8
func (w *wasmADL) is_null(m *wasm.VirtualMachine) reflect.Value {
	body := func(nr common.Node_ref) uint32 {
		n, ok := w.functions[int32(nr)]
		if !ok {
			return 2
		}

		if n.IsNull() {
			return 1
		}
		return 0
	}
	return reflect.ValueOf(body)
}

//as_bool(node_ref) (uint8, errno)
func (w *wasmADL) as_bool(m *wasm.VirtualMachine) reflect.Value {
	body := func(nr common.Node_ref) (uint32, common.Errno) {
		n, ok := w.functions[int32(nr)]
		if !ok {
			return 0, common.ErrUnknown
		}

		b, err := n.AsBool()
		if err != nil {
			if errors.As(err, ipld.ErrNotExists{}) {
				return 0, common.ErrDoesNotExist
			}
			if errors.As(err, ipld.ErrWrongKind{}) {
				return 0, common.ErrWrongKind
			}
			return 0, common.ErrUnknown
		}
		if b {
			return 1, 0
		}
		return 0, 0
	}
	return reflect.ValueOf(body)
}

//as_int(node_ref) (uint64, errno)
func (w *wasmADL) as_int(m *wasm.VirtualMachine) reflect.Value {
	body := func(nr common.Node_ref) (uint64, common.Errno) {
		n, ok := w.functions[int32(nr)]
		if !ok {
			return 0, common.ErrUnknown
		}

		i, err := n.AsInt()
		if err != nil {
			if errors.As(err, ipld.ErrNotExists{}) {
				return 0, common.ErrDoesNotExist
			}
			if errors.As(err, ipld.ErrWrongKind{}) {
				return 0, common.ErrWrongKind
			}
			return 0, common.ErrUnknown
		}

		return uint64(i), 0
	}
	return reflect.ValueOf(body)
}

//as_float(node_ref) (uint64, errno)
func (w *wasmADL) as_float(m *wasm.VirtualMachine) reflect.Value {
	body := func(nr common.Node_ref) (uint64, common.Errno) {
		n, ok := w.functions[int32(nr)]
		if !ok {
			return 0, common.ErrUnknown
		}

		f, err := n.AsFloat()
		if err != nil {
			if errors.As(err, ipld.ErrNotExists{}) {
				return 0, common.ErrDoesNotExist
			}
			if errors.As(err, ipld.ErrWrongKind{}) {
				return 0, common.ErrWrongKind
			}
			return 0, common.ErrUnknown
		}

		return uint64(f), 0
	}
	return reflect.ValueOf(body)
}

//as_string(node_ref, ptr, size) (size, errno)
func (w *wasmADL) as_string(m *wasm.VirtualMachine) reflect.Value {
	body := func(nr common.Node_ref, p common.Ptr, siz common.Size) (common.Size, common.Errno) {
		n, ok := w.functions[int32(nr)]
		if !ok {
			return 0, common.ErrUnknown
		}

		s, err := n.AsString()
		if err != nil {
			if errors.As(err, ipld.ErrNotExists{}) {
				return 0, common.ErrDoesNotExist
			}
			if errors.As(err, ipld.ErrWrongKind{}) {
				return 0, common.ErrWrongKind
			}
			return 0, common.ErrUnknown
		}

		if common.Size(len(s)) > siz {
			return common.Size(len(s)), common.ErrTooBig
		}
		return common.Size(w.setMem([]byte(s), p)), 0
	}
	return reflect.ValueOf(body)
}

//as_bytes(node_ref, ptr, size) (size, errno)
func (w *wasmADL) as_bytes(m *wasm.VirtualMachine) reflect.Value {
	body := func(nr common.Node_ref, p common.Ptr, siz common.Size) (common.Size, common.Errno) {
		n, ok := w.functions[int32(nr)]
		if !ok {
			return 0, common.ErrUnknown
		}

		b, err := n.AsBytes()
		if err != nil {
			if errors.As(err, ipld.ErrNotExists{}) {
				return 0, common.ErrDoesNotExist
			}
			if errors.As(err, ipld.ErrWrongKind{}) {
				return 0, common.ErrWrongKind
			}
			return 0, common.ErrUnknown
		}

		if common.Size(len(b)) > siz {
			return common.Size(len(b)), common.ErrTooBig
		}
		return common.Size(w.setMem(b, p)), 0
	}
	return reflect.ValueOf(body)
}

//as_link(node_ref, ptr, size) (size, errno)
func (w *wasmADL) as_link(m *wasm.VirtualMachine) reflect.Value {
	body := func(nr common.Node_ref, p common.Ptr, siz common.Size) (common.Size, common.Errno) {
		n, ok := w.functions[int32(nr)]
		if !ok {
			return 0, common.ErrUnknown
		}

		l, err := n.AsLink()
		if err != nil {
			if errors.As(err, ipld.ErrNotExists{}) {
				return 0, common.ErrDoesNotExist
			}
			if errors.As(err, ipld.ErrWrongKind{}) {
				return 0, common.ErrWrongKind
			}
			return 0, common.ErrUnknown
		}

		var lb []byte
		if cl, ok := l.(cidlink.Link); ok {
			lb = cl.Cid.Bytes()
		} else {
			lb = []byte(l.String())
		}

		if common.Size(len(lb)) > siz {
			return common.Size(len(lb)), common.ErrTooBig
		}
		return common.Size(w.setMem(lb, p)), 0
	}
	return reflect.ValueOf(body)
}

//map_iterate_next(map_ref) (node_ref, node_ref, errno)
func (w *wasmADL) map_iterate_next(m *wasm.VirtualMachine) reflect.Value {
	return reflect.ValueOf(0)
}

//map_iterate_done(map_ref) uint8
func (w *wasmADL) map_iterate_done(m *wasm.VirtualMachine) reflect.Value {
	return reflect.ValueOf(0)
}

//list_iterate_next(list_ref) (uint64, node_ref, errno)
func (w *wasmADL) list_iterate_next(m *wasm.VirtualMachine) reflect.Value {
	return reflect.ValueOf(0)
}

//export list_iterate_done(list_ref) uint8
func (w *wasmADL) list_iterate_done(m *wasm.VirtualMachine) reflect.Value {
	return reflect.ValueOf(0)
}
