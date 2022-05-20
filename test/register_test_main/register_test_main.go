package main

import (
	"flag"

	"github.com/dxyinme/dfinder-go/discover"
	"github.com/sirupsen/logrus"
)

func main() {
	addr := flag.String("addr", "0.0.0.2", "addr")
	etcd_addr := flag.String("etcd_addr", "127.0.0.1:2379", "etcd addr")
	flag.Parse()
	etcd_conf := discover.DefaultEtcdCfg
	etcd_conf.Endpoints = []string{*etcd_addr}
	r, err := discover.NewRegisterWithEtcdCfg("testname1", *addr, "dev", etcd_conf)
	if err != nil {
		logrus.Error(err)
	}
	r.Serve() // 这是个阻塞函数，如果需要后续有后继语句的话，请直接go r.Serve()
}
