package main

import (
	"github.com/ipld/go-ipld-prime"
	cidlink "github.com/ipld/go-ipld-prime/linking/cid"
)

type proxyPrototype struct {
}

func (p proxyPrototype) NewBuilder() ipld.NodeBuilder {
	return &proxyBuilder{}
}

type proxyBuilder struct {
	node_ref
}

func (pb *proxyBuilder) BeginMap(sizeHint int64) (ipld.MapAssembler, error) {
	mapRef := new_map(size(sizeHint))
	return &proxyMapBuilder{pb, mapRef, nil, nil, 0}, nil
}

func (pb *proxyBuilder) BeginList(sizeHint int64) (ipld.ListAssembler, error) {
	listRef := new_list(size(sizeHint))
	return &proxyListBuilder{pb, listRef, nil, 0}, nil
}

func (pb *proxyBuilder) AssignNull() error {
	if pb.node_ref != 0 {
		return ErrCode{ErrAlreadyExists}
	}
	pb.node_ref = new_null()
	return nil
}
func (pb *proxyBuilder) AssignBool(b bool) error {
	if pb.node_ref != 0 {
		return ErrCode{ErrAlreadyExists}
	}
	if b {
		pb.node_ref = new_bool(1)
	} else {
		pb.node_ref = new_bool(0)
	}
	return nil
}
func (pb *proxyBuilder) AssignInt(i int64) error {
	if pb.node_ref != 0 {
		return ErrCode{ErrAlreadyExists}
	}
	pb.node_ref = new_int(uint64(i))
	return nil
}

func (pb *proxyBuilder) AssignFloat(f float64) error {
	if pb.node_ref != 0 {
		return ErrCode{ErrAlreadyExists}
	}
	pb.node_ref = new_int(uint64(f))
	return nil
}

func (pb *proxyBuilder) AssignString(s string) error {
	if pb.node_ref != 0 {
		return ErrCode{ErrAlreadyExists}
	}
	ptr, size := memFor([]byte(s))
	pb.node_ref = new_string(ptr, size)
	return nil
}

func (pb *proxyBuilder) AssignBytes(b []byte) error {
	if pb.node_ref != 0 {
		return ErrCode{ErrAlreadyExists}
	}
	ptr, size := memFor(b)
	pb.node_ref = new_bytes(ptr, size)
	return nil
}

func (pb *proxyBuilder) AssignLink(l ipld.Link) error {
	if pb.node_ref != 0 {
		return ErrCode{ErrAlreadyExists}
	}
	if cl, ok := l.(cidlink.Link); ok {
		b := cl.Cid.Bytes()
		ptr, size := memFor(b)
		pb.node_ref = new_link(ptr, size)
	} else {
		s := l.String()
		ptr, size := memFor([]byte(s))
		pb.node_ref = new_link(ptr, size)
	}
	return nil
}

func (pb *proxyBuilder) AssignNode(n ipld.Node) error {
	if pb.node_ref != 0 {
		return ErrCode{ErrAlreadyExists}
	}
	if pn, ok := n.(*proxyNode); ok {
		pb.node_ref = pn.node_ref
		return nil
	}
	switch n.Kind() {
	case ipld.Kind_Null:
		return pb.AssignNull()
	case ipld.Kind_Bool:
		b, err := n.AsBool()
		if err != nil {
			return err
		}
		return pb.AssignBool(b)
	case ipld.Kind_Int:
		i, err := n.AsInt()
		if err != nil {
			return err
		}
		return pb.AssignInt(i)
	case ipld.Kind_Float:
		f, err := n.AsFloat()
		if err != nil {
			return err
		}
		return pb.AssignFloat(f)
	case ipld.Kind_String:
		s, err := n.AsString()
		if err != nil {
			return err
		}
		return pb.AssignString(s)
	case ipld.Kind_Bytes:
		b, err := n.AsBytes()
		if err != nil {
			return err
		}
		return pb.AssignBytes(b)
	case ipld.Kind_Link:
		l, err := n.AsLink()
		if err != nil {
			return err
		}
		return pb.AssignLink(l)
	case ipld.Kind_Map:

	case ipld.Kind_List:
	}

	return ErrCode{ErrUnknown}
}

func (pb *proxyBuilder) Prototype() ipld.NodePrototype {
	return proxyPrototype{}
}

func (pb *proxyBuilder) Build() ipld.Node {
	if pb.node_ref == 0 {
		return nil
	}
	return &proxyNode{pb.node_ref}
}
func (pb *proxyBuilder) Reset() {
	pb.node_ref = 0
}

type proxyMapBuilder struct {
	pb     *proxyBuilder
	mapRef map_builder_ref
	kb     *proxyBuilder
	vb     *proxyBuilder
	errno
}

func (pmb *proxyMapBuilder) AssembleKey() ipld.NodeAssembler {
	if pmb.kb != nil {
		if pmb.vb != nil && pmb.kb.node_ref != 0 && pmb.vb.node_ref != 0 {
			// assemble previous entry.
			errno := assemble_map(pmb.mapRef, pmb.kb.node_ref, pmb.vb.node_ref)
			if errno == 0 {
				pmb.kb = nil
				pmb.vb = nil
			} else {
				pmb.errno = errno
			}
		} else {
			return nil
		}
	}
	pmb.kb = &proxyBuilder{}
	return pmb.kb
}

func (pmb *proxyMapBuilder) AssembleValue() ipld.NodeAssembler {
	if pmb.vb != nil {
		return nil
	}
	pmb.vb = &proxyBuilder{}
	return pmb.vb
}

func (pmb *proxyMapBuilder) AssembleEntry(k string) (ipld.NodeAssembler, error) {
	kb := pmb.AssembleKey()
	if kb == nil || pmb.errno != 0 {
		return nil, ErrCode{pmb.errno}
	}
	if err := kb.AssignString(k); err != nil {
		return nil, err
	}
	return pmb.AssembleValue(), nil
}

func (pmb *proxyMapBuilder) Finish() error {
	if pmb.errno != 0 {
		return ErrCode{pmb.errno}
	}

	if pmb.kb != nil && pmb.vb != nil && pmb.kb.node_ref != 0 && pmb.vb.node_ref != 0 {
		// assemble previous entry.
		errno := assemble_map(pmb.mapRef, pmb.kb.node_ref, pmb.vb.node_ref)
		if errno != 0 {
			return ErrCode{errno}
		}
	}

	nr := finish_map(pmb.mapRef)
	if nr == 0 {
		return ErrCode{ErrUnknown}
	}
	pmb.pb.node_ref = nr
	return nil
}

func (pmb *proxyMapBuilder) KeyPrototype() ipld.NodePrototype {
	return proxyPrototype{}
}

func (pmb *proxyMapBuilder) ValuePrototype(k string) ipld.NodePrototype {
	return proxyPrototype{}
}

type proxyListBuilder struct {
	pb      *proxyBuilder
	listRef list_builder_ref
	vb      *proxyBuilder
	errno
}

func (plb *proxyListBuilder) AssembleValue() ipld.NodeAssembler {
	if plb.vb != nil && plb.vb.node_ref != 0 {
		// assemble previous entry.
		errno := assemble_list(plb.listRef, plb.vb.node_ref)
		if errno == 0 {
			plb.vb = nil
		} else {
			plb.errno = errno
		}
	} else if plb.vb != nil {
		return nil
	}
	plb.vb = &proxyBuilder{}
	return plb.vb
}

func (plb *proxyListBuilder) Finish() error {
	if plb.errno != 0 {
		return ErrCode{plb.errno}
	}

	if plb.vb != nil && plb.vb.node_ref != 0 {
		// assemble previous entry.
		errno := assemble_list(plb.listRef, plb.vb.node_ref)
		if errno != 0 {
			return ErrCode{errno}
		}
	}

	nr := finish_list(plb.listRef)
	if nr == 0 {
		return ErrCode{ErrUnknown}
	}
	plb.pb.node_ref = nr
	return nil
}

func (plb *proxyListBuilder) ValuePrototype(idx int64) ipld.NodePrototype {
	return proxyPrototype{}
}
