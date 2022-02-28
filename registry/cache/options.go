package cache

import (
	"time"
)

type Options struct {
	// TTL 缓存的生命周期
	TTL time.Duration
}

type Option func(o *Options)

// WithTTL sets the cache TTL
func WithTTL(t time.Duration) Option {
	return func(o *Options) {
		o.TTL = t
	}
}
