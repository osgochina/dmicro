package graceful

import (
	"github.com/gogf/gf/util/gconv"
	"github.com/osgochina/dmicro/logger"
	"net"
)

type listenerFdMap = map[string]string

// 获取监听列表
func (that *graceful) getGHttpListenerFdMap() map[string]listenerFdMap {
	if that.model == GraceChangeProcess {
		return nil
	}
	if that.inheritedProcListener.Len() <= 0 {
		return nil
	}
	sfm := make(map[string]listenerFdMap)
	m := map[string]string{
		"https": "",
		"http":  "",
	}
	that.inheritedProcListener.Iterator(func(_ int, v interface{}) bool {
		lis, ok := v.(net.Listener)
		if !ok {
			logger.Warningf("inheritedProcListener 不是 net.Listener类型")
			return true
		}
		if that.mwListenAddr == nil {
			return true
		}
		// 判断监听的是否是http协议。如果是http协议则不返回
		data := that.mwListenAddr.Get(lis.Addr().String())
		d, ok := data.(InheritAddr)
		if ok && d.Network != "http" && d.Network != "https" {
			return true
		}
		f, e := lis.(filer).File()
		if e != nil {
			logger.Error(e)
			return true
		}
		str := lis.Addr().String() + "#" + gconv.String(f.Fd()) + ","
		if d.Network == "https" {
			if len(m["https"]) > 0 {
				m["https"] += ","
			}
			m["https"] += str
		} else {
			if len(m["http"]) > 0 {
				m["http"] += ","
			}
			m["http"] += str
		}
		sfm[d.ServerName] = m
		return true
	})
	return sfm
}
