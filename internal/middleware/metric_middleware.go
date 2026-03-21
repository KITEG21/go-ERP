package middleware

import (
	"regexp"
	"strconv"
	"strings"
	"time"
	"user_api/internal/metrics"

	"github.com/gin-gonic/gin"
)

func MetricsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		if c.Request.URL.Path == "/metrics" {
			c.Next()
			return
		}

		start := time.Now()

		c.Next()

		duration := time.Since(start).Seconds()
		status := strconv.Itoa(c.Writer.Status())
		path := normalizePath(c.FullPath(), c.Request.URL.Path)

		metrics.HttpRequestsTotal.WithLabelValues(
			c.Request.Method,
			path,
			status,
		).Inc()

		metrics.HttpRequestDuration.WithLabelValues(
			c.Request.Method,
			path,
		).Observe(duration)
	}
}

func normalizePath(fullPath, rawPath string) string {
	if fullPath != "" {
		return fullPath
	}
	if rawPath == "" {
		return ""
	}

	// Precompile regexes
	digitsRe := regexp.MustCompile(`^\d+$`)
	uuidRe := regexp.MustCompile(`(?i)^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)
	hex24Re := regexp.MustCompile(`(?i)^[0-9a-f]{24}$`)
	hexLongRe := regexp.MustCompile(`(?i)^[0-9a-f]{16,}$`)

	parts := strings.Split(rawPath, "/")
	for i, p := range parts {
		if p == "" {
			continue
		}
		if digitsRe.MatchString(p) || uuidRe.MatchString(p) || hex24Re.MatchString(p) || hexLongRe.MatchString(p) {
			parts[i] = ":id"
		}
	}

	normalized := strings.Join(parts, "/")
	// Ensure leading slash
	if !strings.HasPrefix(normalized, "/") {
		normalized = "/" + normalized
	}
	// Collapse multiple slashes (in case)
	normalized = strings.ReplaceAll(normalized, "//", "/")
	return normalized
}
