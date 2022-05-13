package discover

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/thanhpk/randstr"
	clientv3 "go.etcd.io/etcd/client/v3"
)

const (
	TTLSec                = 15 * time.Second
	TTLInt                = 15
	RevoteSec             = 10 * time.Second
	TokenLength           = 16
	EtcdClientDialTimeout = 5 * time.Second
)

type RegisterImpl struct {
	ServerName string
	Env        string
	Addr       string
	Token      string
	Cfg        clientv3.Config
	cli        *clientv3.Client
}

func (ri *RegisterImpl) Init(servername, addr, env string, cfg clientv3.Config) (err error) {
	ri.ServerName = servername
	ri.Addr = addr
	ri.Env = env
	ri.Token = randstr.String(TokenLength)
	ri.Cfg = cfg
	ri.cli, err = clientv3.New(ri.Cfg)
	if err != nil {
		return err
	}
	return
}

func (ri *RegisterImpl) Serve() {
	grant, err := ri.cli.Grant(context.Background(), TTLInt)
	if err != nil {
		logrus.Errorf("grant lease error %v", err)
		return
	}

	_, err = ri.cli.Put(context.Background(), ri.makeKey(), ri.Addr, clientv3.WithLease(grant.ID))
	if err != nil {
		logrus.Errorf("register error %v", err)
		return
	}
	cnt := 0
	for {
		_, err = ri.cli.KeepAliveOnce(context.Background(), grant.ID)
		if err != nil {
			logrus.Errorf("keep alive error %v", err)
			return
		} else {
			logrus.Infof("keep alive count : %d", cnt)
			cnt++
		}
		time.Sleep(RevoteSec)
	}
}

func (ri *RegisterImpl) makeKey() string {
	return ri.Env + "/" + ri.ServerName + "/" + ri.Token
}
