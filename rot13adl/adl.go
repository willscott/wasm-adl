package rot13adl

import (
	"github.com/ipld/go-ipld-prime"
	rot "github.com/ipld/go-ipld-prime/adl/rot13adl"
)

// Reify creates the reified node version of a node.
func Reify(_ ipld.LinkContext, n ipld.Node, _ *ipld.LinkSystem) (ipld.Node, error) {
	return rot.Reify(n)
}
