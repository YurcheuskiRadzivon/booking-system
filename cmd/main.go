package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/YurcheuskiRadzivon/booking-system/internal/config"
	"github.com/YurcheuskiRadzivon/booking-system/internal/server"
	"github.com/YurcheuskiRadzivon/booking-system/internal/service/booking"
	"github.com/YurcheuskiRadzivon/booking-system/internal/service/notification"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	run(cfg, ctx)
}

func run(cfg *config.Config, ctx context.Context) {
	booking, err := booking.NewService(ctx)
	if err != nil {
		log.Fatalf("Booking: %v", err)
	}

	notification, err := notification.NewService(ctx)
	if err != nil {
		log.Fatalf("Notification: %v", err)
	}

	srv := server.New(cfg.HTTP.PORT, booking, notification)

	srv.RegisterRoutes()

	srv.Start()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case <-interrupt:
		log.Println("Shutdown")

	case err := <-srv.Notify():
		log.Panicf("server: %s", err)
	}

	err = srv.Shutdown()
	if err != nil {
		log.Fatalf("server: %v", err)
	}
}
