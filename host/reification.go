package wasmadl

import (
	"bytes"
	"io/ioutil"
	"reflect"

	"github.com/ipfs/go-cid"
	"github.com/willscott/wasm-adl/common"

	ipld "github.com/ipld/go-ipld-prime"
	cidlink "github.com/ipld/go-ipld-prime/linking/cid"
	basicnode "github.com/ipld/go-ipld-prime/node/basic"
	"github.com/willscott/gasm/hostfunc"
	wasm "github.com/willscott/gasm/wasm"
)

//wasmADL is an instance of an ADL that has been dynamically loaded
type wasmADL struct {
	*wasm.VirtualMachine
	linkSys *ipld.LinkSystem
	// not well named: these are node references in use by the adl
	functions map[int32]ipld.Node

	maprefs  map[int32]ipld.MapIterator
	listrefs map[int32]ipld.ListIterator
}

// newVM creates an ADL instance against a given module.
func newVM(adl string) (*wasmADL, error) {
	buf, err := ioutil.ReadFile(adl)
	if err != nil {
		return nil, err
	}
	mod, err := wasm.DecodeModule(bytes.NewBuffer(buf))
	if err != nil {
		return nil, err
	}

	wadl := &wasmADL{}
	wadl.functions = make(map[int32]ipld.Node)
	wadl.maprefs = make(map[int32]ipld.MapIterator)
	wadl.listrefs = make(map[int32]ipld.ListIterator)

	builder := hostfunc.NewModuleBuilder()
	builder.MustSetFunction("env", "load", wadl.load)

	builder.MustSetFunction("env", "kind", wadl.kind)
	builder.MustSetFunction("env", "lookup_by_string", wadl.lookup_by_string)
	builder.MustSetFunction("env", "lookup_by_node", wadl.lookup_by_node)
	builder.MustSetFunction("env", "lookup_by_index", wadl.lookup_by_index)
	builder.MustSetFunction("env", "lookup_by_string", wadl.lookup_by_string)
	builder.MustSetFunction("env", "map_iterator", wadl.map_iterator)
	builder.MustSetFunction("env", "list_iterator", wadl.list_iterator)
	builder.MustSetFunction("env", "length", wadl.length)
	builder.MustSetFunction("env", "is_absent", wadl.is_absent)
	builder.MustSetFunction("env", "is_null", wadl.is_null)
	builder.MustSetFunction("env", "as_bool", wadl.as_bool)
	builder.MustSetFunction("env", "as_int", wadl.as_int)
	builder.MustSetFunction("env", "as_float", wadl.as_float)
	builder.MustSetFunction("env", "as_string", wadl.as_string)
	builder.MustSetFunction("env", "as_bytes", wadl.as_bytes)
	builder.MustSetFunction("env", "as_link", wadl.as_link)
	builder.MustSetFunction("env", "map_iterate_next", wadl.map_iterate_next)
	builder.MustSetFunction("env", "map_iterate_done", wadl.map_iterate_done)
	builder.MustSetFunction("env", "list_iterate_next", wadl.list_iterate_next)
	builder.MustSetFunction("env", "list_iterate_done", wadl.list_iterate_done)

	builder.MustSetFunction("env", "new_null", wadl.new_null)
	builder.MustSetFunction("env", "new_bool", wadl.new_bool)
	builder.MustSetFunction("env", "new_int", wadl.new_int)
	builder.MustSetFunction("env", "new_float", wadl.new_float)
	builder.MustSetFunction("env", "new_string", wadl.new_string)
	builder.MustSetFunction("env", "new_bytes", wadl.new_bytes)
	builder.MustSetFunction("env", "new_link", wadl.new_link)
	builder.MustSetFunction("env", "new_map", wadl.new_map)
	builder.MustSetFunction("env", "new_list", wadl.new_list)
	builder.MustSetFunction("env", "assemble_map", wadl.assemble_map)
	builder.MustSetFunction("env", "finish_map", wadl.finish_map)
	builder.MustSetFunction("env", "assemble_list", wadl.assemble_list)
	builder.MustSetFunction("env", "finish_list", wadl.finish_list)

	vm, err := wasm.NewVM(mod, builder.Done())
	if err != nil {
		return nil, err
	}
	wadl.VirtualMachine = vm

	return wadl, nil
}

func (w *wasmADL) load(m *wasm.VirtualMachine) reflect.Value {
	body := func(p common.Ptr, s common.Size) common.Node_ref {
		linkBytes := w.getMem(uint64(p), uint64(s))
		if linkBytes == nil {
			return 0
		}
		_, lnk, err := cid.CidFromBytes(linkBytes)
		if err != nil {
			return 0
		}
		node, err := w.linkSys.Load(ipld.LinkContext{}, cidlink.Link{lnk}, basicnode.Prototype.Any)
		if err != nil {
			return 0
		}
		return w.alloc(node)
	}
	return reflect.ValueOf(body)
}

func (w *wasmADL) alloc(n ipld.Node) common.Node_ref {
	k := len(w.functions)
	kn := int32(k + 1)
	w.functions[kn] = n
	return common.Node_ref(kn)
}

func (w *wasmADL) getMem(ptr uint64, size uint64) []byte {
	if int(ptr) > len(w.VirtualMachine.Memory) {
		return nil
	}
	return w.VirtualMachine.Memory[ptr : ptr+size]
}

func (w *wasmADL) getMemPS(ptr common.Ptr, size common.Size) []byte {
	return w.getMem(uint64(ptr), uint64(size))
}

func (w *wasmADL) setMem(src []byte, at common.Ptr) int {
	return copy(w.VirtualMachine.Memory[uint64(at):], src)
}

func New(adl string) *adlBuilder {
	return &adlBuilder{adl}
}

type adlBuilder struct {
	adl string
}

// Reify will reify against known wasm modules
func (a *adlBuilder) Reify(_ ipld.LinkContext, n ipld.Node, ls *ipld.LinkSystem) (ipld.Node, error) {
	vm, err := newVM(a.adl)
	if err != nil {
		return nil, err
	}
	vm.functions[1] = n
	vm.linkSys = ls
	rets, _, err := vm.ExecExportedFunction("adl_reify", 1)
	if err != nil {
		return nil, err
	}
	if rets[0] != 0 {
		return nil, common.ErrCode{common.Errno(rets[0])}
	}
	return &proxyNode{vm}, nil
}
