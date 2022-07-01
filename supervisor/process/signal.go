package process

import (
	"fmt"
	"github.com/osgochina/dmicro/logger"
	"github.com/osgochina/dmicro/utils/signals"
	"os"
)

// Signal 向进程发送信号
// sig: 要发送的信号
// sigChildren: 如果为true，则信号会发送到该进程的子进程
func (that *Process) Signal(sig os.Signal, sigChildren bool) error {
	that.lock.RLock()
	defer that.lock.RUnlock()

	return that.sendSignal(sig, sigChildren)
}

// 发送多个信号到进程
// sig: 要发送的信号列表
// sigChildren: 如果为true，则信号会发送到该进程的子进程
func (that *Process) sendSignals(sigs []string, sigChildren bool) {
	that.lock.RLock()
	defer that.lock.RUnlock()

	for _, strSig := range sigs {
		sig := signals.ToSignal(strSig)
		err := that.sendSignal(sig, sigChildren)
		if err != nil {
			logger.Infof("向进程[%s]发送信号[%s]失败,err:%v", that.GetName(), strSig, err)
		}
	}
}

// sendSignal 向进程发送信号
// sig: 要发送的信号
// sigChildren: 如果为true，则信号会发送到该进程的子进程
func (that *Process) sendSignal(sig os.Signal, sigChildren bool) error {
	if that.cmd != nil && that.cmd.Process != nil {
		logger.Infof("发送信号[%s]到进程[%s]", sig, that.GetName())
		err := signals.Kill(that.cmd.Process, sig, sigChildren)
		return err
	}
	return fmt.Errorf("进程[%s]没有启动", that.GetName())
}
