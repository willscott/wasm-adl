//go:generate tinygo build -opt=s -wasm-abi=generic -target=wasm -o adl.wasm ./
package main

import (
	ipld "github.com/ipld/go-ipld-prime"
	cidlink "github.com/ipld/go-ipld-prime/linking/cid"
)

// ADL is the ADL to compile into wasm
var ADL ipld.NodeReifier

// NodeInst is the active instance of the ADL
var NodeInst ipld.Node

func main() {}

//export adl_reify
func adl_reify(n node_ref) errno {
	if NodeInst != nil {
		return ErrAlreadyExists
	}

	ls := proxyLinkSystem()

	node, err := ADL(ipld.LinkContext{}, &proxyNode{n}, ls)
	if err != nil {
		return toErrno(err)
	}
	NodeInst = node
	return 0
}

//export adl_lookup_by_node
func adl_lookup_by_node(node_ref) (node_ref, errno) {
	if NodeInst == nil {
		return 0, ErrDoesNotExist
	}
	return 0, 0
}

//export adl_lookup_by_index
func adl_lookup_by_index(idx int64) (node_ref, errno) {
	if NodeInst == nil {
		return 0, ErrDoesNotExist
	}
	return 0, 0
}

//export adl_kind
func adl_kind() uint8 {
	if NodeInst == nil {
		return 0
	}
	return uint8(NodeInst.Kind())
}

const NumIterators = 32

var mapIterators [NumIterators]ipld.MapIterator
var listIterators [NumIterators]ipld.ListIterator

//export adl_map_iterator
func adl_map_iterator() uint32 {
	if NodeInst == nil {
		return NumIterators + 1
	}

	for i := 0; i < NumIterators; i++ {
		if mapIterators[i] == nil {
			mapIterators[i] = NodeInst.MapIterator()
			return uint32(i)
		}
	}
	return NumIterators + 1
}

//export adl_map_iterator_next
func adl_map_iterator_next(mi_ref uint32) (node_ref, node_ref, errno) {
	if mi_ref >= NumIterators || mapIterators[mi_ref] == nil {
		return 0, 0, ErrDoesNotExist
	}
	key, val, err := mapIterators[mi_ref].Next()
	if err != nil {
		return 0, 0, toErrno(err)
	}

	// get key and val transformed to nodes on the host side.
	kb := proxyPrototype{}.NewBuilder()
	if err := kb.AssignNode(key); err != nil {
		return 0, 0, toErrno(err)
	}

	vb := proxyPrototype{}.NewBuilder()
	if err := vb.AssignNode(val); err != nil {
		return 0, 0, toErrno(err)
	}

	return kb.(*proxyBuilder).node_ref, vb.(*proxyBuilder).node_ref, 0
}

//export adl_map_iterator_done
func adl_map_iterator_done(mi_ref uint32) uint8 {
	if mi_ref >= NumIterators || mapIterators[mi_ref] == nil {
		return ErrDoesNotExist
	}
	done := mapIterators[mi_ref].Done()
	if done {
		return 1
	}
	return 0
}

//export adl_list_iterator
func adl_list_iterator() uint32 {
	if NodeInst == nil {
		return NumIterators + 1
	}

	for i := 0; i < NumIterators; i++ {
		if listIterators[i] == nil {
			listIterators[i] = NodeInst.ListIterator()
			return uint32(i)
		}
	}
	return NumIterators + 1
}

//export adl_list_iterator_next
func adl_list_iterator_next(li_ref uint32) (uint64, node_ref, errno) {
	if li_ref >= NumIterators || listIterators[li_ref] == nil {
		return 0, 0, ErrDoesNotExist
	}
	idx, val, err := listIterators[li_ref].Next()
	if err != nil {
		return 0, 0, toErrno(err)
	}

	// get key and val transformed to nodes on the host side.
	vb := proxyPrototype{}.NewBuilder()
	if err := vb.AssignNode(val); err != nil {
		return 0, 0, toErrno(err)
	}

	return uint64(idx), vb.(*proxyBuilder).node_ref, 0
}

//export adl_list_iterator_done
func adl_list_iterator_done(li_ref uint32) uint8 {
	if li_ref >= NumIterators || listIterators[li_ref] == nil {
		return ErrDoesNotExist
	}
	done := listIterators[li_ref].Done()
	if done {
		return 1
	}
	return 0
}

//export adl_length
func adl_length() uint64 {
	if NodeInst == nil {
		return 0
	}
	return uint64(NodeInst.Length())
}

//export adl_is_absent
func adl_is_absent() uint8 {
	if NodeInst == nil {
		return ErrDoesNotExist
	}
	if NodeInst.IsAbsent() {
		return 1
	}
	return 0
}

//export adl_is_null
func adl_is_null() uint8 {
	if NodeInst == nil {
		return ErrDoesNotExist
	}
	if NodeInst.IsNull() {
		return 1
	}
	return 0
}

//export adl_as_bool
func adl_as_bool() (uint8, errno) {
	if NodeInst == nil {
		return 0, ErrDoesNotExist
	}
	b, err := NodeInst.AsBool()
	if err != nil {
		return 0, toErrno(err)
	}
	if b {
		return 1, 0
	}
	return 0, 0
}

//export adl_as_int
func adl_as_int() (uint64, errno) {
	if NodeInst == nil {
		return 0, ErrDoesNotExist
	}

	i, err := NodeInst.AsInt()
	if err != nil {
		return 0, toErrno(err)
	}

	return uint64(i), 0
}

//export adl_as_float
func adl_as_float() (uint64, errno) {
	if NodeInst == nil {
		return 0, ErrDoesNotExist
	}

	f, err := NodeInst.AsFloat()
	if err != nil {
		return 0, toErrno(err)
	}

	return uint64(f), 0
}

//export adl_as_string
func adl_as_string() (ptr, size, errno) {
	if NodeInst == nil {
		return 0, 0, ErrDoesNotExist
	}

	s, err := NodeInst.AsString()
	if err != nil {
		return 0, 0, toErrno(err)
	}
	ptr, size := memFor([]byte(s))
	return ptr, size, 0
}

//export adl_as_bytes
func adl_as_bytes() (ptr, size, errno) {
	if NodeInst == nil {
		return 0, 0, ErrDoesNotExist
	}

	b, err := NodeInst.AsBytes()
	if err != nil {
		return 0, 0, toErrno(err)

	}
	ptr, size := memFor(b)
	return ptr, size, 0
}

//export adl_as_link
func adl_as_link() (ptr, size, errno) {
	if NodeInst == nil {
		return 0, 0, ErrDoesNotExist
	}

	l, err := NodeInst.AsLink()
	if err != nil {
		return 0, 0, toErrno(err)
	}

	if cl, ok := l.(cidlink.Link); ok {
		b := cl.Cid.Bytes()
		ptr, size := memFor(b)
		return ptr, size, 0
	}
	s := l.String()
	ptr, size := memFor([]byte(s))
	return ptr, size, 0
}
