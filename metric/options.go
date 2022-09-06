package metric

import "github.com/osgochina/dmicro/drpc"

type Options struct {
	Host        string
	Port        int
	Path        string
	serviceName string
	plugin      drpc.Plugin
}

func OptHost(host string) Option {
	return func(options *Options) {
		options.Host = host
	}
}

func OptPort(port int) Option {
	return func(options *Options) {
		options.Port = port
	}
}

func OptPath(path string) Option {
	return func(options *Options) {
		options.Path = path
	}
}

func OptServiceName(name string) Option {
	return func(options *Options) {
		options.serviceName = name
	}
}
