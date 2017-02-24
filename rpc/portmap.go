package rpc

import (
	"fmt"

	"github.com/fdawg4l/nfs/xdr"
)

// PORTMAP
// RFC 1057 Section A.1

const (
	PMAP_PORT = 111
	PMAP_PROG = 100000
	PMAP_VERS = 2

	PMAPPROC_GETPORT = 3
	PMAPPROC_DUMP    = 4

	IPPROTO_TCP = 6
	IPPROTO_UDP = 17
)

type Mapping struct {
	Prog uint32
	Vers uint32
	Prot uint32
	Port uint32
}

type Portmapper struct {
	*Client
	host string
}

func (p *Portmapper) Getport(mapping Mapping) (int, error) {
	type getport struct {
		Header
		Mapping
	}
	msg := &getport{
		Header{
			Rpcvers: 2,
			Prog:    PMAP_PROG,
			Vers:    PMAP_VERS,
			Proc:    PMAPPROC_GETPORT,
			Cred:    AuthNull,
			Verf:    AuthNull,
		},
		mapping,
	}
	buf, err := p.Call(msg)
	if err != nil {
		return 0, err
	}
	port, _ := xdr.Uint32(buf)
	return int(port), nil
}

func DialPortmapper(net, host string) (*Portmapper, error) {
	client, err := DialTCP(net, nil, fmt.Sprintf("%s:%d", host, PMAP_PORT))
	if err != nil {
		return nil, err
	}
	return &Portmapper{client, host}, nil
}
