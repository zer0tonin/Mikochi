package auth

type RateLimiter struct {
}

func (r *RateLimiter) checkRateLimit(key string) bool {
	return true
}

func (r *RateLimiter) increaseRateLimit(key string) {
}

func (r *RateLimiter) resetRateLimit(key string) {

}

var rateLimiter = RateLimiter{}
