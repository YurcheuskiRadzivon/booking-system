package server

import (
	"net/http"
	"strconv"
	"time"

	bookingModel "github.com/YurcheuskiRadzivon/booking-system/internal/models/booking"
	"github.com/gofiber/fiber/v2"
)

func (s *Server) handleAdminGetRooms(ctx *fiber.Ctx) error {
	rooms, err := s.admin.GetAllRooms(ctx.Context())
	if err != nil {
		return ErrorResponse(ctx, http.StatusInternalServerError, err.Error())
	}
	if rooms == nil {
		rooms = []bookingModel.Room{}
	}
	return ctx.Status(http.StatusOK).JSON(rooms)
}

func (s *Server) handleAdminCreateRoom(ctx *fiber.Ctx) error {
	var room bookingModel.Room
	if err := ctx.BodyParser(&room); err != nil {
		return ErrorResponse(ctx, http.StatusBadRequest, "Invalid request body")
	}

	if err := s.admin.CreateRoom(ctx.Context(), &room); err != nil {
		return ErrorResponse(ctx, http.StatusBadRequest, err.Error())
	}

	return ctx.Status(http.StatusCreated).JSON(room)
}

func (s *Server) handleAdminUpdateRoom(ctx *fiber.Ctx) error {
	id, err := strconv.ParseInt(ctx.Params("id"), 10, 64)
	if err != nil {
		return ErrorResponse(ctx, http.StatusBadRequest, "Invalid room ID")
	}

	var room bookingModel.Room
	if err := ctx.BodyParser(&room); err != nil {
		return ErrorResponse(ctx, http.StatusBadRequest, "Invalid request body")
	}
	room.ID = id

	if err := s.admin.UpdateRoom(ctx.Context(), &room); err != nil {
		return ErrorResponse(ctx, http.StatusBadRequest, err.Error())
	}

	return ctx.Status(http.StatusOK).JSON(room)
}

func (s *Server) handleAdminDeleteRoom(ctx *fiber.Ctx) error {
	id, err := strconv.ParseInt(ctx.Params("id"), 10, 64)
	if err != nil {
		return ErrorResponse(ctx, http.StatusBadRequest, "Invalid room ID")
	}

	if err := s.admin.DeleteRoom(ctx.Context(), id); err != nil {
		return ErrorResponse(ctx, http.StatusBadRequest, err.Error())
	}

	return ctx.Status(http.StatusOK).JSON(fiber.Map{"message": "Room deleted successfully"})
}

func (s *Server) handleAdminGetBookings(ctx *fiber.Ctx) error {
	status := ctx.Query("status")

	bookings := make([]bookingModel.BookingWithRoom, 0)
	var err error

	if status != "" {
		bookings, err = s.admin.GetBookingsByStatus(ctx.Context(), bookingModel.BookingStatus(status))
	} else {
		bookings, err = s.admin.GetAllBookings(ctx.Context())
	}

	if err != nil {
		return ErrorResponse(ctx, http.StatusInternalServerError, err.Error())
	}

	if bookings == nil {
		bookings = []bookingModel.BookingWithRoom{}
	}

	return ctx.Status(http.StatusOK).JSON(bookings)
}

type updateStatusRequest struct {
	Status string `json:"status"`
}

func (s *Server) handleAdminUpdateBookingStatus(ctx *fiber.Ctx) error {
	id, err := strconv.ParseInt(ctx.Params("id"), 10, 64)
	if err != nil {
		return ErrorResponse(ctx, http.StatusBadRequest, "Invalid booking ID")
	}

	var req updateStatusRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ErrorResponse(ctx, http.StatusBadRequest, "Invalid request body")
	}

	if err := s.admin.UpdateBookingStatus(ctx.Context(), id, bookingModel.BookingStatus(req.Status)); err != nil {
		return ErrorResponse(ctx, http.StatusBadRequest, err.Error())
	}

	return ctx.Status(http.StatusOK).JSON(fiber.Map{"message": "Booking status updated successfully"})
}

func (s *Server) handleAdminGetStats(ctx *fiber.Ctx) error {
	stats, err := s.admin.GetStatistics(ctx.Context())
	if err != nil {
		return ErrorResponse(ctx, http.StatusInternalServerError, err.Error())
	}

	return ctx.Status(http.StatusOK).JSON(stats)
}

func (s *Server) handleAdminGetStatus(ctx *fiber.Ctx) error {
	status, err := s.admin.GetHotelStatus(ctx.Context())
	if err != nil {
		return ErrorResponse(ctx, http.StatusInternalServerError, err.Error())
	}

	return ctx.Status(http.StatusOK).JSON(fiber.Map{"status": status})
}

func (s *Server) handleAdminGetSpecialDates(ctx *fiber.Ctx) error {
	dates, err := s.booking.GetSpecialDates(ctx.Context())
	if err != nil {
		return ErrorResponse(ctx, http.StatusInternalServerError, err.Error())
	}
	return ctx.Status(http.StatusOK).JSON(dates)
}

func (s *Server) handleAdminCreateSpecialDate(ctx *fiber.Ctx) error {
	var req bookingModel.CreateSpecialDateRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ErrorResponse(ctx, http.StatusBadRequest, "Invalid request body")
	}

	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return ErrorResponse(ctx, http.StatusBadRequest, "Invalid date format")
	}

	sd := &bookingModel.SpecialDate{
		Date:        date,
		Name:        req.Name,
		Coefficient: req.Coefficient,
	}

	if err := s.booking.CreateSpecialDate(ctx.Context(), sd); err != nil {
		return ErrorResponse(ctx, http.StatusBadRequest, err.Error())
	}

	return ctx.Status(http.StatusCreated).JSON(sd)
}

func (s *Server) handleAdminDeleteSpecialDate(ctx *fiber.Ctx) error {
	id, err := strconv.ParseInt(ctx.Params("id"), 10, 64)
	if err != nil {
		return ErrorResponse(ctx, http.StatusBadRequest, "Invalid ID")
	}

	if err := s.booking.DeleteSpecialDate(ctx.Context(), id); err != nil {
		return ErrorResponse(ctx, http.StatusBadRequest, err.Error())
	}

	return ctx.Status(http.StatusOK).JSON(fiber.Map{"message": "Special date deleted"})
}
