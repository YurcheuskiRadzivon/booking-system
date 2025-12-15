package server

import (
	"encoding/json"
	"time"

	"github.com/YurcheuskiRadzivon/booking-system/internal/service/booking"
	"github.com/YurcheuskiRadzivon/booking-system/internal/service/notification"
	"github.com/gofiber/fiber/v2"
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
}

type Error struct {
	Message string `json:"message" example:"message"`
}

func New(port string, booking booking.Service, notification notification.Service) *Server {
	s := &Server{
		app:     nil,
		notify:  make(chan error, 1),
		address: port,
	}

	app := fiber.New(fiber.Config{
		BodyLimit:    30 * 1024 * 1024,
		ReadTimeout:  _defaultReadTimeout,
		WriteTimeout: _defaultWriteTimeout,
		JSONEncoder:  json.Marshal,
		JSONDecoder:  json.Unmarshal,
	})

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
		ui.Get("/")
	}

	booking := s.app.Group("/booking")
	{
		booking.Get("/")
	}

	notification := s.app.Group("/notification")
	{
		notification.Get("/")
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
