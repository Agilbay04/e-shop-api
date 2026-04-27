package handlers

import (
	"e-shop-api/internal/dtos"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type HealthHandler struct {
	db  *gorm.DB
	rdb *redis.Client
}

func NewHealthHandler(db *gorm.DB, rdb *redis.Client) *HealthHandler {
	return &HealthHandler{
		db:  db,
		rdb: rdb,
	}
}

func (h *HealthHandler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, dtos.HealthResponse{
		Status:   "healthy",
		Service:  "e-shop-api",
		Version:  "1.0.0",
		Details:  nil,
	})
}

func (h *HealthHandler) Readiness(c *gin.Context) {
	checks := make(map[string]any)

	dbSQL, err := h.db.DB()
	if err != nil {
		checks["database"] = dtos.ComponentStatus{
			Status:  "unhealthy",
			Message: "failed to get database connection",
		}
		c.JSON(http.StatusServiceUnavailable, dtos.ReadinessCheck{
			Status: "not_ready",
			Checks: checks,
		})
		return
	}

	if err := dbSQL.Ping(); err != nil {
		checks["database"] = dtos.ComponentStatus{
			Status:  "unhealthy",
			Message: err.Error(),
		}
		c.JSON(http.StatusServiceUnavailable, dtos.ReadinessCheck{
			Status: "not_ready",
			Checks: checks,
		})
		return
	}
	checks["database"] = dtos.ComponentStatus{Status: "healthy"}

	if err := h.rdb.Ping(c.Request.Context()).Err(); err != nil {
		checks["redis"] = dtos.ComponentStatus{
			Status:  "unhealthy",
			Message: err.Error(),
		}
		c.JSON(http.StatusServiceUnavailable, dtos.ReadinessCheck{
			Status: "not_ready",
			Checks: checks,
		})
		return
	}
	checks["redis"] = dtos.ComponentStatus{Status: "healthy"}

	c.JSON(http.StatusOK, dtos.ReadinessCheck{
		Status: "ready",
		Checks: checks,
	})
}