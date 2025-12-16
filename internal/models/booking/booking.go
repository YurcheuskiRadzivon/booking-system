package booking

import (
	"time"
)

type BookingStatus string

const (
	BookingStatusPending   BookingStatus = "pending"
	BookingStatusConfirmed BookingStatus = "confirmed"
	BookingStatusCancelled BookingStatus = "cancelled"
)

type GuestInfo struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
}

type Booking struct {
	ID        int64         `json:"id" db:"id"`
	StartDate time.Time     `json:"start_date" db:"start_date"`
	EndDate   time.Time     `json:"end_date" db:"end_date"`
	RoomID    int64         `json:"room_id" db:"room_id"`
	GuestInfo GuestInfo     `json:"guest_info" db:"guest_info"`
	Price     float64       `json:"price" db:"price"`
	Status    BookingStatus `json:"status" db:"status"`
	CreatedAt time.Time     `json:"created_at" db:"created_at"`
	UpdatedAt time.Time     `json:"updated_at" db:"updated_at"`
}

type BookingWithRoom struct {
	Booking
	Room Room `json:"room"`
}

type CreateBookingRequest struct {
	RoomID    int64     `json:"room_id"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	GuestInfo GuestInfo `json:"guest_info"`
}

type BookingResponse struct {
	Booking
	Room   Room `json:"room"`
	Nights int  `json:"nights"`
}
