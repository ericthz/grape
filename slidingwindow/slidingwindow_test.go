package slidingwindow

import (
	"testing"
	"time"
)

/*
	关键点：​​
	​​队列维护：​​ 使用切片存储事件时间戳，保证事件按时间顺序排列
	​​过期清理：​​ 每次添加事件或查询数量时，清理早于窗口起始时间的事件
	​​并发安全：​​ 使用互斥锁（sync.Mutex）确保多线程安全
	​​时间处理：​​ 使用!time.Before()判断事件是否在窗口内，确保时间窗口左闭右闭
	​​可测试性：​​ 通过Clock接口支持注入模拟时间，便于单元测试
*/

func TestSlidingWindow(t *testing.T) {
	// 创建5秒滑动窗口
	window := NewSlidingWindow(5 * time.Second)

	// 记录事件
	window.AddEvent()
	window.AddEvent()
	window.AddEvent()

	go func() {
		// 模拟事件添加
		for i := 0; i < 10; i++ {
			window.AddEvent()
			time.Sleep(1 * time.Second)
		}
	}()

	// 获取当前窗口内事件数
	t.Log("Current count:", window.Count()) // 输出: 3

	// 等待6秒后，事件过期
	time.Sleep(6 * time.Second)
	t.Log("Count after 6 seconds:", window.Count()) // 输出: 0
}

type MockClock struct {
	now time.Time
}

func (m *MockClock) Now() time.Time {
	return m.now
}

func TestMockSlidingWindow(t *testing.T) {
	clock := &MockClock{}
	sw := NewSlidingWindowWithClock(5*time.Second, clock)

	clock.now = time.Unix(0, 0) // 初始时间设为0
	sw.AddEvent()               // 事件时间 0
	sw.AddEvent()               // 事件时间 0
	sw.AddEvent()               // 事件时间 0

	if count := sw.Count(); count != 3 {
		t.Errorf("Expected 3, got %d", count)
	}

	clock.now = time.Unix(6, 0) // 时间推进到6秒
	if count := sw.Count(); count != 0 {
		t.Errorf("Expected 0, got %d", count)
	}
}
