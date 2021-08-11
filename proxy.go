package main

import (
	"github.com/ipfs/go-cid"
	"github.com/ipld/go-ipld-prime"
	cidlink "github.com/ipld/go-ipld-prime/linking/cid"
)

type proxyNode struct {
	node_ref
}

func (p *proxyNode) Kind() ipld.Kind {
	return ipld.Kind(kind(p.node_ref))
}

func (p *proxyNode) LookupByString(key string) (ipld.Node, error) {
	ptr, l := memFor([]byte(key))
	nr, errno := lookup_by_string(p.node_ref, ptr, l)
	if errno != 0 {
		return nil, ErrCode{errno}
	}

	return &proxyNode{nr}, nil
}

func (p *proxyNode) LookupByNode(key ipld.Node) (ipld.Node, error) {
	if pk, ok := key.(*proxyNode); ok {
		nr, errno := lookup_by_node(p.node_ref, pk.node_ref)
		if errno != 0 {
			return nil, ErrCode{errno}
		}
		return &proxyNode{nr}, nil
	}
	str, err := key.AsString()
	if err != nil {
		return nil, err
	}
	return p.LookupByString(str)
}

func (p *proxyNode) LookupByIndex(idx int64) (ipld.Node, error) {
	nr, errno := lookup_by_index(p.node_ref, uint64(idx))
	if errno != 0 {
		return nil, ErrCode{errno}
	}
	return &proxyNode{nr}, nil
}

func (p *proxyNode) LookupBySegment(seg ipld.PathSegment) (ipld.Node, error) {
	if p.Kind() == ipld.Kind_List {
		idx, err := seg.Index()
		if err != nil {
			return nil, err
		}
		return p.LookupByIndex(idx)
	}
	return p.LookupByString(seg.String())
}

func (p *proxyNode) MapIterator() ipld.MapIterator {
	mr := map_iterator(p.node_ref)
	return &proxyMapIterator{mr}
}

func (p *proxyNode) ListIterator() ipld.ListIterator {
	lr := list_iterator(p.node_ref)
	return &proxyListIterator{lr}
}

func (p *proxyNode) Length() int64 {
	return int64(length(p.node_ref))
}

func (p *proxyNode) IsAbsent() bool {
	return is_absent(p.node_ref) != 0
}

func (p *proxyNode) IsNull() bool {
	return is_null(p.node_ref) != 0
}

func (p *proxyNode) AsBool() (bool, error) {
	bv, errno := as_bool(p.node_ref)
	if errno != 0 {
		return false, ErrCode{errno}
	}
	return bv != 0, nil
}

func (p *proxyNode) AsInt() (int64, error) {
	iv, errno := as_int(p.node_ref)
	if errno != 0 {
		return 0, ErrCode{errno}
	}
	return int64(iv), nil
}

func (p *proxyNode) AsFloat() (float64, error) {
	fv, errno := as_float(p.node_ref)
	if errno != 0 {
		return 0.0, ErrCode{errno}
	}
	return float64(fv), nil
}

func (p *proxyNode) AsString() (string, error) {
	buf, err := getString(p.node_ref)
	if err != nil {
		return "", err
	}
	return string(buf), nil
}

func (p *proxyNode) AsBytes() ([]byte, error) {
	buf, err := getBytes(p.node_ref)
	if err != nil {
		return nil, err
	}
	return buf, nil
}
func (p *proxyNode) AsLink() (ipld.Link, error) {
	buf, err := getLink(p.node_ref)
	if err != nil {
		return nil, err
	}
	_, c, err := cid.CidFromBytes(buf)
	if err != nil {
		return nil, err
	}
	return cidlink.Link{Cid: c}, nil
}

func (p *proxyNode) Prototype() ipld.NodePrototype {
	return proxyPrototype{}
}

var _ (ipld.Node) = (*proxyNode)(nil)

type proxyMapIterator struct {
	map_ref
}

func (p *proxyMapIterator) Next() (key ipld.Node, value ipld.Node, err error) {
	kref, valref, errno := map_iterate_next(p.map_ref)
	if errno != 0 {
		return nil, nil, ErrCode{errno}
	}
	return &proxyNode{kref}, &proxyNode{valref}, nil
}

func (p *proxyMapIterator) Done() bool {
	return map_iterate_done(p.map_ref) != 0
}

type proxyListIterator struct {
	list_ref
}

func (p *proxyListIterator) Next() (int64, ipld.Node, error) {
	idx, valref, errno := list_iterate_next(p.list_ref)
	if errno != 0 {
		return 0, nil, ErrCode{errno}
	}
	return int64(idx), &proxyNode{valref}, nil
}

func (p *proxyListIterator) Done() bool {
	return list_iterate_done(p.list_ref) != 0
}
