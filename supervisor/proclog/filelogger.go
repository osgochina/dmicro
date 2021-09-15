package proclog

import (
	"fmt"
	"github.com/gogf/gf/errors/gerror"
	"os"
	"sync"
)

// FileLogger 写入stdout/stderr到文件
type FileLogger struct {
	//文件名
	name string
	// 日志最大长度
	maxSize int64
	// 保留的日志份数
	backups int
	// 每个文件的长度
	fileSize int64
	//文件句柄
	file *os.File
	//logEventEmitter LogEventEmitter
	locker sync.Locker
}

func NewFileLogger(fileName string, maxSize int64, backups int, locker sync.Locker) *FileLogger {
	logger := &FileLogger{
		name:     fileName,
		maxSize:  maxSize,
		backups:  backups,
		fileSize: 0,
		file:     nil,
		locker:   locker,
	}
	_ = logger.openFile(false)
	return logger
}

//日志文件写入
func (that *FileLogger) Write(p []byte) (int, error) {
	that.locker.Lock()
	defer that.locker.Unlock()

	n, err := that.file.Write(p)

	if err != nil {
		return n, err
	}
	//that.logEventEmitter.emitLogEvent(string(p))
	that.fileSize += int64(n)
	if that.fileSize >= that.maxSize {
		fileInfo, errStat := os.Stat(that.name)
		if errStat == nil {
			that.fileSize = fileInfo.Size()
		} else {
			return n, errStat
		}
	}
	if that.fileSize >= that.maxSize {
		_ = that.Close()
		that.backupFiles()
		_ = that.openFile(true)
	}
	return n, err
}

// Close 关闭文件
func (that *FileLogger) Close() error {
	if that.file != nil {
		err := that.file.Close()
		that.file = nil
		return err
	}
	return nil
}

func (that *FileLogger) SetPid(pid int) {
	// NOTHING TO DO
}

// ClearCurLogFile 清除当前日志文件
func (that *FileLogger) ClearCurLogFile() error {
	that.locker.Lock()
	defer that.locker.Unlock()

	return that.openFile(true)
}

// ClearAllLogFile 清除所有日志文件，包括保留的备份
func (that *FileLogger) ClearAllLogFile() error {
	that.locker.Lock()
	defer that.locker.Unlock()

	for i := that.backups; i > 0; i-- {
		logFile := fmt.Sprintf("%s.%d", that.name, i)
		_, err := os.Stat(logFile)
		if err == nil {
			err = os.Remove(logFile)
			if err != nil {
				return err
			}
		}
	}
	err := that.openFile(true)
	if err != nil {
		return err
	}
	return nil
}

// ReadLog 读取日志
func (that *FileLogger) ReadLog(offset int64, length int64) (string, error) {
	if offset < 0 && length != 0 {
		return "", gerror.New("BAD_ARGUMENTS")
	}
	if offset >= 0 && length < 0 {
		return "", gerror.New("BAD_ARGUMENTS")
	}
	that.locker.Lock()
	defer that.locker.Unlock()
	f, err := os.Open(that.name)
	if err != nil {
		return "", gerror.New("FAILED")
	}
	defer func(f *os.File) {
		_ = f.Close()
	}(f)
	statInfo, err := f.Stat()
	if err != nil {
		return "", gerror.New("FAILED")
	}
	fileLen := statInfo.Size()
	if offset < 0 { // offset < 0 && length == 0
		offset = fileLen + offset
		if offset < 0 {
			offset = 0
		}
		length = fileLen - offset
	} else if length == 0 { // offset >= 0 && length == 0
		if offset > fileLen {
			return "", nil
		}
		length = fileLen - offset
	} else { // offset >= 0 && length > 0

		// if the offset exceeds the length of file
		if offset >= fileLen {
			return "", nil
		}

		// compute actual bytes should be read

		if offset+length > fileLen {
			length = fileLen - offset
		}
	}
	b := make([]byte, length)
	n, err := f.ReadAt(b, offset)
	if err != nil {
		return "", gerror.New("FAILED")
	}
	return string(b[:n]), nil
}

// ReadTailLog 读取尾部日志
func (that *FileLogger) ReadTailLog(offset int64, length int64) (string, int64, bool, error) {
	if offset < 0 {
		return "", offset, false, fmt.Errorf("offset should not be less than 0")
	}
	if length < 0 {
		return "", offset, false, fmt.Errorf("length should be not be less than 0")
	}
	that.locker.Lock()
	defer that.locker.Unlock()

	// open the file
	f, err := os.Open(that.name)
	if err != nil {
		return "", 0, false, err
	}

	defer func(f *os.File) {
		_ = f.Close()
	}(f)

	// get the length of file
	statInfo, err := f.Stat()
	if err != nil {
		return "", 0, false, err
	}

	fileLen := statInfo.Size()

	// check if offset exceeds the length of file
	if offset >= fileLen {
		return "", fileLen, true, nil
	}

	// get the length
	if offset+length > fileLen {
		length = fileLen - offset
	}

	b := make([]byte, length)
	n, err := f.ReadAt(b, offset)
	if err != nil {
		return "", offset, false, err
	}
	return string(b[:n]), offset + int64(n), false, nil

}

//打开要写入的日志文件
func (that *FileLogger) openFile(trunc bool) error {
	if that.file != nil {
		_ = that.file.Close()
	}
	fileInfo, err := os.Stat(that.name)
	if trunc || err != nil {
		that.file, err = os.Create(that.name)
		that.fileSize = 0
	} else {
		that.fileSize = fileInfo.Size()
		that.file, err = os.OpenFile(that.name, os.O_RDWR|os.O_APPEND, 0666)
	}

	if err != nil {
		fmt.Printf("Fail to open log file --%s-- with error %v\n", that.name, err)
	}
	return err
}

// 备份日志文件
func (that *FileLogger) backupFiles() {
	for i := that.backups - 1; i > 0; i-- {
		src := fmt.Sprintf("%s.%d", that.name, i)
		dest := fmt.Sprintf("%s.%d", that.name, i+1)
		if _, err := os.Stat(src); err == nil {
			_ = os.Rename(src, dest)
		}
	}
	dest := fmt.Sprintf("%s.1", that.name)
	_ = os.Rename(that.name, dest)
}
