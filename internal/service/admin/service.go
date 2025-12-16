package admin

import (
	"context"
	"time"

	"github.com/YurcheuskiRadzivon/booking-system/internal/models/booking"
	"github.com/YurcheuskiRadzivon/booking-system/internal/repository"
)

type Statistics struct {
	TotalRooms        int                      `json:"total_rooms"`
	AvailableRooms    int                      `json:"available_rooms"`
	OccupiedRooms     int                      `json:"occupied_rooms"`
	TotalBookings     int                      `json:"total_bookings"`
	PendingBookings   int                      `json:"pending_bookings"`
	ConfirmedBookings int                      `json:"confirmed_bookings"`
	CancelledBookings int                      `json:"cancelled_bookings"`
	RoomsByType       map[booking.RoomType]int `json:"rooms_by_type"`
	TotalRevenue      float64                  `json:"total_revenue"`
}

type Service interface {
	GetAllRooms(ctx context.Context) ([]booking.Room, error)
	CreateRoom(ctx context.Context, room *booking.Room) error
	UpdateRoom(ctx context.Context, room *booking.Room) error
	DeleteRoom(ctx context.Context, id int64) error

	GetAllBookings(ctx context.Context) ([]booking.BookingWithRoom, error)
	GetBookingsByStatus(ctx context.Context, status booking.BookingStatus) ([]booking.BookingWithRoom, error)
	UpdateBookingStatus(ctx context.Context, id int64, status booking.BookingStatus) error

	GetStatistics(ctx context.Context) (*Statistics, error)
	GetHotelStatus(ctx context.Context) (string, error)
}

type service struct {
	ctx  context.Context
	repo repository.Repository
}

func NewService(ctx context.Context, repo repository.Repository) (Service, error) {
	return &service{
		ctx:  ctx,
		repo: repo,
	}, nil
}

func (s *service) GetAllRooms(ctx context.Context) ([]booking.Room, error) {
	return s.repo.Room().GetAll(ctx)
}

func (s *service) CreateRoom(ctx context.Context, room *booking.Room) error {
	if room.Status == "" {
		room.Status = booking.RoomStatusAvailable
	}
	return s.repo.Room().Create(ctx, room)
}

func (s *service) UpdateRoom(ctx context.Context, room *booking.Room) error {
	return s.repo.Room().Update(ctx, room)
}

func (s *service) DeleteRoom(ctx context.Context, id int64) error {
	return s.repo.Room().Delete(ctx, id)
}

func (s *service) GetAllBookings(ctx context.Context) ([]booking.BookingWithRoom, error) {
	return s.repo.Booking().GetAllWithRooms(ctx)
}

func (s *service) GetBookingsByStatus(ctx context.Context, status booking.BookingStatus) ([]booking.BookingWithRoom, error) {
	return s.repo.Booking().GetByStatusWithRooms(ctx, status)
}

func (s *service) UpdateBookingStatus(ctx context.Context, id int64, status booking.BookingStatus) error {
	return s.repo.Booking().UpdateStatus(ctx, id, status)
}

func (s *service) GetStatistics(ctx context.Context) (*Statistics, error) {
	rooms, err := s.repo.Room().GetAll(ctx)
	if err != nil {
		return nil, err
	}

	bookings, err := s.repo.Booking().GetAll(ctx)
	if err != nil {
		return nil, err
	}

	stats := &Statistics{
		TotalRooms:  len(rooms),
		RoomsByType: make(map[booking.RoomType]int),
	}

	now := time.Now()
	for _, room := range rooms {
		stats.RoomsByType[room.RoomType]++
		if room.Status == booking.RoomStatusAvailable {
			stats.AvailableRooms++
		} else {
			stats.OccupiedRooms++
		}
	}

	for _, b := range bookings {
		stats.TotalBookings++
		switch b.Status {
		case booking.BookingStatusPending:
			stats.PendingBookings++
		case booking.BookingStatusConfirmed:
			stats.ConfirmedBookings++
			if b.StartDate.Before(now) && b.EndDate.After(now) {
				stats.OccupiedRooms++
				stats.AvailableRooms--
			}
			stats.TotalRevenue += b.Price
		case booking.BookingStatusCancelled:
			stats.CancelledBookings++
		}
	}

	return stats, nil
}

func (s *service) GetHotelStatus(ctx context.Context) (string, error) {
	stats, err := s.GetStatistics(ctx)
	if err != nil {
		return "", err
	}

	status := "=== СТАТУС ОТЕЛЯ ===\n"
	status += "Всего номеров: " + itoa(stats.TotalRooms) + "\n"
	status += "Доступно: " + itoa(stats.AvailableRooms) + "\n"
	status += "Занято: " + itoa(stats.OccupiedRooms) + "\n"
	status += "\nБронирования:\n"
	status += "  Ожидают: " + itoa(stats.PendingBookings) + "\n"
	status += "  Подтверждено: " + itoa(stats.ConfirmedBookings) + "\n"
	status += "  Отменено: " + itoa(stats.CancelledBookings) + "\n"
	status += "\nНомера по типам:\n"
	for t, count := range stats.RoomsByType {
		status += "  " + string(t) + ": " + itoa(count) + "\n"
	}

	return status, nil
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	result := ""
	for n > 0 {
		result = string(rune('0'+n%10)) + result
		n /= 10
	}
	return result
}
