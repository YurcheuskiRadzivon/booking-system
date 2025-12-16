package server

import (
	"net/http"

	"github.com/YurcheuskiRadzivon/booking-system/internal/models/notification"
	"github.com/gofiber/fiber/v2"
)

type sendNotificationRequest struct {
	Channel   string `json:"channel"`
	Recipient string `json:"recipient"`
	Subject   string `json:"subject,omitempty"`
	Message   string `json:"message"`
}

type broadcastRequest struct {
	Channels  []string `json:"channels"`
	Recipient string   `json:"recipient"`
	Subject   string   `json:"subject,omitempty"`
	Message   string   `json:"message"`
}

func (s *Server) handleSendNotification(ctx *fiber.Ctx) error {
	var req sendNotificationRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ErrorResponse(ctx, http.StatusBadRequest, "Invalid request body")
	}

	if req.Recipient == "" || req.Message == "" {
		return ErrorResponse(ctx, http.StatusBadRequest, "Recipient and message are required")
	}

	var resp *notification.NotificationResponse
	var err error

	switch notification.NotificationChannel(req.Channel) {
	case notification.NotificationChannelEmail:
		resp, err = s.notification.SendEmail(ctx.Context(), req.Recipient, req.Subject, req.Message)
	case notification.NotificationChannelSMS:
		resp, err = s.notification.SendSMS(ctx.Context(), req.Recipient, req.Message)
	case notification.NotificationChannelViber:
		resp, err = s.notification.SendViber(ctx.Context(), req.Recipient, req.Message)
	default:
		return ErrorResponse(ctx, http.StatusBadRequest, "Unknown notification channel")
	}

	if err != nil {
		return ErrorResponse(ctx, http.StatusInternalServerError, err.Error())
	}

	return ctx.Status(http.StatusOK).JSON(resp)
}

func (s *Server) handleBroadcastNotification(ctx *fiber.Ctx) error {
	var req broadcastRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ErrorResponse(ctx, http.StatusBadRequest, "Invalid request body")
	}

	if req.Recipient == "" || req.Message == "" || len(req.Channels) == 0 {
		return ErrorResponse(ctx, http.StatusBadRequest, "Recipient, message, and channels are required")
	}

	channels := make([]notification.NotificationChannel, len(req.Channels))
	for i, ch := range req.Channels {
		channels[i] = notification.NotificationChannel(ch)
	}

	responses, err := s.notification.Broadcast(ctx.Context(), channels, req.Recipient, req.Subject, req.Message)
	if err != nil {
		return ErrorResponse(ctx, http.StatusInternalServerError, err.Error())
	}

	return ctx.Status(http.StatusOK).JSON(responses)
}

func (s *Server) handleGetNotificationTypes(ctx *fiber.Ctx) error {
	types, err := s.notification.GetNotificationTypes(ctx.Context())
	if err != nil {
		return ErrorResponse(ctx, http.StatusInternalServerError, err.Error())
	}

	return ctx.Status(http.StatusOK).JSON(types)
}
