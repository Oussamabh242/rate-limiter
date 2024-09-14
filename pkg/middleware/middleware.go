package middleware

import (
	"fmt"
	"net/http"
	"github.com/Oussamabh242/rate-limiter/pkg/bucket"
	"github.com/Oussamabh242/rate-limiter/pkg/usrmgmt"
	"time"
)

func readUserIP(r *http.Request) string {
	IPAddress := r.Header.Get("X-Real-Ip")
	if IPAddress == "" {
		IPAddress = r.Header.Get("X-Forwarded-For")
	}
	if IPAddress == "" {
		IPAddress = r.RemoteAddr
	}
	return IPAddress
}

func rateLimitedHandler(w http.ResponseWriter ,r *http.Request  ){
  w.WriteHeader(http.StatusTooManyRequests)
  w.Write([]byte("Rate limit exceeded"))
}

type Thing struct {
  buck *bucket.Bucket
}

func InitThing(ClearTime ,BlockTime time.Duration, MaxRequest int ,MaxIncativeTime time.Duration) Thing {
  b :=bucket.NewBucket(ClearTime , BlockTime ,MaxRequest ,MaxIncativeTime) 
  go b.Clean()
  return Thing{
    buck: b,
  }
  
}

func printMap(mp map[string]*bucket.X){
  for k , v := range mp {
    fmt.Println(k , "=>", v)
  }
  fmt.Println("--------------------------------------------")
}

func (t Thing) Middleware(h http.Handler) http.Handler { 
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
      
      ip := readUserIP(r)
      rateLimit:= usrmgmt.IncRequests(t.buck,ip)
      t.buck.Mu.Lock()
      printMap(t.buck.Store)
      t.buck.Mu.Unlock()
      if rateLimit {
        http.HandlerFunc(rateLimitedHandler).ServeHTTP(w , r) 
        return
      }
      h.ServeHTTP(w , r)
    })  
}


