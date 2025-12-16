package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/YurcheuskiRadzivon/booking-system/internal/config"
	"github.com/YurcheuskiRadzivon/booking-system/internal/repository"
	"github.com/YurcheuskiRadzivon/booking-system/internal/server"
	"github.com/YurcheuskiRadzivon/booking-system/internal/service/admin"
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
	repo, err := repository.NewPostgresRepository(cfg.Database.ConnectionString)
	if err != nil {
		log.Fatalf("Repository error: %v", err)
	}
	defer repo.Close()

	log.Println("Connected to database")

	bookingSvc, err := booking.NewService(ctx, repo)
	if err != nil {
		log.Fatalf("Booking service error: %v", err)
	}
	log.Println("Booking service initialized")

	notificationSvc, err := notification.NewService(ctx, repo)
	if err != nil {
		log.Fatalf("Notification service error: %v", err)
	}
	log.Println("Notification service initialized")

	notificationSvc.StartWorker(ctx)

	adminSvc, err := admin.NewService(ctx, repo)
	if err != nil {
		log.Fatalf("Admin service error: %v", err)
	}
	log.Println("Admin service initialized")

	srv := server.New(cfg.HTTP.PORT, bookingSvc, notificationSvc, adminSvc)
	srv.RegisterRoutes()
	srv.Start()

	log.Printf("Server started on %s", cfg.HTTP.PORT)
	log.Printf("Open http://localhost%s/ui/index.html in browser", cfg.HTTP.PORT)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case <-interrupt:
		log.Println("Shutting down...")

	case err := <-srv.Notify():
		log.Panicf("Server error: %s", err)
	}

	err = srv.Shutdown()
	if err != nil {
		log.Fatalf("Server shutdown error: %v", err)
	}
	log.Println("Server stopped")
}
