package notification

import (
	"context"
	"fmt"
	"strings"
	"time"

	bookingModel "github.com/YurcheuskiRadzivon/booking-system/internal/models/booking"
	"github.com/YurcheuskiRadzivon/booking-system/internal/models/notification"
	"github.com/YurcheuskiRadzivon/booking-system/internal/repository"
	"github.com/google/uuid"
)

type NotificationHandler interface {
	Send(event notification.NotificationEvent) bool
	Channel() notification.NotificationChannel
}

type Service interface {
	SendEmail(ctx context.Context, recipient, subject, message string) (*notification.NotificationResponse, error)
	SendSMS(ctx context.Context, recipient, message string) (*notification.NotificationResponse, error)
	SendViber(ctx context.Context, recipient, message string) (*notification.NotificationResponse, error)
	Broadcast(ctx context.Context, channels []notification.NotificationChannel, recipient, subject, message string) ([]notification.NotificationResponse, error)

	NotifyBookingCreated(ctx context.Context, booking *bookingModel.Booking, room *bookingModel.Room) error
	NotifyBookingConfirmed(ctx context.Context, booking *bookingModel.Booking, room *bookingModel.Room) error
	NotifyBookingCancelled(ctx context.Context, booking *bookingModel.Booking, room *bookingModel.Room) error

	GetNotificationTypes(ctx context.Context) ([]notification.NotificationType, error)

	GetBroker() *MessageBroker
	StartWorker(ctx context.Context)
}

type service struct {
	ctx      context.Context
	repo     repository.Repository
	broker   *MessageBroker
	handlers map[notification.NotificationChannel]NotificationHandler
}

func NewService(ctx context.Context, repo repository.Repository) (Service, error) {
	broker := NewMessageBroker()

	srv := &service{
		ctx:    ctx,
		repo:   repo,
		broker: broker,
		handlers: map[notification.NotificationChannel]NotificationHandler{
			notification.NotificationChannelEmail: NewEmailHandler(),
			notification.NotificationChannelSMS:   NewSMSHandler(),
			notification.NotificationChannelViber: NewViberHandler(),
		},
	}

	return srv, nil
}

func (s *service) GetBroker() *MessageBroker {
	return s.broker
}

func (s *service) StartWorker(ctx context.Context) {
	allEvents := s.broker.Subscribe("all", 100)

	go func() {
		fmt.Println(" Notification worker started")
		for {
			select {
			case event, ok := <-allEvents:
				if !ok {
					fmt.Println(" Notification worker stopped")
					return
				}
				s.processEvent(event)
			case <-ctx.Done():
				fmt.Println(" Notification worker stopping...")
				return
			}
		}
	}()
}

func (s *service) processEvent(event notification.NotificationEvent) {
	handler, ok := s.handlers[event.Channel]
	if !ok {
		return
	}

	fmt.Printf("\n--- Processing notification [%s] ---\n", event.Type)
	handler.Send(event)
	fmt.Println("--- Notification sent ---\n")
}

func (s *service) SendEmail(ctx context.Context, recipient, subject, message string) (*notification.NotificationResponse, error) {
	event := notification.NotificationEvent{
		ID:        uuid.New().String(),
		Type:      "manual",
		Channel:   notification.NotificationChannelEmail,
		Recipient: recipient,
		Subject:   subject,
		Message:   message,
		CreatedAt: time.Now(),
	}

	s.broker.Publish(event)

	return &notification.NotificationResponse{
		Success: true,
		Message: "Email notification queued",
		EventID: event.ID,
	}, nil
}

func (s *service) SendSMS(ctx context.Context, recipient, message string) (*notification.NotificationResponse, error) {
	event := notification.NotificationEvent{
		ID:        uuid.New().String(),
		Type:      "manual",
		Channel:   notification.NotificationChannelSMS,
		Recipient: recipient,
		Message:   message,
		CreatedAt: time.Now(),
	}

	s.broker.Publish(event)

	return &notification.NotificationResponse{
		Success: true,
		Message: "SMS notification queued",
		EventID: event.ID,
	}, nil
}

func (s *service) SendViber(ctx context.Context, recipient, message string) (*notification.NotificationResponse, error) {
	event := notification.NotificationEvent{
		ID:        uuid.New().String(),
		Type:      "manual",
		Channel:   notification.NotificationChannelViber,
		Recipient: recipient,
		Message:   message,
		CreatedAt: time.Now(),
	}

	s.broker.Publish(event)

	return &notification.NotificationResponse{
		Success: true,
		Message: "Viber notification queued",
		EventID: event.ID,
	}, nil
}

func (s *service) Broadcast(ctx context.Context, channels []notification.NotificationChannel, recipient, subject, message string) ([]notification.NotificationResponse, error) {
	responses := make([]notification.NotificationResponse, len(channels))

	for i, channel := range channels {
		var resp *notification.NotificationResponse
		var err error

		switch channel {
		case notification.NotificationChannelEmail:
			resp, err = s.SendEmail(ctx, recipient, subject, message)
		case notification.NotificationChannelSMS:
			resp, err = s.SendSMS(ctx, recipient, message)
		case notification.NotificationChannelViber:
			resp, err = s.SendViber(ctx, recipient, message)
		default:
			responses[i] = notification.NotificationResponse{
				Success: false,
				Message: "Unknown channel: " + string(channel),
			}
			continue
		}

		if err != nil {
			responses[i] = notification.NotificationResponse{
				Success: false,
				Message: err.Error(),
			}
		} else {
			responses[i] = *resp
		}
	}

	return responses, nil
}

func (s *service) NotifyBookingCreated(ctx context.Context, booking *bookingModel.Booking, room *bookingModel.Room) error {
	message := s.formatBookingMessage(
		"Booking Created Successfully!",
		booking,
		room,
	)

	event := notification.NotificationEvent{
		ID:        uuid.New().String(),
		Type:      notification.EventTypeBookingCreated,
		Channel:   notification.NotificationChannelEmail,
		Recipient: booking.GuestInfo.Email,
		Subject:   "Booking Created - Room " + room.RoomNumber,
		Message:   message,
		Data: map[string]any{
			"booking_id": booking.ID,
			"room_id":    room.ID,
		},
		CreatedAt: time.Now(),
	}

	s.broker.Publish(event)
	return nil
}

func (s *service) NotifyBookingConfirmed(ctx context.Context, booking *bookingModel.Booking, room *bookingModel.Room) error {
	message := s.formatBookingMessage(
		"Booking Confirmed! We are waiting for you!",
		booking,
		room,
	)

	channels := []notification.NotificationChannel{
		notification.NotificationChannelEmail,
		notification.NotificationChannelSMS,
		notification.NotificationChannelViber,
	}

	s.Broadcast(ctx, channels, booking.GuestInfo.Email, "Booking Confirmed - Room "+room.RoomNumber, message)
	return nil
}

func (s *service) NotifyBookingCancelled(ctx context.Context, booking *bookingModel.Booking, room *bookingModel.Room) error {
	message := fmt.Sprintf(
		"Booking #%d cancelled.\nRoom: %s\nDates: %s - %s",
		booking.ID,
		room.RoomNumber,
		booking.StartDate.Format("02.01.2006"),
		booking.EndDate.Format("02.01.2006"),
	)

	channels := []notification.NotificationChannel{
		notification.NotificationChannelEmail,
		notification.NotificationChannelSMS,
		notification.NotificationChannelViber,
	}

	s.Broadcast(ctx, channels, booking.GuestInfo.Email, "Booking Cancelled - Room "+room.RoomNumber, message)
	return nil
}

func (s *service) formatBookingMessage(header string, booking *bookingModel.Booking, room *bookingModel.Room) string {
	var sb strings.Builder
	sb.WriteString(header)
	sb.WriteString("\n\n")
	sb.WriteString(fmt.Sprintf("Booking ID: #%d\n", booking.ID))
	sb.WriteString(fmt.Sprintf("Room: %s (%s)\n", room.RoomNumber, room.RoomType))
	sb.WriteString(fmt.Sprintf("Check-in: %s\n", booking.StartDate.Format("02.01.2006")))
	sb.WriteString(fmt.Sprintf("Check-out: %s\n", booking.EndDate.Format("02.01.2006")))
	sb.WriteString(fmt.Sprintf("Price: %.2f RUB\n", booking.Price))
	sb.WriteString(fmt.Sprintf("\nGuest: %s\n", booking.GuestInfo.Name))
	return sb.String()
}

func (s *service) GetNotificationTypes(ctx context.Context) ([]notification.NotificationType, error) {
	return s.repo.Notification().GetAll(ctx)
}
