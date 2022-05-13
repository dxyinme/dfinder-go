package main

import (
	"time"

	"github.com/dxyinme/dfinder-go/discover"
	"github.com/sirupsen/logrus"
)

func main() {
	etcd_conf := discover.DefaultEtcdCfg
	etcd_conf.Endpoints = []string{"127.0.0.1:2379"}
	di, err := discover.NewDiscoverWithEtcdCfg("1", "", "dev", etcd_conf)
	if err != nil {
		logrus.Fatal(err)
	}
	for {
		addrs, err := di.GetAllAddrs("testname1", "dev")
		if err != nil {
			logrus.Error(err)
		}
		logrus.Infof("addrs : %v", addrs)
		time.Sleep(time.Second * 5)
	}
}
