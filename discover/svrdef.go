package discover

import (
	"go.etcd.io/etcd/api/v3/mvccpb"
	"google.golang.org/grpc"
)

type SvrDef struct {
	ServerName string
	AddrsMap   map[string]string
	GrpcCliMap map[string]*grpc.ClientConn
}

func NewSvrDef(servername string, kvs *[]*mvccpb.KeyValue) SvrDef {
	ret := SvrDef{
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
