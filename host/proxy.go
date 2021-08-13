package wasmadl

import (
	"fmt"

	"github.com/ipfs/go-cid"
	"github.com/ipld/go-ipld-prime"
	cidlink "github.com/ipld/go-ipld-prime/linking/cid"
	basicnode "github.com/ipld/go-ipld-prime/node/basic"
	"github.com/willscott/wasm-adl/common"
)

var (
	ErrInvalidWasmResponse = fmt.Errorf("Invalid WASM Response")
	ErrUnknownNodeReturned = fmt.Errorf("Unregistered Node Returned by WASM")
)

type proxyNode struct {
	*wasmADL
}

func (p *proxyNode) Kind() ipld.Kind {
	rets, _, err := p.ExecExportedFunction("adl_kind")
	if err != nil {
		return ipld.Kind_Invalid
	}
	return ipld.Kind(rets[0])
}

func (p *proxyNode) LookupByString(key string) (ipld.Node, error) {
	n := basicnode.NewString(key)
	nr := p.alloc(n)
	rets, _, err := p.ExecExportedFunction("adl_lookup_by_node", uint64(nr))
	if err != nil {
		return nil, err
	}
	if len(rets) < 2 {
		return nil, ErrInvalidWasmResponse
	}
	if rets[1] != 0 {
		return nil, common.ErrCode{common.Errno(rets[1])}
	}
	n, ok := p.functions[int32(rets[0])]
	if !ok {
		return nil, ErrUnknownNodeReturned
	}
	return n, nil
}

func (p *proxyNode) LookupByNode(key ipld.Node) (ipld.Node, error) {
	nr := p.alloc(key)
	rets, _, err := p.ExecExportedFunction("adl_lookup_by_node", uint64(nr))
	if err != nil {
		return nil, err
	}
	if len(rets) < 2 {
		return nil, ErrInvalidWasmResponse
	}
	if rets[1] != 0 {
		return nil, common.ErrCode{common.Errno(rets[1])}
	}
	n, ok := p.functions[int32(rets[0])]
	if !ok {
		return nil, ErrUnknownNodeReturned
	}
	return n, nil
}

func (p *proxyNode) LookupByIndex(idx int64) (ipld.Node, error) {
	rets, _, err := p.ExecExportedFunction("adl_lookup_by_index", uint64(idx))
	if err != nil {
		return nil, err
	}
	if len(rets) < 2 {
		return nil, ErrInvalidWasmResponse
	}
	if rets[1] != 0 {
		return nil, common.ErrCode{common.Errno(rets[1])}
	}
	n, ok := p.functions[int32(rets[0])]
	if !ok {
		return nil, ErrUnknownNodeReturned
	}
	return n, nil
}

func (p *proxyNode) LookupBySegment(seg ipld.PathSegment) (ipld.Node, error) {
	if i, err := seg.Index(); err == nil {
		return p.LookupByIndex(i)
	}
	return p.LookupByString(seg.String())
}

func (p *proxyNode) MapIterator() ipld.MapIterator {
	//todo
	return nil
}

func (p *proxyNode) ListIterator() ipld.ListIterator {
	//todo
	return nil
}

func (p *proxyNode) Length() int64 {
	rets, _, err := p.ExecExportedFunction("adl_length")
	if err != nil {
		return -1
	}
	if len(rets) < 1 {
		return -1
	}
	return int64(rets[0])
}

func (p *proxyNode) IsAbsent() bool {
	rets, _, err := p.ExecExportedFunction("adl_is_absent")
	if err != nil {
		return false
	}
	if len(rets) < 1 {
		return false
	}
	if rets[0] != 0 {
		return true
	}
	return false
}

func (p *proxyNode) IsNull() bool {
	rets, _, err := p.ExecExportedFunction("adl_is_null")
	if err != nil {
		return false
	}
	if len(rets) < 1 {
		return false
	}
	if rets[0] != 0 {
		return true
	}
	return false
}

func (p *proxyNode) AsBool() (bool, error) {
	rets, _, err := p.ExecExportedFunction("adl_as_bool")
	if err != nil {
		return false, err
	}
	if len(rets) < 2 {
		return false, ErrInvalidWasmResponse
	}
	if rets[1] != 0 {
		return false, common.ErrCode{common.Errno(rets[1])}
	}
	if rets[0] != 0 {
		return true, nil
	}
	return false, nil
}

func (p *proxyNode) AsInt() (int64, error) {
	rets, _, err := p.ExecExportedFunction("adl_as_int")
	if err != nil {
		return 0, err
	}
	if len(rets) < 2 {
		return 0, ErrInvalidWasmResponse
	}
	if rets[1] != 0 {
		return 0, common.ErrCode{common.Errno(rets[1])}
	}
	return int64(rets[0]), nil
}

func (p *proxyNode) AsFloat() (float64, error) {
	rets, _, err := p.ExecExportedFunction("adl_as_float")
	if err != nil {
		return 0, err
	}
	if len(rets) < 2 {
		return 0, ErrInvalidWasmResponse
	}
	if rets[1] != 0 {
		return 0, common.ErrCode{common.Errno(rets[1])}
	}
	return float64(rets[0]), nil
}

func (p *proxyNode) AsString() (string, error) {
	rets, _, err := p.ExecExportedFunction("adl_as_string")
	if err != nil {
		return "", err
	}
	if len(rets) < 3 {
		return "", ErrInvalidWasmResponse
	}
	if rets[2] != 0 {
		return "", common.ErrCode{common.Errno(rets[1])}
	}

	b := p.getMem(rets[0], rets[1])
	return string(b), nil
}

func (p *proxyNode) AsBytes() ([]byte, error) {
	rets, _, err := p.ExecExportedFunction("adl_as_bytes")
	if err != nil {
		return nil, err
	}
	if len(rets) < 3 {
		return nil, ErrInvalidWasmResponse
	}
	if rets[2] != 0 {
		return nil, common.ErrCode{common.Errno(rets[1])}
	}

	b := p.getMem(rets[0], rets[1])
	return b, nil
}

func (p *proxyNode) AsLink() (ipld.Link, error) {
	rets, _, err := p.ExecExportedFunction("adl_as_link")
	if err != nil {
		return nil, err
	}
	if len(rets) < 3 {
		return nil, ErrInvalidWasmResponse
	}
	if rets[2] != 0 {
		return nil, common.ErrCode{common.Errno(rets[1])}
	}

	b := p.getMem(rets[0], rets[1])
	_, c, err := cid.CidFromBytes(b)
	if err != nil {
		return nil, err
	}
	return cidlink.Link{c}, nil
}

func (p *proxyNode) Prototype() ipld.NodePrototype {
	return nil
}
