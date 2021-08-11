## WASM ADL

This module builds a scaffold around an [ADL](https://ipld.io/docs/advanced-data-layouts/) written in golang that conforms to the [`NodeReifier`](https://github.com/ipld/go-ipld-prime/blob/master/linksystem.go#L224) interface. The scaffold allows the ADL to be compiled into a wasm module with a light weight proxy for access to the IPLD context it finds itself in, while allowing all non-accessed values to remain outside of the wasm memory boundary.
