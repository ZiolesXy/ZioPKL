package seeders

import (
	"fmt"
	"net/http"

	"main/database"
	"main/helpers"
	"main/models"
	"main/respons"

	"github.com/gin-gonic/gin"
)

func RunSeed(c *gin.Context) {

	/* =========================
	   1. SEED ROLES (IDEMPOTENT)
	========================= */

	roleNames := []string{
		"Admin",
		"User",
		"Manager",
		"Staff IT",
		"Sales",
		"Inventory",
	}

	var roles []models.Role
	for _, name := range roleNames {
		role := models.Role{Name: name}
		if err := database.DB.
			Where("name = ?", name).
			FirstOrCreate(&role).Error; err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		roles = append(roles, role)
	}

	/* =========================
	   2. SEED BARANG (IDEMPOTENT)
	========================= */

	namaBarang := []string{
		"MacBook Pro M2", "Logitech MX Master 3", "Keychron K2 V2",
		"Dell UltraSharp 27", "Epson L3210 Printer", "Sony WH-1000XM5",
		"iPad Air 5", "Samsung Galaxy S23", "SteelSeries Mousepad",
		"ThinkPad X1 Carbon", "Wacom Intuos Pro", "Cisco Router ISR",
		"Ubiquiti UniFi AP", "GoPro Hero 11", "BenQ SW271 Monitor",
		"WD Black SN850 SSD", "Seagate IronWolf 4TB", "Asus ROG Strix VGA",
		"Corsair RM850 PSU", "Razer BlackWidow",
	}

	var barangs []models.Barang
	for _, name := range namaBarang {
		barang := models.Barang{Name: name}
		if err := database.DB.
			Where("name = ?", name).
			FirstOrCreate(&barang).Error; err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		barangs = append(barangs, barang)
	}

	/* =========================
	   3. SEED USERS (EMAIL UNIQUE)
	========================= */

	type UserSeed struct {
		Name  string
		Email string
		Role  string
	}

	userSeeds := []UserSeed{
		{"Pasha", "pashaprabasakti@gmail.com", "Admin"},
		{"Acheron", "pasha@test.com", "User"},
		{"Siti", "siti@test.com", "Manager"},
		{"Budi", "budi@test.com", "Sales"},
		{"Dewi", "dewi@test.com", "Staff IT"},
		{"Eko", "eko@test.com", "Inventory"},
		{"Farhan", "farhan@test.com", "User"},
		{"Gita", "gita@test.com", "User"},
		{"Hadi", "hadi@test.com", "User"},
		{"Indah", "indah@test.com", "User"},
	}

	var users []models.User

	for _, u := range userSeeds {
		var role models.Role
		if err := database.DB.
			Where("name = ?", u.Role).
			First(&role).Error; err != nil {
			c.JSON(400, gin.H{"error": "role not found"})
			return
		}

		hashed, _ := helpers.HashPassword("123456")

		user := models.User{
			Name:     u.Name,
			Email:    u.Email,
			Password: hashed,
			RoleID:   role.ID,
		}

		if err := database.DB.
			Where("email = ?", u.Email).
			FirstOrCreate(&user).Error; err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		users = append(users, user)
	}

	/* =========================
	   4. SEED OWNERSHIP (AMAN)
	========================= */

	for i := range users {
		var barangUntukUser []models.Barang

		for j := 0; j < 3; j++ {
			idx := (i + j) % len(barangs)
			barangUntukUser = append(barangUntukUser, barangs[idx])
		}

		if err := database.DB.
			Model(&users[i]).
			Association("Barangs").
			Replace(barangUntukUser); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
	}

	/* =========================
	   5. RESPONSE
	========================= */

	message := fmt.Sprintf(
		"Seeder sukses: %d Role, %d User, %d Barang",
		len(roles),
		len(users),
		len(barangs),
	)

	c.JSON(
		http.StatusOK,
		respons.NewJsonResponse("success", message),
	)
}
