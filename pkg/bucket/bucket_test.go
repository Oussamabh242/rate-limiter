package bucket

import (
	"context"
	"sync"
	"testing"
	"time"
)

func TestBucket(t *testing.T) {
	t.Run("checking If user is blocked", func(t *testing.T) {
		ip := "192.168.1.19"
		b := NewBucket(1*SECOND, 10, time.Hour*1)
		insertUser(b, ip, false, 3, time.Now())
		expected := false
		_, got := b.CheckBlocked("192.168.1.19")
		if got != expected {
			t.Errorf("expected %t got %t", expected, got)
		}
	})

	t.Run("checking If user is not found", func(t *testing.T) {
		ip := "192.168.1.19"
		b := NewBucket(1*SECOND, 10, time.Hour*1)
		insertUser(b, ip, false, 3, time.Now())
		expected := false
		found, _ := b.CheckBlocked("192.168.1.1")
		if found != expected {
			t.Errorf("expected not found got  %t", expected)
		}
	})
	t.Run("Cleaning Inactive Users", func(t *testing.T) {
		b := NewBucket(1*SECOND, 10, time.Hour*1)
		b.Store = generateStore()
		go b.Clean()

		if 1 != 2 {
			t.Errorf("expected not found got  %v ", b.Store)
		}
	})

}

func insertUser(b *Bucket, ip string, blocked bool, reqNum int, lastReq time.Time) {
	user := X{
		Blocked:        blocked,
		RequestNumber:  reqNum,
		BlockStartTime: time.Now(),
		LastRequest:    lastReq,
	}
	b.Store[ip] = &user
}

func generateStore() map[string]*X {
	ipMap := make(map[string]*X)
	ipMap["192.168.1.1"] = &X{
		RequestNumber:  1,
		BlockStartTime: time.Now(),
		Blocked:        false,
		LastRequest:    time.Now(),
	}
	ipMap["192.168.1.2"] = &X{
		RequestNumber:  2,
		BlockStartTime: time.Now().Add(-time.Hour),
		Blocked:        true,
		LastRequest:    time.Now().Add(-time.Minute * 30),
	}
	ipMap["192.168.1.3"] = &X{
		RequestNumber:  3,
		BlockStartTime: time.Now().Add(-time.Hour * 2),
		Blocked:        false,
		LastRequest:    time.Now().Add(-time.Hour * 10),
	}
	ipMap["192.168.1.4"] = &X{
		RequestNumber:  4,
		BlockStartTime: time.Now().Add(-time.Hour * 3),
		Blocked:        true,
		LastRequest:    time.Now().Add(-time.Hour * 20),
	}
	ipMap["192.168.1.5"] = &X{
		RequestNumber:  5,
		BlockStartTime: time.Now().Add(-time.Hour * 4),
		Blocked:        false,
		LastRequest:    time.Now().Add(-time.Hour * 5),
	}
	return ipMap
}
