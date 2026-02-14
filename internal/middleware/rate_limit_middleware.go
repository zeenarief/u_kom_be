package middleware

import (
	"net/http"
	"sync"
	"time"
	"smart_school_be/internal/config"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// IPRateLimiter menyimpan rate limiter untuk setiap IP
type IPRateLimiter struct {
	ips map[string]*rate.Limiter
	mu  *sync.RWMutex
}

// NewIPRateLimiter membuat instance baru IPRateLimiter
func NewIPRateLimiter() *IPRateLimiter {
	return &IPRateLimiter{
		ips: make(map[string]*rate.Limiter),
		mu:  &sync.RWMutex{},
	}
}

// AddLimiter menambahkan atau mendapatkan rate limiter untuk IP tertentu
func (i *IPRateLimiter) AddLimiter(ip string, limiter *rate.Limiter) *rate.Limiter {
	i.mu.Lock()
	defer i.mu.Unlock()

	existingLimiter, exists := i.ips[ip]
	if !exists {
		i.ips[ip] = limiter
		return limiter
	}
	return existingLimiter
}

// GetLimiter mendapatkan rate limiter untuk IP tertentu
func (i *IPRateLimiter) GetLimiter(ip string) *rate.Limiter {
	i.mu.RLock()
	defer i.mu.RUnlock()
	return i.ips[ip]
}

// Cleanup membersihkan IP yang sudah lama tidak aktif
func (i *IPRateLimiter) Cleanup(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		i.mu.Lock()
		for ip, limiter := range i.ips {
			// Jika limiter memiliki token penuh (tidak digunakan), hapus dari map
			if limiter.Tokens() >= float64(limiter.Burst()) {
				delete(i.ips, ip)
			}
		}
		i.mu.Unlock()
	}
}

var (
	ipRateLimiter *IPRateLimiter
	limiterOnce   sync.Once
)

// getRateLimiter mendapatkan instance rate limiter (singleton)
func getRateLimiter() *IPRateLimiter {
	limiterOnce.Do(func() {
		ipRateLimiter = NewIPRateLimiter()
	})
	return ipRateLimiter
}

// RateLimitMiddleware middleware untuk rate limiting berdasarkan IP
func RateLimitMiddleware(cfg *config.Config) gin.HandlerFunc {
	// Jika rate limiting tidak dienable, return middleware kosong
	if !cfg.RateLimitEnabled {
		return func(c *gin.Context) {
			c.Next()
		}
	}

	// Konversi config ke rate.Limit dan burst
	requestsPerSecond := float64(cfg.RateLimitRequests) / cfg.RateLimitTimeWindow.Seconds()
	r := rate.Limit(requestsPerSecond)
	b := cfg.RateLimitRequests

	limiter := getRateLimiter()

	// Jalankan cleanup routine
	go limiter.Cleanup(1 * time.Hour)

	return func(c *gin.Context) {
		// Dapatkan IP address client
		ip := c.ClientIP()
		if ip == "" {
			ip = "unknown"
		}

		// Dapatkan atau buat limiter untuk IP ini
		var clientLimiter *rate.Limiter
		if existing := limiter.GetLimiter(ip); existing != nil {
			clientLimiter = existing
		} else {
			clientLimiter = rate.NewLimiter(r, b)
			limiter.AddLimiter(ip, clientLimiter)
		}

		// Cek apakah request diperbolehkan
		if !clientLimiter.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":       "Too many requests",
				"message":     "You have exceeded the rate limit. Please try again later.",
				"retry_after": time.Now().Add(time.Duration(1/requestsPerSecond) * time.Second).Format(time.RFC1123),
			})
			c.Abort()
			return
		}

		// Set header rate limit info (opsional, untuk client awareness)
		c.Writer.Header().Set("X-RateLimit-Limit", string(rune(cfg.RateLimitRequests)))
		c.Writer.Header().Set("X-RateLimit-Remaining", string(rune(int(clientLimiter.Tokens()))))
		resetTime := time.Now().Add(time.Duration(1/requestsPerSecond) * time.Second)
		c.Writer.Header().Set("X-RateLimit-Reset", resetTime.Format(time.RFC1123))

		c.Next()
	}
}

// Alternative: SimpleRateLimitMiddleware (Lebih sederhana, tanpa dependency external)
func SimpleRateLimitMiddleware(cfg *config.Config) gin.HandlerFunc {
	if !cfg.RateLimitEnabled {
		return func(c *gin.Context) {
			c.Next()
		}
	}

	type ipInfo struct {
		count     int
		lastReset time.Time
	}

	var (
		ips = make(map[string]*ipInfo)
		mu  sync.Mutex
	)

	// Reset hit count setiap time window
	ticker := time.NewTicker(cfg.RateLimitTimeWindow)
	go func() {
		for range ticker.C {
			mu.Lock()
			for ip, info := range ips {
				if time.Since(info.lastReset) >= cfg.RateLimitTimeWindow {
					delete(ips, ip)
				}
			}
			mu.Unlock()
		}
	}()

	return func(c *gin.Context) {
		ip := c.ClientIP()
		if ip == "" {
			ip = "unknown"
		}

		now := time.Now()

		mu.Lock()

		info, exists := ips[ip]
		if !exists {
			info = &ipInfo{
				count:     0,
				lastReset: now,
			}
			ips[ip] = info
		}

		// Reset counter jika sudah melewati time window
		if now.Sub(info.lastReset) > cfg.RateLimitTimeWindow {
			info.count = 0
			info.lastReset = now
		}

		info.count++
		count := info.count

		mu.Unlock()

		if count > cfg.RateLimitRequests {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":       "Too many requests",
				"message":     "Rate limit exceeded",
				"retry_after": now.Add(cfg.RateLimitTimeWindow).Format(time.RFC1123),
			})
			c.Abort()
			return
		}

		// Set headers
		c.Writer.Header().Set("X-RateLimit-Limit", string(rune(cfg.RateLimitRequests)))
		c.Writer.Header().Set("X-RateLimit-Remaining", string(rune(cfg.RateLimitRequests-count)))
		resetTime := info.lastReset.Add(cfg.RateLimitTimeWindow)
		c.Writer.Header().Set("X-RateLimit-Reset", resetTime.Format(time.RFC1123))

		c.Next()
	}
}
