package module

import "time"

// GetTime 转换时间
func GetTime(timestamp int64) string {
	if timestamp == 0 {
		return ""
	}
	datetime := time.Unix(timestamp, 0).Format("2006-01-02 15:04:05")
	return datetime
}
