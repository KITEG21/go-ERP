package health

import (
	"context"
	"net/http"
	"time"

	"user_api/internal/database"

	"github.com/gin-gonic/gin"
)

// Checker defines a function that verifies readiness of an external dependency.
// It's a variable so tests can replace it with a stub.
type Checker func(ctx context.Context) error

// DBChecker is the function used by the readiness handler to check DB reachability.
// Tests can replace this to avoid touching a real DB.
var DBChecker Checker = defaultDBChecker

// defaultDBChecker pings the real database with the provided context.
func defaultDBChecker(ctx context.Context) error {
	sqlDB, err := database.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.PingContext(ctx)
}

// Liveness returns quickly if the process is running. It intentionally does not
// check external dependencies so it is suitable for a Kubernetes liveness probe.
func Liveness(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// Readiness verifies that required external dependencies are available. Currently
// it performs a short DB ping. It returns 200 when ready, or 503 when unavailable.
func Readiness(c *gin.Context) {
	// Keep the readiness check short so probes don't hang.
	ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
	defer cancel()

	if err := DBChecker(ctx); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status": "unavailable",
			"error":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ready"})
}

// RegisterRoutes registers the health endpoints on the provided gin Engine.
//
// Endpoints are intentionally available at the root (no auth) so health probes
// and external systems can reach them without the API's authentication middleware.
func RegisterRoutes(r *gin.Engine) {
	// Provide the conventional /healthz and /readyz aliases
	r.GET("/healthz", Liveness)
	r.GET("/readyz", Readiness)
}
