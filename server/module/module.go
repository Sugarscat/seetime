package module

import (
	"fmt"
	"time"

	"github.com/gorhill/cronexpr"
)

// LeakyBucket 漏桶算法的桶定义
type LeakyBucket struct {
	capacity     float64   // 桶容量
	rate         float64   // 漏水速率
	water        float64   // 当前水量
	lastLeakTime time.Time // 上次漏水时间
}

// GetTime 转换时间
func GetTime(timestamp int64) string {
	if timestamp == 0 {
		return ""
	}
	datetime := time.Unix(timestamp, 0).Format("2006-01-02 15:04:05")
	return datetime
}

/* NewLeakyBucket 漏桶算法 限制 api 请求次数 */
func NewLeakyBucket(capacity, rate float64) *LeakyBucket {
	return &LeakyBucket{
		capacity:     capacity,
		rate:         rate,
		water:        0,
		lastLeakTime: time.Now(),
	}
}

func (b *LeakyBucket) AddWater(amount float64) bool {
	// 先漏水
	b.Leak()
	// 再加水
	if b.water+amount <= b.capacity {
		b.water += amount
		return true // 添加成功
	} else {
		return false // 添加失败
	}
}

func (b *LeakyBucket) Leak() {
	now := time.Now()
	elapsed := now.Sub(b.lastLeakTime).Seconds() // 计算距离上次漏水时间
	b.water = Max(b.water-elapsed*b.rate, 0)     // 漏掉一定量的水
	b.lastLeakTime = now                         // 更新上次漏水时间
}

func Max(x, y float64) float64 {
	if x > y {
		return x
	}
	return y
}

// ParsingCrontab 解析定时器的表达式，得到下次执行时间
func ParsingCrontab(last int64, cron string) int64 {
	defer func() { // 捕获异常，避免任务执行错误，导致整个系统挂掉
		if err := recover(); err != nil {
			fmt.Println("无法解析 Cron")
		}
	}()
	expr := cronexpr.MustParse(cron)
	return expr.Next(time.Unix(last, 0)).Unix()
}
