package main

import (
	"time"

	"github.com/dxyinme/dfinder-go/discover"
	"github.com/sirupsen/logrus"
)

func main() {
	etcd_conf := discover.DefaultEtcdCfg
	etcd_conf.Endpoints = []string{"127.0.0.1:2379"}
	di, err := discover.NewDiscoverWithEtcdCfg("dev", etcd_conf)
	if err != nil {
		logrus.Fatal(err)
	}
	for {
		addr, err := di.GetRandomAddr("test_py")
		if err != nil {
			logrus.Error(err)
		}
		logrus.Infof("addr : %s", addr)
		time.Sleep(time.Second * 5)
	}
}
