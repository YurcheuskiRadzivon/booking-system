package notification

import (
	"fmt"

	"github.com/YurcheuskiRadzivon/booking-system/internal/models/notification"
)

type SMSHandler struct{}

func NewSMSHandler() *SMSHandler {
	return &SMSHandler{}
}

func (h *SMSHandler) Send(event notification.NotificationEvent) bool {
	fmt.Printf(" SMS для %s:\n", event.Recipient)
	fmt.Printf("  Сообщение: %s\n", event.Message)
	return true
}

func (h *SMSHandler) Channel() notification.NotificationChannel {
	return notification.NotificationChannelSMS
}
