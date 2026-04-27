package services

import (
	"e-shop-api/internal/pkg/logger"
	"e-shop-api/internal/pkg/utils"

	"github.com/sony/gobreaker"
	"go.uber.org/zap"
)

type NotificationService interface {
    QueueSendEmail(to, subject, body string)
}

type notificationService struct {
	emailBreaker *gobreaker.CircuitBreaker
}

func NewNotificationService() NotificationService {
    return &notificationService{
		emailBreaker: utils.NewCircuitBreaker("Email-Notification"),
	}
}

func (s *notificationService) QueueSendEmail(to, subject, body string) {
	// Run in goroutine
	utils.SafeGo(func() {
		// Send email with auto retry
		_, err := s.emailBreaker.Execute(func() (interface{}, error) {
			return nil, utils.AutoRetry(func() error {
				return utils.SendEmail(to, subject, body)
			})
		})

		if err != nil {
			// Check if circuit breaker is open
			if err == gobreaker.ErrOpenState {
				logger.L.Warn("[Email Breaker] Circuit is OPEN, skipping send", 
					zap.String("to", to))
				return
			}
			
			// Log error send email after fail all retry
			logger.L.Info("[Email Error] Failed to send", 
				zap.String("to", to), 
				zap.Error(err))
			return
		}

		// Log success send email
		logger.L.Info("[Email Success] Sent", zap.String("to", to))
	})
}
