package ratelimiter

import (
	"sync"
	"time"
)

//Safe map so that only one client can access it a time
type limitersMap struct {
	ips map[string]limiter
	mu  sync.Mutex
}

//set the map element by first locking the map
func (m *limitersMap) set(key string, val limiter) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.ips[key] = val
}

func (m *limitersMap) Get(key string) limiter {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.ips[key]
}

func (m *limitersMap) release(ip string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.ips, ip)
}

//Map element to store values related to rate limitter
type limiter struct {
	lastRequest time.Time
	capacity    int
	Colldowned  bool
}
