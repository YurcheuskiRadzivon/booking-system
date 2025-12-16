package notification

import (
	"fmt"

	"github.com/YurcheuskiRadzivon/booking-system/internal/models/notification"
)

type EmailHandler struct{}

func NewEmailHandler() *EmailHandler {
	return &EmailHandler{}
}

func (h *EmailHandler) Send(event notification.NotificationEvent) bool {
	fmt.Printf("   Email для %s:\n", event.Recipient)
	fmt.Printf("   Тема: %s\n", event.Subject)
	fmt.Printf("   Сообщение: %s\n", event.Message)
	return true
}

func (h *EmailHandler) Channel() notification.NotificationChannel {
	return notification.NotificationChannelEmail
}
