package drpc

import (
	"crypto/tls"
	"github.com/gogf/gf/container/gtype"
	"github.com/osgochina/dmicro/logger"
	"net"
	"time"
)

// Dialer 拨号器
type Dialer struct {
	//网络类型
	network string
	//本地使用的地址端口
	localAddr net.Addr
	//tls配置信息
	tlsConfig *tls.Config
	//拨号器拨号超时时间
	dialTimeout time.Duration
	//拨号器重复拨号的时间间隔
	redialInterval time.Duration
	//拨号器重复拨号的最大次数
	redialTimes int
}

// NewDialer 创建一个拨号器
func NewDialer(localAddr net.Addr, tlsConfig *tls.Config, dialTimeout, redialInterval time.Duration, redialTimes int) *Dialer {
	return &Dialer{
		network:        localAddr.Network(),
		localAddr:      localAddr,
		tlsConfig:      tlsConfig,
		dialTimeout:    dialTimeout,
		redialTimes:    redialTimes,
		redialInterval: redialInterval,
	}
}

// Network 获取拨号器的网络类型
func (that *Dialer) Network() string {
	return that.network
}

// LocalAddr 获取拨号器本地使用的端口地址
func (that *Dialer) LocalAddr() net.Addr {
	return that.localAddr
}

// TLSConfig 获取tls配置信息
func (that *Dialer) TLSConfig() *tls.Config {
	return that.tlsConfig
}

// DialTimeout 获取拨号器拨号时候的超时时间
func (that *Dialer) DialTimeout() time.Duration {
	return that.dialTimeout
}

// RedialInterval 返回拨号器重试拨号时候的间隔
func (that *Dialer) RedialInterval() time.Duration {
	return that.redialInterval
}

// RedialTimes 拨号器重复拨号的最大次数
func (that *Dialer) RedialTimes() int {
	return that.redialTimes
}

// Dial 拨号链接地址 addr
func (that *Dialer) Dial(addr string) (net.Conn, error) {
	return that.dialWithRetry(addr, "", nil)
}

// 拨号，如果拨号失败，则重试 redialTimes 次
func (that *Dialer) dialWithRetry(addr, sessID string, fn func(conn net.Conn) error) (net.Conn, error) {
	conn, err := that.dialOne(addr)
	if err == nil {
		if fn == nil {
			return conn, nil
		} else {
			err = fn(conn)
			if err == nil {
				return conn, nil
			}
		}
	}
	redialTimes := that.newRedialCounter()

	for redialTimes.Add(-1) > 0 {
		time.Sleep(that.redialInterval)
		if sessID == "" {
			logger.Debugf("trying to redial... (network:%s, addr:%s)", that.network, addr)
		} else {
			logger.Debugf("trying to redial... (network:%s, addr:%s, id:%s)", that.network, addr, sessID)
		}
		conn, err = that.dialOne(addr)
		if err == nil {
			if fn == nil {
				return conn, nil
			} else {
				err = fn(conn)
				if err == nil {
					return conn, nil
				}
			}
		}
	}
	return nil, err
}

//拨号一次
func (that *Dialer) dialOne(addr string) (net.Conn, error) {
	//if network := asQUIC(d.network); network != "" {
	//	ctx := context.Background()
	//	if d.dialTimeout > 0 {
	//		ctx, _ = context.WithTimeout(ctx, d.dialTimeout)
	//	}
	//	var tlsConf = d.tlsConfig
	//	if tlsConf == nil {
	//		tlsConf = GenerateTLSConfigForClient()
	//	}
	//	return quic.DialAddrContext(ctx, network, d.localAddr.(*FakeAddr).udpAddr, addr, tlsConf, nil)
	//}
	//
	//if network := asKCP(d.network); network != "" {
	//	return kcp.DialAddrContext(network, d.localAddr.(*FakeAddr).udpAddr, addr, d.tlsConfig, dataShards, parityShards)
	//}

	dialer := &net.Dialer{
		LocalAddr: that.localAddr,
		Timeout:   that.dialTimeout,
	}
	//使用tls加密拨号
	if that.tlsConfig != nil {
		return tls.DialWithDialer(dialer, that.network, addr, that.tlsConfig)
	}
	return dialer.Dial(that.network, addr)
}

func (that *Dialer) newRedialCounter() *gtype.Int {
	return gtype.NewInt(that.redialTimes)
}
