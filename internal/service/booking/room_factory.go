package booking

import (
	"github.com/YurcheuskiRadzivon/booking-system/internal/models/booking"
)

type RoomFactory struct{}

func NewRoomFactory() *RoomFactory {
	return &RoomFactory{}
}

func (f *RoomFactory) CreateRoom(roomType booking.RoomType, roomNumber string) *booking.Room {
	switch roomType {
	case booking.RoomTypeStandard:
		return f.CreateStandardRoom(roomNumber)
	case booking.RoomTypeDeluxe:
		return f.CreateDeluxeRoom(roomNumber)
	case booking.RoomTypeSuite:
		return f.CreateSuite(roomNumber)
	case booking.RoomTypeFamily:
		return f.CreateFamilyRoom(roomNumber)
	default:
		return f.CreateStandardRoom(roomNumber)
	}
}

func (f *RoomFactory) CreateStandardRoom(roomNumber string) *booking.Room {
	return &booking.Room{
		RoomNumber:  roomNumber,
		RoomType:    booking.RoomTypeStandard,
		BasePrice:   2500.0,
		Capacity:    2,
		Status:      booking.RoomStatusAvailable,
		Description: "Стандартный номер " + roomNumber,
	}
}

func (f *RoomFactory) CreateDeluxeRoom(roomNumber string) *booking.Room {
	return &booking.Room{
		RoomNumber:  roomNumber,
		RoomType:    booking.RoomTypeDeluxe,
		BasePrice:   4500.0,
		Capacity:    3,
		Status:      booking.RoomStatusAvailable,
		Description: "Делюкс номер " + roomNumber + " с видом на море",
	}
}

func (f *RoomFactory) CreateSuite(roomNumber string) *booking.Room {
	return &booking.Room{
		RoomNumber:  roomNumber,
		RoomType:    booking.RoomTypeSuite,
		BasePrice:   7500.0,
		Capacity:    4,
		Status:      booking.RoomStatusAvailable,
		Description: "Люкс номер " + roomNumber + " с гостиной и джакузи",
	}
}

func (f *RoomFactory) CreateFamilyRoom(roomNumber string) *booking.Room {
	return &booking.Room{
		RoomNumber:  roomNumber,
		RoomType:    booking.RoomTypeFamily,
		BasePrice:   5500.0,
		Capacity:    6,
		Status:      booking.RoomStatusAvailable,
		Description: "Семейный номер " + roomNumber + " с двумя спальнями",
	}
}
