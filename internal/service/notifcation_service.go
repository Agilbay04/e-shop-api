package service

import (
	"e-shop-api/internal/pkg/util"
	"log"
)

type NotificationService interface {
    QueueSendEmail(to, subject, body string)
}

type notificationService struct{}

func NewNotificationService() NotificationService {
    return &notificationService{}
}

func (s *notificationService) QueueSendEmail(to, subject, body string) {
    // Use goroutine to send email
    util.SafeGo(func() {
        err := util.SendEmail(to, subject, body)
        if err != nil {
            log.Printf("[Email Error] Failed to send to %s: %v", to, err)
            return
        }
        log.Printf("[Email Success] Sent to %s", to)
    })
}
