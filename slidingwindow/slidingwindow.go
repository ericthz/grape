package slidingwindow

import (
	"sync"
	"time"
)

// Clock 接口用于获取当前时间，便于测试
type Clock interface {
	Now() time.Time
}

// RealClock 实现实际的时间获取
type RealClock struct{}

func (c *RealClock) Now() time.Time {
	return time.Now()
}

// SlidingWindow 滑动窗口结构体
type SlidingWindow struct {
	events []time.Time   // 存储事件时间戳
	window time.Duration // 窗口时长
	mu     sync.Mutex    // 互斥锁保证并发安全
	clock  Clock         // 时间源
}

// NewSlidingWindowWithClock 创建滑动窗口实例，可注入时间源
func NewSlidingWindowWithClock(window time.Duration, clock Clock) *SlidingWindow {
	return &SlidingWindow{
		window: window,
		clock:  clock,
		events: make([]time.Time, 0),
	}
}

// NewSlidingWindow 创建滑动窗口实例，使用实际时间源
func NewSlidingWindow(window time.Duration) *SlidingWindow {
	return NewSlidingWindowWithClock(window, &RealClock{})
}

// AddEvent 记录新事件
func (sw *SlidingWindow) AddEvent() {
	sw.mu.Lock()
	defer sw.mu.Unlock()
	now := sw.clock.Now()
	sw.events = append(sw.events, now)
	sw.cleanup(now)
}

// Count 返回当前窗口内的事件数量
func (sw *SlidingWindow) Count() int {
	sw.mu.Lock()
	defer sw.mu.Unlock()
	now := sw.clock.Now()
	sw.cleanup(now)
	return len(sw.events)
}

// cleanup 清理过期事件
func (sw *SlidingWindow) cleanup(now time.Time) {
	cutoff := now.Add(-sw.window)
	i := 0
	for ; i < len(sw.events); i++ {
		if !sw.events[i].Before(cutoff) {
			break
		}
	}
	sw.events = sw.events[i:]
}
