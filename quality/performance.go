package quality

import (
  "sync"
  "time"
  "golang.org/x/time/rate"
)

// ConnectionConfig holds configuration for connection management
type ConnectionConfig struct {
  MaxConns      int           // Maximum concurrent connections
  IdleTimeout   time.Duration // Time before idle connections are closed
  WriteTimeout  time.Duration // Timeout for write operations
}

// ConnectionManager manages concurrent connections
type ConnectionManager struct {
  config     ConnectionConfig
  semaphore  chan struct{}
  current    int
  currentMu  sync.Mutex
}

// NewConnectionManager creates a new connection manager
func NewConnectionManager(config ConnectionConfig) *ConnectionManager {
  return &ConnectionManager{
    config:    config,
    semaphore: make(chan struct{}, config.MaxConns),
  }
}

// Acquire reserves a connection slot
func (cm *ConnectionManager) Acquire() {
  cm.semaphore <- struct{}{}
  cm.currentMu.Lock()
  cm.current++
  cm.currentMu.Unlock()
}

// Release frees a connection slot
func (cm *ConnectionManager) Release() {
  cm.currentMu.Lock()
  cm.current--
  cm.currentMu.Unlock()
  <-cm.semaphore
}

// Current returns current connection count
func (cm *ConnectionManager) Current() int {
  cm.currentMu.Lock()
  defer cm.currentMu.Unlock()
  return cm.current
}

// Current returns current connection count
func (cm *ConnectionManager) Close() error {
  return nil
}

// RateLimiter implements connection rate limiting
type RateLimiter struct {
  limiter       *rate.Limiter
  visitors      map[string]*rate.Limiter
  visitorsMu    sync.Mutex
  burst         int
  cleanupTicker *time.Ticker
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(limit rate.Limit, burst int, cleanupInterval time.Duration) *RateLimiter {
  rl := &RateLimiter{
    limiter:  rate.NewLimiter(limit, burst),
    visitors: make(map[string]*rate.Limiter),
    burst:    burst,
  }

  // Start cleanup goroutine
  rl.cleanupTicker = time.NewTicker(cleanupInterval)
  go func() {
    for range rl.cleanupTicker.C {
      rl.cleanupVisitors()
    }
  }()

  return rl
}

// GetVisitor gets or creates a rate limiter for a visitor
func (rl *RateLimiter) GetVisitor(ip string) *rate.Limiter {
  rl.visitorsMu.Lock()
  defer rl.visitorsMu.Unlock()

  limiter, exists := rl.visitors[ip]
  if !exists {
    limiter = rate.NewLimiter(rl.limiter.Limit(), rl.burst)
    rl.visitors[ip] = limiter
  }

  return limiter
}

// Allow checks if a visitor is allowed to connect
func (rl *RateLimiter) Allow() bool {
  return rl.limiter.Allow()
}

// cleanupVisitors removes old visitors
func (rl *RateLimiter) cleanupVisitors() {
  rl.visitorsMu.Lock()
  defer rl.visitorsMu.Unlock()

  for ip := range rl.visitors {
    // In real implementation, you might track last activity
    delete(rl.visitors, ip)
  }
}

// Stop cleans up the rate limiter
func (rl *RateLimiter) Stop() {
  rl.cleanupTicker.Stop()
}
