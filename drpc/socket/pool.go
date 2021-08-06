package socket

import "sync"

var socketPool = sync.Pool{
	New: func() interface{} {
		s := newSocket(nil, nil)
		s.fromPool = true
		return s
	},
}
