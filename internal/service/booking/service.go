package booking

import (
	"context"
	"errors"

	"github.com/YurcheuskiRadzivon/booking-system/internal/models/booking"
	"github.com/YurcheuskiRadzivon/booking-system/internal/repository"
)

var (
	ErrRoomNotFound     = errors.New("room not found")
	ErrRoomNotAvailable = errors.New("room is not available for selected dates")
	ErrBookingNotFound  = errors.New("booking not found")
	ErrInvalidDates     = errors.New("invalid booking dates")
	ErrInvalidGuestInfo = errors.New("invalid guest information")
)

type Service interface {
	GetAllRooms(ctx context.Context) ([]booking.Room, error)
	GetRoomByID(ctx context.Context, id int64) (*booking.Room, error)
	FindAvailableRooms(ctx context.Context, req booking.RoomSearchRequest) ([]booking.RoomWithAvailability, error)

	CreateBooking(ctx context.Context, req booking.CreateBookingRequest) (*booking.BookingResponse, error)
	GetBookingByID(ctx context.Context, id int64) (*booking.BookingWithRoom, error)
	GetAllBookings(ctx context.Context) ([]booking.Booking, error)
	GetBookingsByEmail(ctx context.Context, email string) ([]booking.Booking, error)
	ConfirmBooking(ctx context.Context, id int64) (*booking.Booking, error)
	CancelBooking(ctx context.Context, id int64) (*booking.Booking, error)

	CalculatePrice(ctx context.Context, req booking.PriceCalculationRequest) (*booking.PriceCalculationResponse, error)

	CreateRoom(ctx context.Context, room *booking.Room) error
	UpdateRoom(ctx context.Context, room *booking.Room) error
	DeleteRoom(ctx context.Context, id int64) error

	GetSpecialDates(ctx context.Context) ([]booking.SpecialDate, error)
	CreateSpecialDate(ctx context.Context, sd *booking.SpecialDate) error
	DeleteSpecialDate(ctx context.Context, id int64) error
}

type service struct {
	ctx         context.Context
	repo        repository.Repository
	roomFactory *RoomFactory
}

func NewService(ctx context.Context, repo repository.Repository) (Service, error) {
	srv := &service{
		ctx:         ctx,
		repo:        repo,
		roomFactory: NewRoomFactory(),
	}

	return srv, nil
}

func (s *service) GetAllRooms(ctx context.Context) ([]booking.Room, error) {
	return s.repo.Room().GetAll(ctx)
}

func (s *service) GetRoomByID(ctx context.Context, id int64) (*booking.Room, error) {
	room, err := s.repo.Room().GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if room == nil {
		return nil, ErrRoomNotFound
	}
	return room, nil
}

func (s *service) FindAvailableRooms(ctx context.Context, req booking.RoomSearchRequest) ([]booking.RoomWithAvailability, error) {
	if req.CheckIn.IsZero() || req.CheckOut.IsZero() || req.CheckOut.Before(req.CheckIn) {
		return nil, ErrInvalidDates
	}

	var rooms []booking.Room
	var err error

	if req.RoomType != "" {
		rooms, err = s.repo.Room().GetAvailableByType(ctx, req.RoomType, req.CheckIn, req.CheckOut)
	} else if req.Capacity > 0 {
		rooms, err = s.repo.Room().GetAvailableByCapacity(ctx, req.Capacity, req.CheckIn, req.CheckOut)
	} else {
		rooms, err = s.repo.Room().GetAvailable(ctx, req.CheckIn, req.CheckOut)
	}

	if err != nil {
		return nil, err
	}

	specialDates, _ := s.repo.SpecialDate().GetByDateRange(ctx, req.CheckIn, req.CheckOut)
	calculator := NewPriceCalculator(specialDates)

	result := make([]booking.RoomWithAvailability, len(rooms))
	for i, room := range rooms {
		priceInfo := calculator.CalculateTotalPrice(room.BasePrice, req.CheckIn, req.CheckOut)
		result[i] = booking.RoomWithAvailability{
			Room:        room,
			IsAvailable: true,
			TotalPrice:  priceInfo.TotalPrice,
		}
	}

	return result, nil
}

func (s *service) CreateBooking(ctx context.Context, req booking.CreateBookingRequest) (*booking.BookingResponse, error) {
	if req.StartDate.IsZero() || req.EndDate.IsZero() || req.EndDate.Before(req.StartDate) {
		return nil, ErrInvalidDates
	}

	if req.GuestInfo.Name == "" || req.GuestInfo.Email == "" {
		return nil, ErrInvalidGuestInfo
	}

	room, err := s.repo.Room().GetByID(ctx, req.RoomID)
	if err != nil {
		return nil, err
	}
	if room == nil {
		return nil, ErrRoomNotFound
	}

	available, err := s.repo.Booking().IsRoomAvailable(ctx, req.RoomID, req.StartDate, req.EndDate)
	if err != nil {
		return nil, err
	}
	if !available {
		return nil, ErrRoomNotAvailable
	}

	specialDates, _ := s.repo.SpecialDate().GetByDateRange(ctx, req.StartDate, req.EndDate)
	calculator := NewPriceCalculator(specialDates)
	priceInfo := calculator.CalculateTotalPrice(room.BasePrice, req.StartDate, req.EndDate)

	newBooking := &booking.Booking{
		StartDate: req.StartDate,
		EndDate:   req.EndDate,
		RoomID:    req.RoomID,
		GuestInfo: req.GuestInfo,
		Price:     priceInfo.TotalPrice,
		Status:    booking.BookingStatusPending,
	}

	if err := s.repo.Booking().Create(ctx, newBooking); err != nil {
		return nil, err
	}

	return &booking.BookingResponse{
		Booking: *newBooking,
		Room:    *room,
		Nights:  priceInfo.Nights,
	}, nil
}

func (s *service) GetBookingByID(ctx context.Context, id int64) (*booking.BookingWithRoom, error) {
	b, err := s.repo.Booking().GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if b == nil {
		return nil, ErrBookingNotFound
	}

	room, err := s.repo.Room().GetByID(ctx, b.RoomID)
	if err != nil {
		return nil, err
	}

	return &booking.BookingWithRoom{
		Booking: *b,
		Room:    *room,
	}, nil
}

func (s *service) GetAllBookings(ctx context.Context) ([]booking.Booking, error) {
	return s.repo.Booking().GetAll(ctx)
}

func (s *service) GetBookingsByEmail(ctx context.Context, email string) ([]booking.Booking, error) {
	return s.repo.Booking().GetByEmail(ctx, email)
}

func (s *service) ConfirmBooking(ctx context.Context, id int64) (*booking.Booking, error) {
	b, err := s.repo.Booking().GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if b == nil {
		return nil, ErrBookingNotFound
	}

	if err := s.repo.Booking().UpdateStatus(ctx, id, booking.BookingStatusConfirmed); err != nil {
		return nil, err
	}

	b.Status = booking.BookingStatusConfirmed
	return b, nil
}

func (s *service) CancelBooking(ctx context.Context, id int64) (*booking.Booking, error) {
	b, err := s.repo.Booking().GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if b == nil {
		return nil, ErrBookingNotFound
	}

	if err := s.repo.Booking().UpdateStatus(ctx, id, booking.BookingStatusCancelled); err != nil {
		return nil, err
	}

	b.Status = booking.BookingStatusCancelled
	return b, nil
}

func (s *service) CalculatePrice(ctx context.Context, req booking.PriceCalculationRequest) (*booking.PriceCalculationResponse, error) {
	if req.CheckIn.IsZero() || req.CheckOut.IsZero() || req.CheckOut.Before(req.CheckIn) {
		return nil, ErrInvalidDates
	}

	room, err := s.repo.Room().GetByID(ctx, req.RoomID)
	if err != nil {
		return nil, err
	}
	if room == nil {
		return nil, ErrRoomNotFound
	}

	specialDates, _ := s.repo.SpecialDate().GetByDateRange(ctx, req.CheckIn, req.CheckOut)
	calculator := NewPriceCalculator(specialDates)
	priceInfo := calculator.CalculateTotalPrice(room.BasePrice, req.CheckIn, req.CheckOut)

	return &priceInfo, nil
}

func (s *service) CreateRoom(ctx context.Context, room *booking.Room) error {
	return s.repo.Room().Create(ctx, room)
}

func (s *service) UpdateRoom(ctx context.Context, room *booking.Room) error {
	return s.repo.Room().Update(ctx, room)
}

func (s *service) DeleteRoom(ctx context.Context, id int64) error {
	return s.repo.Room().Delete(ctx, id)
}

func (s *service) GetSpecialDates(ctx context.Context) ([]booking.SpecialDate, error) {
	return s.repo.SpecialDate().GetAll(ctx)
}

func (s *service) CreateSpecialDate(ctx context.Context, sd *booking.SpecialDate) error {
	return s.repo.SpecialDate().Create(ctx, sd)
}

func (s *service) DeleteSpecialDate(ctx context.Context, id int64) error {
	return s.repo.SpecialDate().Delete(ctx, id)
}
