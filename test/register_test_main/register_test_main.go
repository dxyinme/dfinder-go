package main

import (
	"github.com/dxyinme/dfinder-go/discover"
	"github.com/sirupsen/logrus"
)

func main() {
	etcd_conf := discover.DefaultEtcdCfg
	etcd_conf.Endpoints = []string{"127.0.0.1:2379"}
	r, err := discover.NewRegisterWithEtcdCfg("testname1", "0.0.0.2", "dev", etcd_conf)
	if err != nil {
		logrus.Error(err)
	}
	r.Serve() // 这是个阻塞函数，如果需要后续有后继语句的话，请直接go r.Serve()
}
