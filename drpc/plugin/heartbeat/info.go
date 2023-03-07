package heartbeat

import (
	"github.com/gogf/gf/v2/container/gmap"
	"sync"
	"time"
)

const (
	//心跳最小频率
	minRateSecond = 3
	//心跳组件记录在会话的数据key
	heartbeatSwapKey = "heartbeatSwapKey"
)

// 心跳数据元信息
type heartbeatInfo struct {
	// 心跳频率
	rate time.Duration
	// 上一次心跳的时间
	last time.Time
	mu   sync.RWMutex
}

//获取心跳频率
func (that *heartbeatInfo) getRate() time.Duration {
	that.mu.RLock()
	rate := that.rate
	that.mu.RUnlock()
	return rate
}

//获取上一次心跳时间
func (that *heartbeatInfo) getLast() time.Time {
	that.mu.RLock()
	last := that.last
	that.mu.RUnlock()
	return last
}

// copy信息
func (that *heartbeatInfo) elemCopy() heartbeatInfo {
	that.mu.RLock()
	defer that.mu.RUnlock()

	return heartbeatInfo{
		rate: that.rate,
		last: that.last,
	}
}

//初始化心跳源信息
func initHeartbeatInfo(m *gmap.Map, rate time.Duration) {
	m.Set(heartbeatSwapKey, &heartbeatInfo{
		rate: rate,
		last: time.Now(),
	})
}

//获取心跳原信息
func getHeartbeatInfo(m *gmap.Map) (*heartbeatInfo, bool) {
	_info, ok := m.Search(heartbeatSwapKey)
	if !ok {
		return nil, false
	}
	return _info.(*heartbeatInfo), true
}

//更新心跳源信息
func updateHeartbeatInfo(m *gmap.Map, rate time.Duration) (isFirst bool) {
	info, ok := getHeartbeatInfo(m)
	if !ok {
		isFirst = true
		if rate > 0 {
			initHeartbeatInfo(m, rate)
		}
		return
	}
	info.mu.Lock()
	if rate > 0 {
		info.rate = rate
	}
	info.last = time.Now()
	info.mu.Unlock()
	return
}
