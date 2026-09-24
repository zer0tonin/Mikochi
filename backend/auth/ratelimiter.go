package auth

import (
	"math"
	"sync"
	"time"
)

type accessLimit struct {
	attempts    int
	nextAttempt time.Time
}

// RateLimiter is used to limit failed login attempts on a username
type RateLimiter struct {
	mutex sync.Mutex
	accessMap map[string]accessLimit
}

func (r *RateLimiter) checkRateLimit(key string) bool {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	limit, ok := r.accessMap[key]
	if !ok {
		return true
	}

	if limit.nextAttempt.After(time.Now()) {
		return false
	}

	return true
}

func (r *RateLimiter) increaseRateLimit(key string) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	limit, ok := r.accessMap[key]
	if !ok {
		r.accessMap[key] = accessLimit{
			attempts:    1,
			nextAttempt: time.Now().Add(time.Second),
		}
	} else {
		r.accessMap[key] = accessLimit{
			attempts:    limit.attempts + 1,
			nextAttempt: time.Now().Add(time.Duration(math.Pow(2, float64(limit.attempts))) * time.Second),
		}
	}
}

func (r *RateLimiter) resetRateLimit(key string) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.accessMap[key] = accessLimit{
		attempts:    0,
		nextAttempt: time.Now(),
	}
}

var rateLimiter = RateLimiter{
	accessMap: map[string]accessLimit{},
}
