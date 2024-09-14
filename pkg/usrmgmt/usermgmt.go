package usrmgmt

import (
	"github.com/Oussamabh242/rate-limiter/pkg/bucket"
	"time"
)

func RestoreActivity(b *bucket.Bucket, IPAddress string) {
	user, found := b.Store[IPAddress]
  if !found {
    return
  }
	start := user.BlockStartTime.Add(b.BlockTime)
	if user.Blocked && start.Before(time.Now()) {
		user.Blocked = false
		user.RequestNumber = 0
	}
}

func createUser(b *bucket.Bucket, ip string) {
	b.Store[ip] = &bucket.X{
		Blocked:       false,
		RequestNumber: 0,
	}
}

func IncRequests(b *bucket.Bucket, ip string) bool {

	//creating user if first request
  b.Mu.Lock()
	_, found := b.Store[ip]
	if !found {
		createUser(b, ip)
	}

  b.Store[ip].LastRequest = time.Now()


	// trying to restore user actitvity in case he is blocked
	if b.Store[ip].Blocked {
		RestoreActivity(b, ip)
	}
	// in case he hit the max limit and (keep blocking the user)
	if b.Store[ip].RequestNumber+1 > b.MaxReq {
		b.Store[ip].BlockStartTime = time.Now()
		b.Store[ip].Blocked = true
		b.Store[ip].RequestNumber++
    b.Mu.Unlock()
		return true

	} else {
		b.Store[ip].RequestNumber++
	}
  b.Mu.Unlock()
	return false

}
