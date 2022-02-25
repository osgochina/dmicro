// Package backoff 提供阻塞功能
package backoff

import (
	"math"
	"time"
)

// Do 函数是 x^e乘于0.1秒的倍数，最大两分钟
func Do(attempts int) time.Duration {
	if attempts > 13 {
		return 2 * time.Minute
	}
	return time.Duration(math.Pow(float64(attempts), math.E)) * time.Millisecond * 100
}
