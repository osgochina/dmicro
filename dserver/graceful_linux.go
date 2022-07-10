// +build linux

package dserver

// GraceSignal 监听信号
func (that *graceful) graceSignal() {
	// 单进程模型
	if that.dServer.procModel == ProcessModelSingle {
		that.onStart()
		that.graceSingle()
		return
	}
	// 多进程模型
	if that.dServer.procModel == ProcessModelMulti {
		that.graceMultiSignal()
		return
	}
}

// 单进程模式平滑重启
func (that *graceful) graceSingle() {
	signal.Notify(
		that.signal,
		syscall.SIGINT,
		syscall.SIGQUIT,
		syscall.SIGKILL,
		syscall.SIGTERM,
		syscall.SIGABRT,
		syscall.SIGUSR1,
		syscall.SIGUSR2,
	)
	pid := os.Getpid()
	for {
		sig := <-that.signal
		logger.Printf(`进程:%d,收到信号: %s`, pid, sig.String())
		switch sig {
		// 强制关闭服务
		case syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGABRT:
			signal.Reset(syscall.SIGINT, syscall.SIGQUIT, syscall.SIGKILL, syscall.SIGABRT, syscall.SIGTERM)
			// 强制关闭的时候，设置超时时间为1秒，表示1秒后强制结束
			that.shutdownSingle(time.Second)
			continue
		// 平滑的关闭服务
		case syscall.SIGQUIT:
			signal.Reset(syscall.SIGINT, syscall.SIGQUIT, syscall.SIGKILL, syscall.SIGABRT, syscall.SIGTERM)
			// 平滑重启的时候使用默认超时时间，能够等待业务处理完毕，优雅的结束
			that.shutdownSingle()
			continue
		// 平滑重启服务
		case syscall.SIGUSR2:
			that.rebootSingle()
			continue
		default:
		}
	}
}

// MasterWorker模式平滑重启
func (that *graceful) graceMultiSignal() {
	pid := os.Getpid()
	if that.isChild() {
		signal.Notify(
			that.signal,
			syscall.SIGINT,
			syscall.SIGQUIT,
			syscall.SIGKILL,
			syscall.SIGTERM,
			syscall.SIGABRT,
		)
		for {
			sig := <-that.signal
			logger.Printf(`进程:%d,收到信号: %s`, pid, sig.String())
			switch sig {
			// 强制关闭服务
			case syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGABRT:
				signal.Reset(syscall.SIGINT, syscall.SIGQUIT, syscall.SIGKILL, syscall.SIGABRT, syscall.SIGTERM)
				// 强制关闭的时候，设置超时时间为1秒，表示1秒后强制结束
				that.shutdownMultiChild(time.Second)
				continue
			// 平滑的关闭服务
			case syscall.SIGQUIT:
				signal.Reset(syscall.SIGINT, syscall.SIGQUIT, syscall.SIGKILL, syscall.SIGABRT, syscall.SIGTERM)
				// 平滑重启的时候使用默认超时时间，能够等待业务处理完毕，优雅的结束
				that.shutdownMultiChild()
				continue
			default:
			}
		}
	} else {

		signal.Notify(
			that.signal,
			syscall.SIGINT,
			syscall.SIGQUIT,
			syscall.SIGKILL,
			syscall.SIGTERM,
			syscall.SIGABRT,
			syscall.SIGUSR1,
			syscall.SIGUSR2,
		)
		for {
			sig := <-that.signal
			logger.Printf(`进程:%d,收到信号: %s`, pid, sig.String())
			switch sig {
			// 关闭服务
			case syscall.SIGINT, syscall.SIGKILL, syscall.SIGABRT, syscall.SIGTERM:
				signal.Reset(syscall.SIGINT, syscall.SIGQUIT, syscall.SIGKILL, syscall.SIGABRT, syscall.SIGTERM)
				// 强制关闭的时候，设置超时时间为1秒，表示1秒后强制结束
				that.shutdownMultiMaster()
				continue
			//优化的关闭服务
			case syscall.SIGQUIT:
				signal.Reset(syscall.SIGINT, syscall.SIGQUIT, syscall.SIGKILL, syscall.SIGABRT, syscall.SIGTERM)
				that.shutdownMultiMaster()
				continue
			// 平滑重启服务
			case syscall.SIGUSR2:
				that.rebootMulti()
				continue
			default:
			}
		}
	}
}
