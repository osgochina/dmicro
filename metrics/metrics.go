package metrics

type Option func(*Options)

type Metrics interface {
	Init(...Option)
	Options() Options
	Enabled() bool
	Start()
	String() string
}
