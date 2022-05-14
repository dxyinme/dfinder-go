package discover

import (
	"sync"

	"go.etcd.io/etcd/api/v3/mvccpb"
	"google.golang.org/grpc"
)

type SvrDef struct {
	ServerName string
	AddrsMap   map[string]string
	AddrsMapMu sync.RWMutex
	GrpcCliMap map[string]*grpc.ClientConn
}

func NewSvrDef(servername string, kvs *[]*mvccpb.KeyValue) *SvrDef {
	ret := &SvrDef{
		ServerName: servername,
		AddrsMap:   make(map[string]string),
	}
	if kvs == nil || len(*kvs) == 0 {
	} else {
		for _, v := range *kvs {
			ret.AddrsMap[string(v.Key)] = string(v.Value)
		}
	}
	return ret
}

func (sd *SvrDef) Update(kvs *[]*mvccpb.KeyValue) {
	sd.AddrsMapMu.Lock()
	defer sd.AddrsMapMu.Unlock()
	for _, v := range *kvs {
		sd.AddrsMap[string(v.Key)] = string(v.Value)
	}
}

func (sd *SvrDef) AddrsList() (addrs []string) {
	addrs = make([]string, 0)
	sd.AddrsMapMu.RLock()
	defer sd.AddrsMapMu.RUnlock()
	for _, v := range sd.AddrsMap {
		addrs = append(addrs, v)
	}
	return
}
