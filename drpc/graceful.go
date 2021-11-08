package drpc

import (
	"github.com/gogf/gf/container/gset"
	"github.com/osgochina/dmicro/utils/errors"
	"github.com/osgochina/dmicro/utils/graceful"
	"time"
)

func init() {
	graceful.Graceful().SetShutdownEndpoint(func(endpointList *gset.Set) error {
		var count int
		var errCh = make(chan error, endpointList.Size())
		//异步关闭端点
		for _, val := range endpointList.Slice() {
			count++
			e := val.(*endpoint)
			go func(e *endpoint) {
				errCh <- e.Close()
			}(e)
		}
		var err error
		for i := 0; i < count; i++ {
			err = errors.Merge(err, <-errCh)
		}
		close(errCh)
		return err
	})
}

func Shutdown(timeout ...time.Duration) {
	graceful.Graceful().Shutdown(timeout...)
}
