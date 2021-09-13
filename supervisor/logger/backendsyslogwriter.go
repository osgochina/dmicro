package logger

import "log/syslog"

type BackendSysLogWriter struct {
	network    string
	raddr      string
	priority   syslog.Priority
	tag        string
	logChannel chan []byte
}

// NewBackendSysLogWriter creates background syslog writer
func NewBackendSysLogWriter(network, raddr string, priority syslog.Priority, tag string) *BackendSysLogWriter {
	bs := &BackendSysLogWriter{network: network, raddr: raddr, priority: priority, tag: tag, logChannel: make(chan []byte)}
	bs.start()
	return bs
}

func (bs *BackendSysLogWriter) start() {
	go func() {
		var writer *syslog.Writer = nil
		for {
			b, ok := <-bs.logChannel
			// if channel is closed
			if !ok {
				if writer != nil {
					_ = writer.Close()
				}
				break
			}
			// if not connect to syslog, try to connect to it
			if writer == nil {
				writer, _ = syslog.Dial(bs.network, bs.raddr, bs.priority, bs.tag)
			}
			if writer != nil {
				_, _ = writer.Write(b)
			}

		}
	}()
}

// Write data to the backend syslog writer
func (bs *BackendSysLogWriter) Write(b []byte) (int, error) {
	bs.logChannel <- b
	return len(b), nil
}

// Close background write channel
func (bs *BackendSysLogWriter) Close() error {
	close(bs.logChannel)
	return nil
}
