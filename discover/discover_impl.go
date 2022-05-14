package discover

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
)

const (
	UpdateSvrDefIntervalSec = 5 // Sec
)

func makePrefix(server, env string) string {
	return env + "/" + server
}

type DiscoverImpl struct {
	Cfg        clientv3.Config
	mpMu       sync.RWMutex
	Str2SvrDef map[string]*SvrDef
	updMu      sync.Mutex
	// 某个server 上一次update SvrDef 的时间
	updDef map[string]int64
	cli    *clientv3.Client
}

func (di *DiscoverImpl) Init() (err error) {
	di.Str2SvrDef = make(map[string]*SvrDef)
	di.updDef = make(map[string]int64)
	di.cli, err = clientv3.New(di.Cfg)
	return
}

func (di *DiscoverImpl) GetAllGrpcConn(servername, env string) (cliList []*grpc.ClientConn, err error) {
	return
}
func (di *DiscoverImpl) GetRandomGrpcConn(servername, env string) (cli *grpc.ClientConn, err error) {
	return
}
func (di *DiscoverImpl) GetAllAddrs(servername, env string) (addrs []string, err error) {
	di.updateAddr(servername, env)
	di.mpMu.RLock()
	defer di.mpMu.RUnlock()
	svrDef, ok := di.Str2SvrDef[servername]
	if !ok {
		return nil, errors.New("server not found")
	}
	return svrDef.AddrsList(), nil
}

func (di *DiscoverImpl) updateAddr(servername, env string) {
	now_time := time.Now().Unix()
	di.updMu.Lock()
	last_time := di.updDef[servername]
	di.updMu.Unlock()
	if now_time-last_time > UpdateSvrDefIntervalSec {
		resp, err := di.cli.Get(context.Background(), makePrefix(servername, env), clientv3.WithPrefix())
		logrus.Infof("resp: %v", resp)
		if err != nil {
			logrus.Error(err)
		}
		di.mpMu.Lock()
		if svr_def, ok := di.Str2SvrDef[servername]; !ok {
			di.Str2SvrDef[servername] = NewSvrDef(servername, &resp.Kvs)
		} else {
			svr_def.Update(&resp.Kvs)
		}
		di.mpMu.Unlock()

		di.updMu.Lock()
		di.updDef[servername] = now_time
		di.updMu.Unlock()
	}
}
