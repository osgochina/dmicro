package metric

import "github.com/osgochina/dmicro/drpc"

type Options struct {
	Host        string
	Port        int
	Path        string
	ServiceName string
	Plugins     []drpc.Plugin
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
		options.ServiceName = name
	}
}

func OptPlugin(plugin drpc.Plugin) Option {
	return func(options *Options) {
		options.Plugins = append(options.Plugins, plugin)
	}
}
