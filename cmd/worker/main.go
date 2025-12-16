package main

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"
	"time"

	notificationModel "github.com/YurcheuskiRadzivon/booking-system/internal/models/notification"
	notificationSvc "github.com/YurcheuskiRadzivon/booking-system/internal/service/notification"
)

func main() {
	fmt.Println("🔔 Notification Worker Starting...")
	fmt.Println("======================================")

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	broker := notificationSvc.NewMessageBroker()

	emailEvents := broker.Subscribe(notificationModel.NotificationChannelEmail, 100)
	smsEvents := broker.Subscribe(notificationModel.NotificationChannelSMS, 100)
	viberEvents := broker.Subscribe(notificationModel.NotificationChannelViber, 100)

	emailHandler := notificationSvc.NewEmailHandler()
	smsHandler := notificationSvc.NewSMSHandler()
	viberHandler := notificationSvc.NewViberHandler()

	go func() {
		for event := range emailEvents {
			fmt.Printf("\n[%s] 📧 Email Event Received\n", time.Now().Format("15:04:05"))
			emailHandler.Send(event)
		}
	}()

	go func() {
		for event := range smsEvents {
			fmt.Printf("\n[%s] 📱 SMS Event Received\n", time.Now().Format("15:04:05"))
			smsHandler.Send(event)
		}
	}()

	go func() {
		for event := range viberEvents {
			fmt.Printf("\n[%s] 💬 Viber Event Received\n", time.Now().Format("15:04:05"))
			viberHandler.Send(event)
		}
	}()

	fmt.Println("🎧 Listening for notification events...")
	fmt.Println("Press Ctrl+C to stop")

	go func() {
		time.Sleep(2 * time.Second)
		fmt.Println("\n📬 Sending test notification...")
		broker.Publish(notificationModel.NotificationEvent{
			Type:      notificationModel.EventTypeBookingCreated,
			Channel:   notificationModel.NotificationChannelEmail,
			Recipient: "test@example.com",
			Subject:   "Тестовое уведомление",
			Message:   "Это тестовое сообщение от системы бронирования",
			CreatedAt: time.Now(),
		})
	}()

	<-ctx.Done()
	fmt.Println("\n🛑 Worker stopping...")
	broker.Close()
	fmt.Println("👋 Worker stopped")
}
