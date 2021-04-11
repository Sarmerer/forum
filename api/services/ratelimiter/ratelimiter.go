package ratelimiter

import (
	"time"
)

var Limiter *limitersMap = &limitersMap{ips: make(map[string]limiter)}

// LimitExceeded : checks if the limit has been exceeded, by checking against the
// time of the first request. If the first request was made more than the time limit time agoa
// then the first request time is updated to latest request's time
func LimitExceeded(ip string, capacity int, timeLimit time.Duration, cooldown time.Duration) bool {
	requestor := Limiter.Get(ip)
	if (requestor == limiter{}) {
		l := limiter{time.Now(), capacity, false}
		Limiter.set(ip, l)
		requestor = l
	}
	if requestor.Colldowned {
		return true
	}
	requestor.lastRequest = time.Now()
	if requestor.capacity <= 0 && !requestor.Colldowned {
		requestor.Colldowned = true
		go func() {
			cd := timeLimit
			if cooldown > 0*time.Second {
				cd = cooldown
			}
			time.Sleep(cd)
			Limiter.release(ip)
		}()
		Limiter.set(ip, requestor)
		return true
	}
	requestor.capacity--
	Limiter.set(ip, requestor)
	return false
}
