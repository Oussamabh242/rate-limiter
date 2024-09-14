package usrmgmt

import (
	"rate-limiter/pkg/bucket"
	"testing"
	"time"
)

func TestUsrMgmt(t *testing.T) {
	t.Run("Restoring activity for blocked ip adress after time passed", func(t *testing.T) {
		ip := "192.168.1.19"
		b := bucket.NewBucket(1*bucket.SECOND , 1*time.Second, 3, time.Hour)
		insertUser(b, ip, true, 3)
		time.Sleep(2 * time.Second)
		RestoreActivity(b, "192.168.1.19")
		ipUser := "192.168.1.19"
		found, got := b.CheckBlocked(ipUser)
		expected := false
		if !found || !got == expected || b.Store[ip].RequestNumber != 0 {
			t.Errorf("expected user found got => %t \nexpected value to be %t got %t \n request number = %v", found, expected, got, b.Store[ip].RequestNumber)
		}
	})

	t.Run("trying to Restore activity for blocked ip adress before time passed", func(t *testing.T) {
		ip := "192.168.1.19"
		b := bucket.NewBucket(1*bucket.SECOND,1*time.Second, 10, time.Hour)
		insertUser(b, ip, true, 3)
		time.Sleep(0 * time.Second)
		RestoreActivity(b, "192.168.1.19")
		ipUser := "192.168.1.19"
		found, got := b.CheckBlocked(ipUser)
		expected := true
		if !found || !got == expected {
			t.Errorf("expected user found got => %t \nexpected value to be %t got %t ", found, expected, got)
		}
	})

	t.Run("incrementing user requests not blocked user", func(t *testing.T) {
		ip := "192.168.1.19"
		b := bucket.NewBucket(1*bucket.SECOND,1*time.Second, 10, time.Hour)
		insertUser(b, ip, false, 3)
		ratelimited := IncRequests(b, ip)
		expected := 4
		expectedRateLimited := false
		got := b.Store[ip].RequestNumber

		if got != expected || ratelimited != expectedRateLimited {
			t.Errorf("expected %d , got %d \n  %v", expected, got, b.Store[ip].RequestNumber)
		}
	})

	t.Run("incrementing user requests blocked user and restoring activity", func(t *testing.T) {
		ip := "192.168.1.19"
		b := bucket.NewBucket(1*bucket.SECOND,1*time.Second, 10, time.Hour)
		insertUser(b, ip, true, 10)

		time.Sleep(2 * time.Second)

		ratelimited := IncRequests(b, ip)
		expectedReq := 1
		gotReqNum := b.Store[ip].RequestNumber
		_, gotBlock := b.CheckBlocked(ip)
		expectedBlock := false
		if gotReqNum != expectedReq || expectedBlock != gotBlock || ratelimited != false {
			t.Errorf("expected %d requests , got %d requests \n expected blocked to be %t got %t ", expectedReq, gotReqNum, expectedBlock, gotBlock)
		}
	})

	t.Run("incrementing user requests and blocking user", func(t *testing.T) {
		ip := "192.168.1.19"
		b := bucket.NewBucket(1*bucket.SECOND,1*time.Second, 10, time.Hour)
		insertUser(b, ip, false, 10)

		time.Sleep(2 * time.Second)

		ratelimited := IncRequests(b, ip)
		expectedReq := 11
		gotReqNum := b.Store[ip].RequestNumber
		timeBlock := b.Store[ip].BlockStartTime
		_, gotBlock := b.CheckBlocked(ip)
		expectedBlock := true
		if gotReqNum != expectedReq || expectedBlock != gotBlock || ratelimited != true {
			t.Errorf("expected %d requests , got %d requests \n expected blocked to be %t got %t time blocked %v", expectedReq, gotReqNum, expectedBlock, gotBlock, timeBlock)
		}
	})

}

func insertUser(b *bucket.Bucket, ip string, blocked bool, reqNum int) {
	user := bucket.X{
		Blocked:        blocked,
		RequestNumber:  reqNum,
		BlockStartTime: time.Now(),
		LastRequest:    time.Now(),
	}
	b.Store[ip] = &user
}
