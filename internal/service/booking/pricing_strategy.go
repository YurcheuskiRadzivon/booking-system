package booking

import (
	"context"
	"time"

	"github.com/YurcheuskiRadzivon/booking-system/internal/models/booking"
)

const (
	CoeffWeekend      = 1.25
	CoeffSummerSpring = 1.3
	CoeffWinterFall   = 0.9
)

type PriceCalculator struct {
	specialDates map[string]booking.SpecialDate
}

func NewPriceCalculator(specialDates []booking.SpecialDate) *PriceCalculator {
	dateMap := make(map[string]booking.SpecialDate)
	for _, sd := range specialDates {
		key := sd.Date.Format("2006-01-02")
		dateMap[key] = sd
	}
	return &PriceCalculator{specialDates: dateMap}
}

func (pc *PriceCalculator) CalculateTotalPrice(basePrice float64, checkIn, checkOut time.Time) booking.PriceCalculationResponse {
	nights := int(checkOut.Sub(checkIn).Hours() / 24)
	if nights <= 0 {
		nights = 1
	}

	breakdown := make([]booking.DayPriceInfo, 0, nights)
	totalPrice := 0.0
	currentDate := checkIn

	for i := 0; i < nights; i++ {
		dayInfo := pc.calculateDayPrice(basePrice, currentDate)
		breakdown = append(breakdown, dayInfo)
		totalPrice += dayInfo.DayPrice
		currentDate = currentDate.AddDate(0, 0, 1)
	}

	return booking.PriceCalculationResponse{
		BasePrice:      basePrice,
		TotalPrice:     totalPrice,
		Nights:         nights,
		DailyBreakdown: breakdown,
	}
}

func (pc *PriceCalculator) calculateDayPrice(basePrice float64, date time.Time) booking.DayPriceInfo {
	coefficient := 1.0
	reasons := []string{}

	dateKey := date.Format("2006-01-02")
	if special, ok := pc.specialDates[dateKey]; ok {
		return booking.DayPriceInfo{
			Date:        dateKey,
			BasePrice:   basePrice,
			Coefficient: special.Coefficient,
			Reason:      special.Name,
			DayPrice:    basePrice * special.Coefficient,
		}
	}

	if date.Weekday() == time.Saturday || date.Weekday() == time.Sunday {
		coefficient *= CoeffWeekend
		reasons = append(reasons, "Vyhodnoy")
	}

	month := date.Month()
	if month >= time.April && month <= time.September {
		coefficient *= CoeffSummerSpring
		reasons = append(reasons, "Vysokiy sezon")
	} else {
		coefficient *= CoeffWinterFall
		reasons = append(reasons, "Nizkiy sezon")
	}

	reason := "Obychnyy den"
	if len(reasons) > 0 {
		reason = ""
		for i, r := range reasons {
			if i > 0 {
				reason += ", "
			}
			reason += r
		}
	}

	return booking.DayPriceInfo{
		Date:        dateKey,
		BasePrice:   basePrice,
		Coefficient: coefficient,
		Reason:      reason,
		DayPrice:    basePrice * coefficient,
	}
}

type PriceService interface {
	CalculatePrice(ctx context.Context, roomID int64, checkIn, checkOut time.Time) (*booking.PriceCalculationResponse, error)
}
