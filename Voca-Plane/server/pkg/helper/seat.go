package helper

import (
	"errors"
	"fmt"
	"strings"
	"voca-plane/internal/domain/dto/request"
	"voca-plane/internal/domain/models"
)

const (
	FirstClassRatio    = 0.20
	BusinessClassRatio = 0.30
)

type ClassAlloc struct {
	ClassType string
	Price     float64
	SeatCount int
}

func ValidateFlightInput(flight *models.Flight, classCount int, classPrice []request.ClassPriceRequest) error {
	const MaxSeatRows = 7
	maxCapacity := flight.TotalRows * flight.TotalColumns

	if flight.TotalRows > MaxSeatRows {
		return fmt.Errorf("total_rows cannot exceed 8")
	}

	if flight.TotalSeats != maxCapacity {
		return fmt.Errorf("total_seats must equal rows × columns (%d), got %d", maxCapacity, flight.TotalSeats,)
	}

	if classCount < 1 || classCount > 3 {
		return errors.New("class_count must be 1, 2, or 3")
	}

	if len(classPrice) != classCount {
		return fmt.Errorf("class_prices count (%d) must match class_count (%d)", len(classPrice), classCount)
	}

	return nil
}

func MapClassPrices(classPrices []request.ClassPriceRequest) map[string]float64 {
	classMap := make(map[string]float64)

	for _, cp := range classPrices {
		classType := strings.ToLower(cp.ClassType)
		classMap[classType] = cp.Price
	}

	return classMap
}

func CalculateSeatAllocation(totalSeats int, classMap map[string]float64, classCount int) []ClassAlloc {
	var allocations []ClassAlloc

	switch classCount {

	case 1:
		for classType, price := range classMap {
			allocations = append(allocations, ClassAlloc{
				ClassType: classType,
				Price:     price,
				SeatCount: totalSeats,
			})
		}

	case 2:
		var types []string
		for classType := range classMap {
			types = append(types, classType)
		}

		seatsA := int(float64(totalSeats) * BusinessClassRatio)
		seatsB := totalSeats - seatsA

		allocations = append(allocations,
			ClassAlloc{types[0], classMap[types[0]], seatsA},
			ClassAlloc{types[1], classMap[types[1]], seatsB},
		)

	case 3:
		firstSeats := int(float64(totalSeats) * FirstClassRatio)
		businessSeats := int(float64(totalSeats) * BusinessClassRatio)
		economySeats := totalSeats - firstSeats - businessSeats

		allocations = append(allocations,
			ClassAlloc{"first", classMap["first"], firstSeats},
			ClassAlloc{"business", classMap["business"], businessSeats},
			ClassAlloc{"economy", classMap["economy"], economySeats},
		)
	}

	return allocations
}

// GenerateSeatCodes generates seat codes using row numbers (1-based) and column letters (A-based).
// Example: rows=9, columns=6 → ["1A","1B","1C","1D","1E","1F","2A",...,"9F"]
func GenerateSeatCodes(rows, columns int) []string {
	codes := make([]string, 0, rows*columns)
	for r := 1; r <= rows; r++ {
		for c := 0; c < columns; c++ {
			letter := string(rune('A' + c))
			codes = append(codes, fmt.Sprintf("%d%s", r, letter))
		}
	}
	return codes
}

// GenerateFlightSeatModels creates FlightSeat pivot entries for a given flight,
// distributing seats across class allocations in order.
func GenerateFlightSeatModels(flightID uint, seats []models.Seat, allocations []ClassAlloc) []models.FlightSeat {
	flightSeats := make([]models.FlightSeat, 0, len(seats))

	seatIndex := 0
	for _, alloc := range allocations {
		for i := 0; i < alloc.SeatCount && seatIndex < len(seats); i++ {
			flightSeats = append(flightSeats, models.FlightSeat{
				FlightID:    flightID,
				SeatID:      seats[seatIndex].ID,
				ClassType:   alloc.ClassType,
				IsAvailable: true,
			})
			seatIndex++
		}
	}

	return flightSeats
}
