package seeders

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"
	"voca-plane/internal/domain/models"
	"voca-plane/pkg/helper"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func SeedAll(db *gorm.DB) {
	log.Println(">>> STARTING FULL SEEDING...")
	SeedUsers(db)
	SeedAirlines(db)
	SeedAirports(db)
	SeedFlights(db)
	SeedPromos(db)
	log.Println(">>> SEEDING COMPLETED")
}

func SeedMain(db *gorm.DB) {
	log.Println(">>> STARTING FULL SEEDING...")
	SeedUsers(db)
	log.Println(">>> SEEDING COMPLETED")
}

func DropAll(db *gorm.DB) {
    log.Println(">>> DROPPING ALL TABLES...")

    err := db.Migrator().DropTable(
        &models.FlightSeat{},
        &models.FlightClass{},
        &models.Seat{},
        &models.Flight{},
        &models.Airline{},
        &models.Airport{},
        &models.TransactionItem{},
        &models.Transaction{},
        &models.PromoCode{},
        &models.User{},
    )

    if err != nil {
        log.Fatalf(">>> FAILED TO DROP TABLES: %v", err)
    }
    
    log.Println(">>> ALL TABLES DROPPED SUCCESSFULLY")
}

func ResetDatabase(db *gorm.DB) {
    DropAll(db)

    log.Println(">>> RE-MIGRATING TABLES...")
    db.AutoMigrate(
        &models.User{},
        &models.Airline{},
        &models.Airport{},
        &models.Flight{},
        &models.FlightClass{},
        &models.Seat{},
        &models.FlightSeat{},
        &models.Transaction{},
        &models.TransactionItem{},
        &models.PromoCode{},
    )

    // 3. Isi ulang data
    SeedAll(db)
}

func SeedUsers(db *gorm.DB) {
	log.Println(">>> Seeding Users...")
	
	// Super Admin
	superAdminPwd, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	superAdmin := models.User{
		Name:     "Super Admin",
		Email:    "superadmin@flightbooking.com",
		Password: string(superAdminPwd),
		Role:     models.RoleSuperAdmin,
	}
	db.FirstOrCreate(&superAdmin, models.User{Email: superAdmin.Email})

	// Admin
	adminPwd, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	admin := models.User{
		Name:     "Admin User",
		Email:    "admin@flightbooking.com",
		Password: string(adminPwd),
		Role:     models.RoleAdmin,
	}
	db.FirstOrCreate(&admin, models.User{Email: admin.Email})

	// Regular User
	userPwd, _ := bcrypt.GenerateFromPassword([]byte("user123"), bcrypt.DefaultCost)
	user := models.User{
		Name:     "John Doe",
		Email:    "user@flightbooking.com",
		Password: string(userPwd),
		Role:     models.RoleUser,
	}
	db.FirstOrCreate(&user, models.User{Email: user.Email})
}

func SeedAirlines(db *gorm.DB) {
	log.Println(">>> Seeding Airlines from JSON...")
	data, err := os.ReadFile("internal/seeders/airlines.json")
	if err != nil {
		log.Printf("Failed to read airlines.json: %v", err)
		return
	}

	var airlines []models.Airline
	if err := json.Unmarshal(data, &airlines); err != nil {
		log.Printf("Failed to unmarshal airlines.json: %v", err)
		return
	}

	for _, a := range airlines {
		db.FirstOrCreate(&a, models.Airline{Code: a.Code})
	}
}

func SeedAirports(db *gorm.DB) {
	log.Println(">>> Seeding Airports from JSON...")
	data, err := os.ReadFile("internal/seeders/airports.json")
	if err != nil {
		log.Printf("Failed to read airports.json: %v", err)
		return
	}

	var airportMap map[string]struct {
		IATA string `json:"iata"`
		Name string `json:"name"`
		City string `json:"city"`
	}
	if err := json.Unmarshal(data, &airportMap); err != nil {
		log.Printf("Failed to unmarshal airports.json: %v", err)
		return
	}

	// We only seed a subset if it's too many, or only those with IATA code
	count := 0
	for _, v := range airportMap {
		if v.IATA == "" {
			continue
		}
		airport := models.Airport{
			Code: v.IATA,
			Name: v.Name,
			City: v.City,
		}
		db.FirstOrCreate(&airport, models.Airport{Code: airport.Code})
		count++
		if count >= 100 { // Limit to 100 for now to avoid massive DB bloat if not needed
			break
		}
	}
}

func SeedFlights(db *gorm.DB) {
	i := 0
	var airlines []models.Airline
	var airports []models.Airport

	db.Find(&airlines)
	db.Find(&airports)

	if len(airlines) == 0 || len(airports) < 2 {
		log.Println("Insufficient airlines or airports for flight seeding")
		return
	}

	type classConfig struct {
		ClassType string
		Price     float64
		Percent   float64
	}

	for i = 1; i <= 100; i++ {
		airline := airlines[i%len(airlines)]
		origin := airports[i%len(airports)]
		dest := airports[(i+1)%len(airports)]

		// Konfigurasi kursi per pesawat
		totalSeats := 20
		totalRows := 5
		totalColumns := 4

		// Logika Variasi Kelas:
		// i 1-7   : Full Economy
		// i 8-14  : Business & Economy
		// i 15-20 : First, Business, & Economy
		var classes []classConfig
		switch {
		default:
			classes = []classConfig{
				{"First", 5000000 + float64(i*150000), 0.20},
				{"Business", 3000000 + float64(i*100000), 0.30},
				{"Economy", 1000000 + float64(i*50000), 0.50},
			}
		}

		flight := models.Flight{
			AirlineID:     airline.ID,
			OriginID:      origin.ID,
			DestinationID: dest.ID,
			DepartureTime: time.Now().Add(time.Duration(i*6) * time.Hour),
			ArrivalTime:   time.Now().Add(time.Duration(i*6+2) * time.Hour),
			FlightNumber:  fmt.Sprintf("FL%03d", i),
			TotalSeats:    totalSeats,
			TotalRows:     totalRows,
			TotalColumns:  totalColumns,
		}

		// Gunakan Where agar flight.ID terisi ke struct meskipun data sudah ada (FirstOrCreate)
		if err := db.Where(models.Flight{FlightNumber: flight.FlightNumber}).FirstOrCreate(&flight).Error; err != nil {
			log.Printf("Failed to seed flight %s: %v", flight.FlightNumber, err)
			continue
		}

		// Build class allocations for seat generation
		var allocations []helper.ClassAlloc

		// Create flight classes and build allocations
		for _, c := range classes {
			seatCount := int(float64(totalSeats) * c.Percent)

			fClass := models.FlightClass{
				FlightID:  flight.ID,
				ClassType: c.ClassType,
				Price:     c.Price,
			}
			
			// Pastikan Class terbuat dan ID-nya didapat
			db.Where(models.FlightClass{FlightID: flight.ID, ClassType: fClass.ClassType}).FirstOrCreate(&fClass)

			allocations = append(allocations, helper.ClassAlloc{
				ClassType: c.ClassType,
				Price:     c.Price,
				SeatCount: seatCount,
			})
		}

		// Fix last class allocation to take remaining seats
		if len(allocations) > 0 {
			used := 0
			for i := 0; i < len(allocations)-1; i++ {
				used += allocations[i].SeatCount
			}
			allocations[len(allocations)-1].SeatCount = totalSeats - used
		}

		// Generate seat codes (e.g. "1A", "1B", ..., "5D")
		seatCodes := helper.GenerateSeatCodes(totalRows, totalColumns)

		// Upsert master Seat rows
		seatModels := make([]models.Seat, len(seatCodes))
		for j, code := range seatCodes {
			seatModels[j] = models.Seat{SeatCode: code}
			db.Where(models.Seat{SeatCode: code}).FirstOrCreate(&seatModels[j])
		}

		// Order seats by code to match seatCodes order
		seatMap := make(map[string]models.Seat)
		for _, s := range seatModels {
			seatMap[s.SeatCode] = s
		}
		orderedSeats := make([]models.Seat, 0, len(seatCodes))
		for _, code := range seatCodes {
			if s, ok := seatMap[code]; ok {
				orderedSeats = append(orderedSeats, s)
			}
		}

		// Generate FlightSeat pivot entries distributed across allocations
		flightSeats := helper.GenerateFlightSeatModels(flight.ID, orderedSeats, allocations)

		for _, fs := range flightSeats {
			db.Where(models.FlightSeat{FlightID: fs.FlightID, SeatID: fs.SeatID}).FirstOrCreate(&fs)
		}
	}
	log.Printf(">>> Seeding %v Sample Flights with Diverse Classes...", i-1)
}

func SeedPromos(db *gorm.DB) {
	log.Println(">>> Seeding Promos...")
	promos := []models.PromoCode{
		{Code: "HEMAT50", Discount: 50, IsActive: true},
		{Code: "NEWUSER", Discount: 20, IsActive: true},
	}
	for _, p := range promos {
		db.FirstOrCreate(&p, models.PromoCode{Code: p.Code})
	}
}

func InitSeeders(db *gorm.DB) {
	var count int64
	db.Model(&models.User{}).Count(&count)
	if count == 0 {
		SeedMain(db)
	}
}
