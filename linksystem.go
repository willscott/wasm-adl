package main

import (
	"bytes"
	"io"

	"github.com/ipld/go-ipld-prime"
	cidlink "github.com/ipld/go-ipld-prime/linking/cid"
)

// hostStorageReadOpener will just put the bytes of the link into a buffer/
// the correspnding `DecoderChooser`/`decoder` will expect this form and use it
// to make a load call out to the host for a node_ref
func hostStorageReadOpener(_ ipld.LinkContext, l ipld.Link) (io.Reader, error) {
	if cl, ok := l.(cidlink.Link); ok {
		b := cl.Cid.Bytes()
		return bytes.NewBuffer(b), nil
	}
	s := l.String()
	return bytes.NewBufferString(s), nil
}

// hostDecoderChooser returns a hostDecoder
func hostDecoderChooser(_ ipld.Link) (ipld.Decoder, error) {
	return hostDecoder, nil
}

// hostDecoder proxies creation of proxy nodes against a host linksystem.
func hostDecoder(na ipld.NodeAssembler, r io.Reader) error {
	if bf, ok := r.(*bytes.Buffer); ok {
		by := bf.Bytes()
		p, s := memFor(by)
		nr := load(p, s)
		return na.AssignNode(&proxyNode{nr})
	}
	return ErrCode{ErrUnknown}
}

// TODO: writing blocks / mutation not yet supported.
func hostStorageWriteOpener(ipld.LinkContext) (io.Writer, ipld.BlockWriteCommitter, error) {
	return nil, nil, nil
}

func proxyLinkSystem() *ipld.LinkSystem {
	return &ipld.LinkSystem{
		EncoderChooser:     nil,
		DecoderChooser:     hostDecoderChooser,
		HasherChooser:      nil,
		StorageWriteOpener: hostStorageWriteOpener,
		StorageReadOpener:  hostStorageReadOpener,
		TrustedStorage:     true,
		NodeReifier:        nil,
	}
}
