package server

import (
	"net/http"
	"strconv"
	"time"

	bookingModel "github.com/YurcheuskiRadzivon/booking-system/internal/models/booking"
	"github.com/gofiber/fiber/v2"
)

func (s *Server) handleGetRooms(ctx *fiber.Ctx) error {
	rooms, err := s.booking.GetAllRooms(ctx.Context())
	if err != nil {
		return ErrorResponse(ctx, http.StatusInternalServerError, err.Error())
	}
	return ctx.Status(http.StatusOK).JSON(rooms)
}

func (s *Server) handleSearchRooms(ctx *fiber.Ctx) error {
	checkInStr := ctx.Query("check_in")
	checkOutStr := ctx.Query("check_out")
	roomType := ctx.Query("room_type")
	capacityStr := ctx.Query("capacity")

	var checkIn, checkOut time.Time
	var err error

	if checkInStr != "" {
		checkIn, err = time.Parse("2006-01-02", checkInStr)
		if err != nil {
			return ErrorResponse(ctx, http.StatusBadRequest, "Invalid check_in date format (use YYYY-MM-DD)")
		}
	} else {
		checkIn = time.Now()
	}

	if checkOutStr != "" {
		checkOut, err = time.Parse("2006-01-02", checkOutStr)
		if err != nil {
			return ErrorResponse(ctx, http.StatusBadRequest, "Invalid check_out date format (use YYYY-MM-DD)")
		}
	} else {
		checkOut = checkIn.AddDate(0, 0, 1)
	}

	var capacity int
	if capacityStr != "" {
		capacity, _ = strconv.Atoi(capacityStr)
	}

	req := bookingModel.RoomSearchRequest{
		CheckIn:  checkIn,
		CheckOut: checkOut,
		RoomType: bookingModel.RoomType(roomType),
		Capacity: capacity,
	}

	rooms, err := s.booking.FindAvailableRooms(ctx.Context(), req)
	if err != nil {
		return ErrorResponse(ctx, http.StatusInternalServerError, err.Error())
	}

	return ctx.Status(http.StatusOK).JSON(rooms)
}

func (s *Server) handleGetRoomByID(ctx *fiber.Ctx) error {
	id, err := strconv.ParseInt(ctx.Params("id"), 10, 64)
	if err != nil {
		return ErrorResponse(ctx, http.StatusBadRequest, "Invalid room ID")
	}

	room, err := s.booking.GetRoomByID(ctx.Context(), id)
	if err != nil {
		return ErrorResponse(ctx, http.StatusNotFound, err.Error())
	}

	return ctx.Status(http.StatusOK).JSON(room)
}

func (s *Server) handleCreateBooking(ctx *fiber.Ctx) error {
	var req bookingModel.CreateBookingRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ErrorResponse(ctx, http.StatusBadRequest, "Invalid request body")
	}

	booking, err := s.booking.CreateBooking(ctx.Context(), req)
	if err != nil {
		return ErrorResponse(ctx, http.StatusBadRequest, err.Error())
	}

	s.notification.NotifyBookingCreated(ctx.Context(), &booking.Booking, &booking.Room)

	return ctx.Status(http.StatusCreated).JSON(booking)
}

func (s *Server) handleGetBooking(ctx *fiber.Ctx) error {
	id, err := strconv.ParseInt(ctx.Params("id"), 10, 64)
	if err != nil {
		return ErrorResponse(ctx, http.StatusBadRequest, "Invalid booking ID")
	}

	booking, err := s.booking.GetBookingByID(ctx.Context(), id)
	if err != nil {
		return ErrorResponse(ctx, http.StatusNotFound, err.Error())
	}

	return ctx.Status(http.StatusOK).JSON(booking)
}

func (s *Server) handleGetMyBookings(ctx *fiber.Ctx) error {
	email := ctx.Query("email")
	if email == "" {
		return ErrorResponse(ctx, http.StatusBadRequest, "Email is required")
	}

	bookings, err := s.booking.GetBookingsByEmail(ctx.Context(), email)
	if err != nil {
		return ErrorResponse(ctx, http.StatusInternalServerError, err.Error())
	}

	return ctx.Status(http.StatusOK).JSON(bookings)
}

func (s *Server) handleConfirmBooking(ctx *fiber.Ctx) error {
	id, err := strconv.ParseInt(ctx.Params("id"), 10, 64)
	if err != nil {
		return ErrorResponse(ctx, http.StatusBadRequest, "Invalid booking ID")
	}

	booking, err := s.booking.ConfirmBooking(ctx.Context(), id)
	if err != nil {
		return ErrorResponse(ctx, http.StatusBadRequest, err.Error())
	}

	bookingWithRoom, _ := s.booking.GetBookingByID(ctx.Context(), id)
	if bookingWithRoom != nil {
		s.notification.NotifyBookingConfirmed(ctx.Context(), booking, &bookingWithRoom.Room)
	}

	return ctx.Status(http.StatusOK).JSON(booking)
}

func (s *Server) handleCancelBooking(ctx *fiber.Ctx) error {
	id, err := strconv.ParseInt(ctx.Params("id"), 10, 64)
	if err != nil {
		return ErrorResponse(ctx, http.StatusBadRequest, "Invalid booking ID")
	}

	bookingWithRoom, _ := s.booking.GetBookingByID(ctx.Context(), id)

	booking, err := s.booking.CancelBooking(ctx.Context(), id)
	if err != nil {
		return ErrorResponse(ctx, http.StatusBadRequest, err.Error())
	}

	if bookingWithRoom != nil {
		s.notification.NotifyBookingCancelled(ctx.Context(), booking, &bookingWithRoom.Room)
	}

	return ctx.Status(http.StatusOK).JSON(booking)
}

func (s *Server) handleCalculatePrice(ctx *fiber.Ctx) error {
	var req bookingModel.PriceCalculationRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ErrorResponse(ctx, http.StatusBadRequest, "Invalid request body")
	}

	price, err := s.booking.CalculatePrice(ctx.Context(), req)
	if err != nil {
		return ErrorResponse(ctx, http.StatusBadRequest, err.Error())
	}

	return ctx.Status(http.StatusOK).JSON(price)
}
