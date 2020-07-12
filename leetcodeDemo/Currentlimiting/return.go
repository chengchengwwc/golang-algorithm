package Currentlimiting

/*
令牌桶算法：
假设每100ms生产一个令牌，按user_id/IP记录访问最近一次访问的时间戳 t_last 和令牌数，
每次请求时如果 now - last > 100ms, 增加 (now - last) / 100ms个令牌。
然后，如果令牌数 > 0，令牌数 -1 继续执行后续的业务逻辑，否则返回请求频率超限的错误码或页面。
*/

import (
	"sync"
	"time"
)

var recordMu map[string]*sync.RWMutex

func init() {
	recordMu = make(map[string]*sync.RWMutex)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

type record struct {
	last  time.Time
	token int
}

type TokenBucket struct {
	BucketSize int                // 木桶中的容量，最多可以存放多少令牌
	TokenRate  time.Duration      //多长时间生成一个令牌
	records    map[string]*record //报错的uid的访问记录
}

func NewTokenBucket(bucketSize int, tokenRate time.Duration) *TokenBucket {
	return &TokenBucket{
		BucketSize: bucketSize,
		TokenRate:  tokenRate,
		records:    make(map[string]*record),
	}
}

func (t *TokenBucket) getUid() string {
	return "127.0.0.1"
}

//获取这个uid上次访问的时间戳和令牌
func (t *TokenBucket) getRecord(uid string) *record {
	if r, ok := t.records[uid]; ok {
		return r
	}
	return &record{}
}

func (t *TokenBucket) storeRecord(uid string, r *record) {
	t.records[uid] = r
}

//验证是否保存了一个令牌
func (t *TokenBucket) validata(uid string) bool {
	rl, ok := recordMu[uid]
	if !ok {
		var mu sync.RWMutex
		rl = &mu
		recordMu[uid] = rl
	}
	rl.Lock()
	defer rl.Unlock()
	r := t.getRecord(uid)
	now := time.Now()
	if r.last.IsZero() {
		//第一次访问初始化最大令牌数
		r.last, r.token = now, t.BucketSize
	} else {
		if r.last.Add(t.TokenRate).Before(now) {
			//如果上次请求的时间超过了token rate 增加令牌，同时更新last
			r.token += max(int(now.Sub(r.last)/t.TokenRate), t.BucketSize)
			r.last = now
		}
	}
	var result bool
	if r.token > 0 {
		r.token--
		result = true
	}
	t.storeRecord(uid, r)
	return result
}

//返回是否被限流
func (t *TokenBucket) IsLimited() bool {
	return !t.validata(t.getUid())
}
