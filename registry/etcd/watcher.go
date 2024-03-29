package etcd

import (
	"context"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/osgochina/dmicro/registry"
	clientv3 "go.etcd.io/etcd/client/v3"
	"time"
)

//
type etcdWatcher struct {
	stop    chan bool
	w       clientv3.WatchChan
	client  *clientv3.Client
	timeout time.Duration
}

func newEtcdWatcher(r *etcdRegistry, timeout time.Duration, opts ...registry.WatchOption) (registry.Watcher, error) {
	var wo registry.WatchOptions
	for _, o := range opts {
		o(&wo)
	}
	ctx, cancel := context.WithCancel(context.Background())
	stop := make(chan bool, 1)
	go func() {
		<-stop
		cancel()
	}()
	watchPath := prefix
	if len(wo.Service) > 0 {
		watchPath = servicePath(wo.Service) + "/"
	}
	return &etcdWatcher{
		stop:    stop,
		w:       r.client.Watch(ctx, watchPath, clientv3.WithPrefix(), clientv3.WithPrevKV()),
		client:  r.client,
		timeout: timeout,
	}, nil
}

func (that *etcdWatcher) Next() (*registry.Result, error) {
	for wresp := range that.w {
		if wresp.Err() != nil {
			return nil, wresp.Err()
		}
		if wresp.Canceled {
			return nil, gerror.New("could not get next")
		}
		for _, ev := range wresp.Events {
			service := decode(ev.Kv.Value)
			var action registry.EventType

			switch ev.Type {
			case clientv3.EventTypePut:
				if ev.IsCreate() {
					action = registry.Create
				} else if ev.IsModify() {
					action = registry.Update
				}
			case clientv3.EventTypeDelete:
				action = registry.Delete

				// get service from prevKv
				service = decode(ev.PrevKv.Value)
			}

			if service == nil {
				continue
			}
			return &registry.Result{
				Action:  action,
				Service: service,
			}, nil
		}
	}
	return nil, gerror.New("could not get next")
}

func (that *etcdWatcher) Stop() {
	select {
	case <-that.stop:
		return
	default:
		close(that.stop)
	}
}
