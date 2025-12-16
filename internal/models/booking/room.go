package booking

import (
	"time"
)

type RoomType string

const (
	RoomTypeStandard RoomType = "standard"
	RoomTypeDeluxe   RoomType = "deluxe"
	RoomTypeSuite    RoomType = "suite"
	RoomTypeFamily   RoomType = "family"
)

type RoomStatus string

const (
	RoomStatusAvailable   RoomStatus = "available"
	RoomStatusOccupied    RoomStatus = "occupied"
	RoomStatusMaintenance RoomStatus = "maintenance"
)

type Room struct {
	ID          int64      `json:"id" db:"id"`
	RoomNumber  string     `json:"room_number" db:"room_number"`
	RoomType    RoomType   `json:"room_type" db:"room_type"`
	BasePrice   float64    `json:"base_price" db:"base_price"`
	Capacity    int        `json:"capacity" db:"capacity"`
	Status      RoomStatus `json:"status" db:"status"`
	Description string     `json:"description" db:"description"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
}

type RoomTypeInfo struct {
	ID        int64   `json:"id" db:"id"`
	Name      string  `json:"name" db:"name"`
	BasePrice float64 `json:"base_price" db:"base_price"`
	Breakfast bool    `json:"breakfast" db:"breakfast"`
	Lunch     bool    `json:"lunch" db:"lunch"`
	Dinner    bool    `json:"dinner" db:"dinner"`
	FastWifi  bool    `json:"fast_wifi" db:"fast_wifi"`
	Pool      bool    `json:"pool" db:"pool"`
	Gym       bool    `json:"gym" db:"gym"`
}

type RoomSearchRequest struct {
	CheckIn  time.Time `json:"check_in"`
	CheckOut time.Time `json:"check_out"`
	RoomType RoomType  `json:"room_type,omitempty"`
	Capacity int       `json:"capacity,omitempty"`
}

type RoomWithAvailability struct {
	Room        Room    `json:"room"`
	IsAvailable bool    `json:"is_available"`
	TotalPrice  float64 `json:"total_price,omitempty"`
}
