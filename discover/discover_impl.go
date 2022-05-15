package discover

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	clientv3 "go.etcd.io/etcd/client/v3"
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
	updMu      sync.RWMutex
	// 某个server 上一次update SvrDef 的时间
	updDef map[string]int64
	cli    *clientv3.Client

	// 初始化的时候必须要设置
	Env string
}

func (di *DiscoverImpl) Init() (err error) {
	di.Str2SvrDef = make(map[string]*SvrDef)
	di.updDef = make(map[string]int64)
	di.cli, err = clientv3.New(di.Cfg)
	if di.Env == "" {
		return errors.New("env is no set")
	}
	return
}

func (di *DiscoverImpl) GetRandomAddr(servername string) (addr string, err error) {
	di.updateDef(servername)
	di.mpMu.RLock()
	defer di.mpMu.RUnlock()
	svrDef, ok := di.Str2SvrDef[servername]
	if !ok {
		return "", errors.New("server not found")
	}
	return svrDef.RandomAddr(), nil
}
func (di *DiscoverImpl) GetAllAddrs(servername string) (addrs []string, err error) {
	di.updateDef(servername)
	di.mpMu.RLock()
	defer di.mpMu.RUnlock()
	svrDef, ok := di.Str2SvrDef[servername]
	if !ok {
		return nil, errors.New("server not found")
	}
	return svrDef.AddrsList(), nil
}

func (di *DiscoverImpl) updateDef(servername string) {
	now_time := time.Now().Unix()
	last_time := di.getDefUpdateTime(servername)
	if now_time-last_time > UpdateSvrDefIntervalSec {
		resp, err := di.cli.Get(context.Background(), makePrefix(servername, di.Env), clientv3.WithPrefix())
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
		di.setDefUpdateTime(servername, now_time)
	}
}

func (di *DiscoverImpl) getDefUpdateTime(servername string) int64 {
	di.updMu.RLock()
	defer di.updMu.RUnlock()
	time_v, ok := di.updDef[servername]
	if !ok {
		return int64(0)
	}
	return time_v
}

func (di *DiscoverImpl) setDefUpdateTime(servername string, time_v int64) {
	di.updMu.Lock()
	defer di.updMu.Unlock()
	di.updDef[servername] = time_v
}
