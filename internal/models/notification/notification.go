package notification

import (
	"time"
)

type NotificationChannel string

const (
	NotificationChannelEmail NotificationChannel = "email"
	NotificationChannelSMS   NotificationChannel = "sms"
	NotificationChannelViber NotificationChannel = "viber"
)

type EventType string

const (
	EventTypeBookingCreated   EventType = "booking_created"
	EventTypeBookingConfirmed EventType = "booking_confirmed"
	EventTypeBookingCancelled EventType = "booking_cancelled"
)

type NotificationType struct {
	ID      int64  `json:"id" db:"id"`
	Name    string `json:"name" db:"name"`
	Message string `json:"message" db:"message"`
}

type NotificationEvent struct {
	ID        string              `json:"id"`
	Type      EventType           `json:"type"`
	Channel   NotificationChannel `json:"channel"`
	Recipient string              `json:"recipient"`
	Subject   string              `json:"subject"`
	Message   string              `json:"message"`
	Data      map[string]any      `json:"data,omitempty"`
	CreatedAt time.Time           `json:"created_at"`
}

type SendNotificationRequest struct {
	Channel   NotificationChannel `json:"channel"`
	Recipient string              `json:"recipient"`
	Subject   string              `json:"subject,omitempty"`
	Message   string              `json:"message"`
}

type NotificationResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	EventID string `json:"event_id,omitempty"`
}
