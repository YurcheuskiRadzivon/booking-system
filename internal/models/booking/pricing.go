package booking

import (
	"time"
)

type AlgorithmType string

const (
	AlgorithmTypeRegular  AlgorithmType = "regular"
	AlgorithmTypeWeekend  AlgorithmType = "weekend"
	AlgorithmTypeSeasonal AlgorithmType = "seasonal"
	AlgorithmTypeSpecial  AlgorithmType = "special"
)

type PricingAlgorithm struct {
	ID            int64         `json:"id" db:"id"`
	Date          time.Time     `json:"date" db:"date"`
	AlgorithmType AlgorithmType `json:"algorithm_type" db:"algorithm_type"`
}

type SpecialDate struct {
	ID          int64     `json:"id" db:"id"`
	Date        time.Time `json:"date" db:"date"`
	Name        string    `json:"name" db:"name"`
	Coefficient float64   `json:"coefficient" db:"coefficient"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

type CreateSpecialDateRequest struct {
	Date        string  `json:"date"`
	Name        string  `json:"name"`
	Coefficient float64 `json:"coefficient"`
}

type PriceCalculationRequest struct {
	RoomID   int64     `json:"room_id"`
	CheckIn  time.Time `json:"check_in"`
	CheckOut time.Time `json:"check_out"`
}

type PriceCalculationResponse struct {
	BasePrice      float64        `json:"base_price"`
	TotalPrice     float64        `json:"total_price"`
	Nights         int            `json:"nights"`
	DailyBreakdown []DayPriceInfo `json:"daily_breakdown"`
}

type DayPriceInfo struct {
	Date        string  `json:"date"`
	BasePrice   float64 `json:"base_price"`
	Coefficient float64 `json:"coefficient"`
	Reason      string  `json:"reason"`
	DayPrice    float64 `json:"day_price"`
}
