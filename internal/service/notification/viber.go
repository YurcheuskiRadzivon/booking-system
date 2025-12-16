package notification

import (
	"fmt"

	"github.com/YurcheuskiRadzivon/booking-system/internal/models/notification"
)

type ViberHandler struct{}

func NewViberHandler() *ViberHandler {
	return &ViberHandler{}
}

func (h *ViberHandler) Send(event notification.NotificationEvent) bool {
	fmt.Printf("💬 Viber для %s:\n", event.Recipient)
	fmt.Printf("   Сообщение: %s\n", event.Message)
	return true
}

func (h *ViberHandler) Channel() notification.NotificationChannel {
	return notification.NotificationChannelViber
}
