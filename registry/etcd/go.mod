module github.com/osgochina/dmicro/registry/etcd

require (
	github.com/gogf/gf v1.16.6
	github.com/osgochina/dmicro v0.0.0-00010101000000-000000000000
	go.etcd.io/etcd/api/v3 v3.5.2
	go.etcd.io/etcd/client/v3 v3.5.2
	go.uber.org/zap v1.17.0
)

replace github.com/osgochina/dmicro => ../../../dmicro

go 1.15
