package discover

import clientv3 "go.etcd.io/etcd/client/v3"

var (
	DefaultEtcdCfg = clientv3.Config{
		DialTimeout: EtcdClientDialTimeout,
	}
)

type Register interface {
	Init(servername, addr, env string, cfg clientv3.Config) error
	Serve()
}

func NewRegisterWithEtcdCfg(servername, addr, env string, cfg clientv3.Config) (r Register, err error) {
	r = &RegisterImpl{}
	err = r.Init(servername, addr, env, cfg)
	if err != nil {
		return nil, err
	}
	return
}
