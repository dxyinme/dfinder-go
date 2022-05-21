package main

import (
	"flag"
	"time"

	"github.com/dxyinme/dfinder-go/discover"
	"github.com/sirupsen/logrus"
)

func main() {
	etcd_addr := flag.String("etcd_addr", "127.0.0.1:2379", "etcd addr")
	svrname := flag.String("svr", "testname1", "etcd addr")
	flag.Parse()
	etcd_conf := discover.DefaultEtcdCfg
	etcd_conf.Endpoints = []string{*etcd_addr}
	di, err := discover.NewDiscoverWithEtcdCfg("dev", etcd_conf)
	if err != nil {
		logrus.Fatal(err)
	}
	for {
		addr, err := di.GetRandomAddr(*svrname)
		if err != nil {
			logrus.Error(err)
		}
		addrs, err := di.GetAllAddrs(*svrname)
		if err != nil {
			logrus.Error(err)
		}
		logrus.Infof("addrs : %v", addrs)
		logrus.Infof("addr  : %s", addr)
		time.Sleep(time.Second * 5)
	}
}
