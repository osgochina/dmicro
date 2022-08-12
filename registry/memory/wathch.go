package memory

import (
	"errors"
	"github.com/osgochina/dmicro/registry"
)

type memWatcher struct {
	id        string
	watchOpts registry.WatchOptions
	result    chan *registry.Result
	exit      chan bool
}

// Next 获取监听结果
func (that *memWatcher) Next() (*registry.Result, error) {
	for {
		select {
		case r := <-that.result:
			// 如果要事件不是要监听的服务的,则忽略
			if len(that.watchOpts.Service) > 0 && that.watchOpts.Service != r.Service.Name {
				continue
			}
			return r, nil
		case <-that.exit:
			return nil, errors.New("watcher stopped")
		}
	}
}

// Stop 停止监听器
func (that *memWatcher) Stop() {
	select {
	case <-that.exit:
		return
	default:
		close(that.exit)
	}
}
