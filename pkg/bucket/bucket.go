package bucket

import (
	"fmt"
	"sync"
	"time"
)

const (
	SECOND = time.Second
	MINUTE = time.Minute
	HOUR   = time.Hour
)

type X struct {
	RequestNumber  int
	BlockStartTime time.Time
	Blocked        bool
	LastRequest    time.Time
}

type Bucket struct {
	Store            map[string]*X
	ClearTime        time.Duration
	BlockTime        time.Duration
	MaxReq           int
	MaxIncactiveTime time.Duration
	Mu               sync.Mutex
}

func NewBucket(Period, BlockTime time.Duration, MaxReq int, MaxIncactiveTime time.Duration) *Bucket {
	var store map[string]*X = map[string]*X{}
	b := Bucket{
		BlockTime:        BlockTime,
		MaxReq:           MaxReq,
		Store:            store,
		ClearTime:        Period,
		MaxIncactiveTime: MaxIncactiveTime,
	}
	return &b
}

func (b *Bucket) CheckBlocked(ip string) (found bool, value bool) {
	b.Mu.Lock()
	val, found := b.Store[ip]
	b.Mu.Unlock()
	if !found {
		value = true
		return found, value
	}
	return found, val.Blocked
}

func Kill(lastReq time.Time, MaxInactive time.Duration) bool {
	return time.Since(lastReq) > MaxInactive
}

func (b *Bucket) Clean() {
	ticker := time.NewTicker(b.ClearTime)
	for {
		fmt.Println("+++WORKING+++")
		select {
		case <-ticker.C:
			for ip, user := range b.Store {
				b.Mu.Lock()
				if Kill(user.LastRequest, b.MaxIncactiveTime) {
					delete(b.Store, ip)
				}
				b.Mu.Unlock()
			}

		}
	}

}
