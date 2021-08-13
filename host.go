package main

import (
	"github.com/willscott/wasm-adl/common"
)

type ptr = common.Ptr
type size = common.Size
type node_ref = common.Node_ref
type map_ref = common.Map_ref
type list_ref = common.List_ref
type map_builder_ref = common.Map_builder_ref
type list_builder_ref = common.List_builder_ref

//export load
func load(ptr, size) node_ref

//export kind
func kind(node_ref) uint8

//export lookup_by_string
func lookup_by_string(node_ref, ptr, size) (node_ref, errno)

//export lookup_by_node
func lookup_by_node(node_ref, node_ref) (node_ref, errno)

//export lookup_by_index
func lookup_by_index(node_ref, uint64) (node_ref, errno)

//export map_iterator
func map_iterator(node_ref) map_ref

//export list_iterator
func list_iterator(node_ref) list_ref

//export length
func length(node_ref) uint64

//export is_absent
func is_absent(node_ref) uint8

//export is_null
func is_null(node_ref) uint8

//export as_bool
func as_bool(node_ref) (uint8, errno)

//export as_int
func as_int(node_ref) (uint64, errno)

//export as_float
func as_float(node_ref) (uint64, errno)

//export as_string
func as_string(node_ref, ptr, size) (size, errno)

//export as_bytes
func as_bytes(node_ref, ptr, size) (size, errno)

//export as_link
func as_link(node_ref, ptr, size) (size, errno)

//export map_iterate_next
func map_iterate_next(map_ref) (node_ref, node_ref, errno)

//export map_iterate_done
func map_iterate_done(map_ref) uint8

//export list_iterate_next
func list_iterate_next(list_ref) (uint64, node_ref, errno)

//export list_iterate_done
func list_iterate_done(list_ref) uint8

//export new_null
func new_null() node_ref

//export new_bool
func new_bool(uint8) node_ref

//export new_int
func new_int(uint64) node_ref

//export new_float
func new_float(uint64) node_ref

//export new_string
func new_string(ptr, size) node_ref

//export new_bytes
func new_bytes(ptr, size) node_ref

//export new_link
func new_link(ptr, size) node_ref

//export new_map
func new_map(size) map_builder_ref

//export new_list
func new_list(size) list_builder_ref

//export assemble_map
func assemble_map(map_builder_ref, node_ref, node_ref) errno

//export finish_map
func finish_map(map_builder_ref) node_ref

//export assemble_list
func assemble_list(list_builder_ref, node_ref) errno

//export finish_list
func finish_list(list_builder_ref) node_ref
