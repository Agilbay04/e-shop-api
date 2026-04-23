package util

import (
	"e-shop-api/internal/pkg/logger"

	"go.uber.org/zap"
)

// SafeGo runs a goroutine and recovers from panics
func SafeGo(task func()) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				logger.Log.Info("[Goroutine Panic] Recovered", zap.Any("error", r))
			}
		}()

		task()
	}()
}