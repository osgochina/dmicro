package metric

import "github.com/osgochina/dmicro/drpc"

type Option func(*Options)

type Metrics interface {
	Init(...Option)
	Options() Options
	Enabled() bool
	Plugin() []drpc.Plugin
	Start()
	String() string
}
