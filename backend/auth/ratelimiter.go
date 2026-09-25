package auth

import (
	"log"
	"math"
	"sync"
	"time"
)

// RateLimiter is used to limit failed login attempts on a username
type RateLimiter struct {
	mutex            sync.RWMutex
	attemptsCounters map[string]Expirable[int]
}

// Initializes a new in-memory rate limiter
func NewRateLimiter() *RateLimiter {
	return &RateLimiter{
		mutex:            sync.RWMutex{},
		attemptsCounters: map[string]Expirable[int]{},
	}
}

func (r *RateLimiter) checkRateLimit(key string) bool {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	attempts, ok := r.attemptsCounters[key]
	if !ok || attempts.IsExpired() {
		return true
	}

	return false
}

func (r *RateLimiter) increaseRateLimit(key string) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	limit, ok := r.attemptsCounters[key]
	if !ok {
		r.attemptsCounters[key] = NewExpirable(
			1,
			time.Second,
		)
	} else {
		r.attemptsCounters[key] = NewExpirable(
			limit.GetValue()+1,
			time.Duration(math.Pow(2, float64(limit.GetValue())))*time.Second,
		)
	}
}

func (r *RateLimiter) resetRateLimit(key string) {
	r.mutex.Lock()
	delete(r.attemptsCounters, key)
	r.mutex.Unlock()
}

func (r *RateLimiter) Cleanup() {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	expired := []string{}
	for key, value := range r.attemptsCounters {
		if value.IsExpired() {
			expired = append(expired, key)
		}
	}

	for _, e := range expired {
		delete(r.attemptsCounters, e)
	}

	if len(expired) != 0 {
		log.Printf("Cleaned up %d expired keys from rate limiter", len(expired))
	}
}
