module main

require (
	github.com/gogf/gf v1.16.6
	github.com/osgochina/dmicro v0.0.0-00010101000000-000000000000
	github.com/osgochina/dmicro/registry/etcd v0.0.0-00010101000000-000000000000
)

replace github.com/osgochina/dmicro => ../../../dmicro

replace github.com/osgochina/dmicro/registry/etcd => ../../../dmicro/registry/etcd

go 1.15
