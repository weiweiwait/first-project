package rate_liniter

import (
	"context"
	"fmt"
	"sync"
	"time"
)

const defaultWindowCount = 10

type SlidingWindow struct {
	windowSize   int64
	qps          int32           //窗口内的最大qps
	windowCount  int64           //子窗口个数
	windowsArray []int32         //子窗口数组
	count        int32           //当前窗口内请求数
	ctx          context.Context //控制限流协程的上下文
	mu           sync.Mutex
}

func NewSlidingWindow(windowSize int64, qps int32, ctx context.Context) *SlidingWindow {
	window := &SlidingWindow{
		windowSize:   windowSize,
		qps:          qps,
		ctx:          ctx,
		windowsArray: make([]int32, defaultWindowCount),
		windowCount:  defaultWindowCount,
	}
	go window.start()
	return window
}

func (w *SlidingWindow) start() {
	for {
		w.mu.Lock()
		select {
		case <-w.ctx.Done():
			fmt.Println("sliding window end")
		default:
			time.Sleep(time.Duration(w.windowSize/w.windowCount) * time.Millisecond)
			w.windowsArray = append(w.windowsArray, 0)
			if len(w.windowsArray) > int(w.windowCount) {
				w.count = w.count - w.windowsArray[0]
				w.windowsArray = w.windowsArray[1:]
			}
		}
		w.mu.Unlock()
	}
}

func (w *SlidingWindow) TryQuery() bool {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.count >= w.qps {
		return false
	}
	w.windowsArray[int(w.windowCount)-1]++
	w.count++
	return true
}
