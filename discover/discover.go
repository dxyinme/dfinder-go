package discover

import (
	"github.com/sirupsen/logrus"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
)

type Discover interface {
	Init() (err error)

	// 获取这个服务的所有grpc连接
	GetAllGrpcConn(servername, env string) (cliList []*grpc.ClientConn, err error)

	// 获取这个服务的随机一个grpc连接
	GetRandomGrpcConn(servername, env string) (cli *grpc.ClientConn, err error)

	// 获取这个服务的所有服务地址
	GetAllAddrs(servername, env string) (addrs []string, err error)
}

func NewDiscoverWithEtcdCfg(servername, addr, env string, cfg clientv3.Config) (d Discover, err error) {
	d = &DiscoverImpl{
		Cfg: cfg,
	}
	err = d.Init()
	if err != nil {
		logrus.Error(err)
	}
	return
}
