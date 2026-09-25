package auth

import "time"

// Expirable is used to clean up values from long-lived maps and avoid memory leaks
type Expirable[T any] struct {
	validUntil time.Time
	value      T
}

func NewExpirable[T any](value T, duration time.Duration) Expirable[T] {
	return Expirable[T]{
		validUntil: time.Now().Add(duration),
		value:      value,
	}
}

func (e Expirable[T]) IsExpired() bool {
	return time.Now().After(e.validUntil)
}

func (e Expirable[T]) GetValue() T {
	return e.value
}
