package discover

import (
	"github.com/sirupsen/logrus"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type Discover interface {
	Init() (err error)

	// 获取这个服务的随机一个addr
	GetRandomAddr(servername string) (addr string, err error)

	// 获取这个服务的所有服务地址
	GetAllAddrs(servername string) (addrs []string, err error)
}

func NewDiscoverWithEtcdCfg(env string, cfg clientv3.Config) (d Discover, err error) {
	d = &DiscoverImpl{
		Cfg: cfg,
		Env: env,
	}
	err = d.Init()
	if err != nil {
		logrus.Error(err)
	}
	return
}
