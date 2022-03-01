// Package backoff 提供阻塞功能
package backoff

import (
	"math"
	"time"
)

// DoMul 函数是 x^e乘于0.1秒的倍数，最大两分钟
func DoMul(attempts int) time.Duration {
	if attempts > 13 {
		return 2 * time.Minute
	}
	return time.Duration(math.Pow(float64(attempts), math.E)) * time.Millisecond * 100
}

// Do 计算阻塞时间
func Do(attempts int) time.Duration {
	if attempts == 0 {
		return time.Duration(0)
	}
	return time.Duration(math.Pow(10, float64(attempts))) * time.Millisecond
}
