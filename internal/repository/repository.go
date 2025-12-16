package repository

import (
	"context"
	"time"

	"github.com/YurcheuskiRadzivon/booking-system/internal/models/booking"
	"github.com/YurcheuskiRadzivon/booking-system/internal/models/notification"
)

type Repository interface {
	Room() RoomRepository
	Booking() BookingRepository
	Notification() NotificationRepository
	SpecialDate() SpecialDateRepository
	Close() error
}

type RoomRepository interface {
	GetAll(ctx context.Context) ([]booking.Room, error)
	GetByID(ctx context.Context, id int64) (*booking.Room, error)
	GetByNumber(ctx context.Context, roomNumber string) (*booking.Room, error)
	GetAvailable(ctx context.Context, checkIn, checkOut time.Time) ([]booking.Room, error)
	GetAvailableByType(ctx context.Context, roomType booking.RoomType, checkIn, checkOut time.Time) ([]booking.Room, error)
	GetAvailableByCapacity(ctx context.Context, capacity int, checkIn, checkOut time.Time) ([]booking.Room, error)
	Create(ctx context.Context, room *booking.Room) error
	Update(ctx context.Context, room *booking.Room) error
	Delete(ctx context.Context, id int64) error
	UpdateStatus(ctx context.Context, id int64, status booking.RoomStatus) error
}

type BookingRepository interface {
	GetAll(ctx context.Context) ([]booking.Booking, error)
	GetAllWithRooms(ctx context.Context) ([]booking.BookingWithRoom, error)
	GetByID(ctx context.Context, id int64) (*booking.Booking, error)
	GetByRoomID(ctx context.Context, roomID int64) ([]booking.Booking, error)
	GetByStatus(ctx context.Context, status booking.BookingStatus) ([]booking.Booking, error)
	GetByStatusWithRooms(ctx context.Context, status booking.BookingStatus) ([]booking.BookingWithRoom, error)
	GetByEmail(ctx context.Context, email string) ([]booking.Booking, error)
	GetActiveForRoom(ctx context.Context, roomID int64, checkIn, checkOut time.Time) ([]booking.Booking, error)
	Create(ctx context.Context, b *booking.Booking) error
	Update(ctx context.Context, b *booking.Booking) error
	UpdateStatus(ctx context.Context, id int64, status booking.BookingStatus) error
	Delete(ctx context.Context, id int64) error
	IsRoomAvailable(ctx context.Context, roomID int64, checkIn, checkOut time.Time) (bool, error)
}

type NotificationRepository interface {
	GetAll(ctx context.Context) ([]notification.NotificationType, error)
	GetByID(ctx context.Context, id int64) (*notification.NotificationType, error)
	GetByName(ctx context.Context, name string) (*notification.NotificationType, error)
	Create(ctx context.Context, nt *notification.NotificationType) error
	Update(ctx context.Context, nt *notification.NotificationType) error
	Delete(ctx context.Context, id int64) error
}

type SpecialDateRepository interface {
	GetAll(ctx context.Context) ([]booking.SpecialDate, error)
	GetByID(ctx context.Context, id int64) (*booking.SpecialDate, error)
	GetByDate(ctx context.Context, date time.Time) (*booking.SpecialDate, error)
	GetByDateRange(ctx context.Context, start, end time.Time) ([]booking.SpecialDate, error)
	Create(ctx context.Context, sd *booking.SpecialDate) error
	Update(ctx context.Context, sd *booking.SpecialDate) error
	Delete(ctx context.Context, id int64) error
}
