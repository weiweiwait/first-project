package rate_liniter

import (
	"sync"
	"time"
)

type TokenBucket struct {
	capacity int64 //桶容量(每秒钟处理请求的数量)
	token    int64 //令牌数量
	rate     int64 //每秒钟生成令牌数量
	time     time.Time
	mu       sync.Mutex
}

func min(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

func NewTokenBucket(capacity int64, token int64, rate int64) *TokenBucket {
	return &TokenBucket{
		capacity: capacity,
		token:    token,
		rate:     rate,
		time:     time.Now(),
	}
}

func (b *TokenBucket) TryAcquired() bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	//计算生成令牌数量
	createToken := int64(time.Now().Second()-b.time.Second()) * b.rate
	b.token = min(b.token+createToken, b.capacity)
	b.time = time.Now()
	if b.token > 0 {
		b.token--
		return true
	} else {
		return false
	}

}
