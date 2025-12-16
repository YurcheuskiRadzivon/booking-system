package server

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/YurcheuskiRadzivon/booking-system/internal/service/admin"
	"github.com/YurcheuskiRadzivon/booking-system/internal/service/booking"
	"github.com/YurcheuskiRadzivon/booking-system/internal/service/notification"
	"github.com/YurcheuskiRadzivon/booking-system/web"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
)

const (
	_defaultAddr            = ":80"
	_defaultReadTimeout     = 40 * time.Second
	_defaultWriteTimeout    = 40 * time.Second
	_defaultShutdownTimeout = 40 * time.Second
)

type Server struct {
	app          *fiber.App
	notify       chan error
	address      string
	booking      booking.Service
	notification notification.Service
	admin        admin.Service
}

type Error struct {
	Message string `json:"message" example:"message"`
}

func New(port string, bookingSvc booking.Service, notificationSvc notification.Service, adminSvc admin.Service) *Server {
	s := &Server{
		app:          nil,
		notify:       make(chan error, 1),
		address:      port,
		booking:      bookingSvc,
		notification: notificationSvc,
		admin:        adminSvc,
	}

	app := fiber.New(fiber.Config{
		BodyLimit:    30 * 1024 * 1024,
		ReadTimeout:  _defaultReadTimeout,
		WriteTimeout: _defaultWriteTimeout,
		JSONEncoder:  json.Marshal,
		JSONDecoder:  json.Unmarshal,
	})

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	s.app = app

	return s
}

func (s *Server) Start() {
	go func() {
		s.notify <- s.app.Listen(s.address)
		close(s.notify)
	}()
}

func (s *Server) RegisterRoutes() {
	ui := s.app.Group("/ui")
	{
		ui.Use("/", filesystem.New(filesystem.Config{
			Root:       http.FS(web.Assets),
			PathPrefix: "",
			Browse:     true,
		}))
	}

	bookingGroup := s.app.Group("/booking")
	{
		bookingGroup.Get("/rooms", s.handleGetRooms)
		bookingGroup.Get("/rooms/search", s.handleSearchRooms)
		bookingGroup.Get("/rooms/:id", s.handleGetRoomByID)
		bookingGroup.Post("/", s.handleCreateBooking)
		bookingGroup.Get("/:id", s.handleGetBooking)
		bookingGroup.Put("/:id/confirm", s.handleConfirmBooking)
		bookingGroup.Put("/:id/cancel", s.handleCancelBooking)
		bookingGroup.Post("/price", s.handleCalculatePrice)
		bookingGroup.Get("/my", s.handleGetMyBookings)
	}

	notificationGroup := s.app.Group("/notification")
	{
		notificationGroup.Post("/send", s.handleSendNotification)
		notificationGroup.Post("/broadcast", s.handleBroadcastNotification)
		notificationGroup.Get("/types", s.handleGetNotificationTypes)
	}

	adminGroup := s.app.Group("/admin")
	{
		adminGroup.Get("/rooms", s.handleAdminGetRooms)
		adminGroup.Post("/rooms", s.handleAdminCreateRoom)
		adminGroup.Put("/rooms/:id", s.handleAdminUpdateRoom)
		adminGroup.Delete("/rooms/:id", s.handleAdminDeleteRoom)
		adminGroup.Get("/bookings", s.handleAdminGetBookings)
		adminGroup.Put("/bookings/:id/status", s.handleAdminUpdateBookingStatus)
		adminGroup.Get("/stats", s.handleAdminGetStats)
		adminGroup.Get("/status", s.handleAdminGetStatus)

		adminGroup.Get("/dates", s.handleAdminGetSpecialDates)
		adminGroup.Post("/dates", s.handleAdminCreateSpecialDate)
		adminGroup.Delete("/dates/:id", s.handleAdminDeleteSpecialDate)
	}
}

func (s *Server) Notify() <-chan error {
	return s.notify
}

func (s *Server) Shutdown() error {
	return s.app.ShutdownWithTimeout(_defaultShutdownTimeout)
}

func ErrorResponse(ctx *fiber.Ctx, code int, msg string) error {
	return ctx.Status(code).JSON(Error{Message: msg})
}
