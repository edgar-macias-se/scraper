package scraper

import "time"

type RateLimiter struct {
	tokens   chan struct{}
	rate     time.Duration
	capacity int
}

func NewRateLimiter(rate time.Duration, capacity int) *RateLimiter {
	rl := &RateLimiter{
		tokens:   make(chan struct{}, capacity),
		rate:     rate,
		capacity: capacity,
	}
	
	// Llenar burst inicial
	for i := 0; i < capacity; i++ {
		rl.tokens <- struct{}{}
	}
	
	// Recargar tokens
	go func() {
		ticker := time.NewTicker(rate)
		defer ticker.Stop()
		
		for range ticker.C {
			select {
			case rl.tokens <- struct{}{}:
			default:
			}
		}
	}()
	
	return rl
}

func (rl *RateLimiter) Wait() {
	<-rl.tokens
}
