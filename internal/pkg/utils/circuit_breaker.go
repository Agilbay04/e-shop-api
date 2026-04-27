package utils

import (
	"e-shop-api/internal/pkg/logger"
	"runtime"
	"strings"
	"time"

	"github.com/sony/gobreaker"
	"go.uber.org/zap"
)

func NewCircuitBreaker(name string) *gobreaker.CircuitBreaker {
	// Config circuit breaker
	maxReq := GetEnvInt("CB_MAX_REQUESTS", "3")
	interval := GetEnvTime("CB_INTERVAL", "5s")
	timeout := GetEnvTime("CB_TIMEOUT", "30s")
	threshold := GetEnvInt("CB_THRESHOLD", "3")

	settings := gobreaker.Settings{
		Name:        name,
		MaxRequests: uint32(maxReq),
		Interval:    interval,
		Timeout:     timeout,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			return counts.ConsecutiveFailures > uint32(threshold)
		},
		OnStateChange: func(name string, from gobreaker.State, to gobreaker.State) {
			logger.L.Warn("Circuit Breaker State Change",
				zap.String("name", name),
				zap.String("from", from.String()),
				zap.String("to", to.String()),
			)
		},
	}
	return gobreaker.NewCircuitBreaker(settings)
}

// Auto retry mechanism
func AutoRetry(fn func() error) error {
	attempts := GetEnvInt("RETRY_ATTEMPTS", "3")
	delay := GetEnvTime("RETRY_DELAY", "2s")

	pc, _, _, ok := runtime.Caller(1)
	actionName := "unknown"
	if ok {
		fullFuncName := runtime.FuncForPC(pc).Name()
		parts := strings.Split(fullFuncName, "/")
		actionName = parts[len(parts)-1]
	}
	
	return RetryHelper(actionName, attempts, delay, fn)
}

// Auto retry helper
func RetryHelper(actionName string, attempts int, sleep time.Duration, fn func() error) error {
	var err error
	for i := range attempts {
		if err = fn(); err == nil {
			return nil
		}
		
		logger.L.Warn("Retrying action...", 
			zap.String("func", actionName),
			zap.Int("attempt", i+1), 
			zap.Error(err),
		)

		if i < attempts-1 {
			time.Sleep(sleep)
		}
	}
	return err
}