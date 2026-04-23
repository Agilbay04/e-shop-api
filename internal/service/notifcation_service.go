package service

import (
	"e-shop-api/internal/pkg/logger"
	"e-shop-api/internal/pkg/util"

	"go.uber.org/zap"
)

type NotificationService interface {
    QueueSendEmail(to, subject, body string)
}

type notificationService struct{}

func NewNotificationService() NotificationService {
    return &notificationService{}
}

func (s *notificationService) QueueSendEmail(to, subject, body string) {
	util.SafeGo(func() {
		err := util.SendEmail(to, subject, body)
		if err != nil {
			logger.L.Info("[Email Error] Failed to send", zap.String("to", to), zap.Error(err))
			return
		}
		logger.L.Info("[Email Success] Sent", zap.String("to", to))
	})
}
