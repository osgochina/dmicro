package registry

import "sync"

type mdnsRegistry struct {
	sync.Mutex
	opts   Options
	domain string
}
