package middleware

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// PrometheusMiddleware returns a Gin middleware that observes request count
// and duration. It must NOT register collectors itself; call RegisterMetrics()
// once at startup (e.g. in your server setup) before using the middleware.
func PrometheusMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip metrics endpoint itself to avoid recursion and noise
		if c.FullPath() == "/metrics" || c.Request.URL.Path == "/metrics" {
			c.Next()
			return
		}

		start := time.Now()
		c.Next() // process request

		// Obtain a stable route label — FullPath returns the registered route
		// pattern (e.g. /api/v1/users/:id). Fall back to raw path for unmatched routes.
		path := c.FullPath()
		if path == "" {
			path = c.Request.URL.Path
		}

		method := c.Request.Method
		status := strconv.Itoa(c.Writer.Status())
		duration := time.Since(start).Seconds()

		// Defensive nil-checks (shouldn't be nil, but safe)
		if RequestDuration != nil {
			RequestDuration.WithLabelValues(method, path, status).Observe(duration)
		}
		if RequestCounter != nil {
			RequestCounter.WithLabelValues(method, path, status).Inc()
		}
	}
}
