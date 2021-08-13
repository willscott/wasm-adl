package wasmadl

import (
	"reflect"

	"github.com/ipfs/go-cid"
	cidlink "github.com/ipld/go-ipld-prime/linking/cid"
	basicnode "github.com/ipld/go-ipld-prime/node/basic"
	wasm "github.com/willscott/gasm/wasm"
	"github.com/willscott/wasm-adl/common"
)

//new_null() node_ref
func (w *wasmADL) new_null(m *wasm.VirtualMachine) reflect.Value {
	body := func() common.Node_ref {
		b := basicnode.Prototype.Any.NewBuilder()
		_ = b.AssignNull()
		n := b.Build()
		return w.alloc(n)
	}
	return reflect.ValueOf(body)
}

//new_bool(uint32) node_ref
func (w *wasmADL) new_bool(m *wasm.VirtualMachine) reflect.Value {
	body := func(v uint32) common.Node_ref {
		if v != 0 {
			b := basicnode.NewBool(true)
			return w.alloc(b)
		}
		b := basicnode.NewBool(false)
		return w.alloc(b)
	}
	return reflect.ValueOf(body)
}

//new_int(uint64) node_ref
func (w *wasmADL) new_int(m *wasm.VirtualMachine) reflect.Value {
	body := func(v uint64) common.Node_ref {
		i := basicnode.NewInt(int64(v))
		return w.alloc(i)
	}
	return reflect.ValueOf(body)
}

//new_float(uint64) node_ref
func (w *wasmADL) new_float(m *wasm.VirtualMachine) reflect.Value {
	body := func(v uint64) common.Node_ref {
		f := basicnode.NewFloat(float64(v))
		return w.alloc(f)
	}
	return reflect.ValueOf(body)
}

//new_string(ptr, size) node_ref
func (w *wasmADL) new_string(m *wasm.VirtualMachine) reflect.Value {
	body := func(ptr common.Ptr, size common.Size) common.Node_ref {
		buf := w.getMemPS(ptr, size)
		if buf == nil {
			return 0
		}
		s := basicnode.NewString(string(buf))
		return w.alloc(s)
	}
	return reflect.ValueOf(body)
}

//new_bytes(ptr, size) node_ref
func (w *wasmADL) new_bytes(m *wasm.VirtualMachine) reflect.Value {
	body := func(ptr common.Ptr, size common.Size) common.Node_ref {
		buf := w.getMemPS(ptr, size)
		if buf == nil {
			return 0
		}
		b := basicnode.NewBytes(buf)
		return w.alloc(b)
	}
	return reflect.ValueOf(body)
}

//new_link(ptr, size) node_ref
func (w *wasmADL) new_link(m *wasm.VirtualMachine) reflect.Value {
	body := func(ptr common.Ptr, size common.Size) common.Node_ref {
		buf := w.getMemPS(ptr, size)
		if buf == nil {
			return 0
		}
		_, c, err := cid.CidFromBytes(buf)
		if err != nil {
			return 0
		}
		return w.alloc(basicnode.NewLink(cidlink.Link{Cid: c}))
	}
	return reflect.ValueOf(body)
}

//new_map(size) map_builder_ref
func (w *wasmADL) new_map(m *wasm.VirtualMachine) reflect.Value {
	return reflect.ValueOf(0)
}

//new_list(size) list_builder_ref
func (w *wasmADL) new_list(m *wasm.VirtualMachine) reflect.Value {
	return reflect.ValueOf(0)
}

//assemble_map(map_builder_ref, node_ref, node_ref) errno
func (w *wasmADL) assemble_map(m *wasm.VirtualMachine) reflect.Value {
	return reflect.ValueOf(0)
}

//finish_map(map_builder_ref) node_ref
func (w *wasmADL) finish_map(m *wasm.VirtualMachine) reflect.Value {
	return reflect.ValueOf(0)
}

//assemble_list(list_builder_ref, node_ref) errno
func (w *wasmADL) assemble_list(m *wasm.VirtualMachine) reflect.Value {
	return reflect.ValueOf(0)
}

//finish_list(list_builder_ref) node_ref
func (w *wasmADL) finish_list(m *wasm.VirtualMachine) reflect.Value {
	return reflect.ValueOf(0)
}
