package utils

import "strconv"

// GetBytes 通过可识别字符串，返回数字容量
//
//	logSize=1MB
//	logSize=1GB
//	logSize=1KB
//	logSize=1024
//
func GetBytes(value string, defValue int) int {

	if len(value) > 2 {
		lastTwoBytes := value[len(value)-2:]
		if lastTwoBytes == "MB" {
			return toInt(value[:len(value)-2], 1024*1024, defValue)
		} else if lastTwoBytes == "GB" {
			return toInt(value[:len(value)-2], 1024*1024*1024, defValue)
		} else if lastTwoBytes == "KB" {
			return toInt(value[:len(value)-2], 1024, defValue)
		}
		return toInt(value, 1, defValue)
	}
	return defValue
}

func toInt(s string, factor int, defValue int) int {
	i, err := strconv.Atoi(s)
	if err == nil {
		return i * factor
	}
	return defValue
}
