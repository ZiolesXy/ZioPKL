package seeders

import (
	// "errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
	"voca-store/internal/domain/models"
	"voca-store/internal/helper"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func SeedRoles(db *gorm.DB) error {
	roles := []models.Role{
		{Name: "Admin"},
		{Name: "User"},
	}

	return db.Clauses(clause.OnConflict{DoNothing: true}).Create(&roles).Error
}

func SeedCategories(db *gorm.DB) error {

	names := []string{
		"Laptop",
		"Smartphone",
		"Accessories",
		"Makanan",
		"Minuman",
		"Pertanian",
		"Mainan",
		"Lainnya",
	}

	var categories []models.Category

	for _, name := range names {

		slug, err := helper.GenerateUniqueCategorySlug(db, name)

		if err != nil {
			return err
		}

		categories = append(categories, models.Category{
			Name: name,
			Slug: slug,
		})

	}

	return db.
		Clauses(clause.OnConflict{DoNothing: true}).
		Create(&categories).Error
}

func SeedAdmin(db *gorm.DB) error {
	var adminRole models.Role
	if err := db.Where("name = ?", "Admin").First(&adminRole).Error; err != nil {
		return err
	}

	var existingAdmin models.User
	if err := db.Where("email = ?", "pashaprabasakti@gmail.com").First(&existingAdmin).Error; err == nil {
		return nil // already exists
	}

	hashedPassword, err := helper.HashPassword("360589")
	if err != nil {
		return err
	}

	admin := models.User{
		Name:            "Pasha",
		Email:           "pashaprabasakti@gmail.com",
		Password:        hashedPassword,
		TelephoneNumber: "081234567890",
		RoleID:          adminRole.ID,
	}

	if err := db.Create(&admin).Error; err != nil {
		return err
	}

	adminAddresses := []struct {
		label       string
		addressLine string
		city        string
		province    string
		postalCode  string
	}{
		{"Rumah Dinas", "Jl. Menteri No. 1", "Jakarta Pusat", "DKI Jakarta", "10110"},
		{"Kantor Pusat", "Gedung Voca Lantai 1", "Jakarta Selatan", "DKI Jakarta", "12190"},
		{"Rumah Pribadi", "Jl. Boulevard Indah 5", "Tangerang", "Banten", "15810"},
		{"Gudang Utama", "Kawasan Industri MM2100", "Bekasi", "Jawa Barat", "17530"},
		{"Drop Point", "Jl. Ahmad Yani No. 10", "Semarang", "Jawa Tengah", "50131"},
	}

	for i, addr := range adminAddresses {
		uid, err := helper.NewGenerateAddressUID(db)
		if err != nil {
			return err
		}

		address := models.Address{
			UID:           uid,
			UserID:        admin.ID,
			Label:         addr.label,
			RecipientName: admin.Name,
			Phone:         admin.TelephoneNumber,
			AddressLine:   addr.addressLine,
			City:          addr.city,
			Province:      addr.province,
			PostalCode:    addr.postalCode,
			IsPrimary:     i == 0,
		}

		if err := db.Create(&address).Error; err != nil {
			return err
		}
	}

	return nil
}

func SeedUsers(db *gorm.DB) error {
	var userRole models.Role
	if err := db.Where("name = ?", "User").First(&userRole).Error; err != nil {
		return err
	}

	users := []struct {
		name      string
		email     string
		pass      string
		telp      string
		addresses []struct {
			label       string
			addressLine string
			city        string
			province    string
			postalCode  string
		}
	}{
		{
			name:  "ZiolesXy",
			email: "eaglegaming3605@gmail.com",
			pass:  "123456",
			telp:  "081584012152",
			addresses: []struct {
				label       string
				addressLine string
				city        string
				province    string
				postalCode  string
			}{
				{"Kantor Pusat", "Jl. Jend. Sudirman Kav. 52-53", "Jakarta Selatan", "DKI Jakarta", "12190"},
				{"Gudang Logistik", "Kawasan Industri Bukit Indah Blok C2", "Purwakarta", "Jawa Barat", "41181"},
				{"Villa", "Jl. Labuan Sait No. 10", "Pecatu", "Bali", "80361"},
				{"Kost", "Jl. Dipati Ukur No. 45", "Bandung", "Jawa Barat", "40132"},
				{"Rumah Kebun", "Jl. Raya Trans Sulawesi KM 15", "Maros", "Sulawesi Selatan", "90552"},
			},
		},
		{
			name:  "Amtir",
			email: "petir@gmail.com",
			pass:  "123456",
			telp:  "081111111111",
			addresses: []struct {
				label       string
				addressLine string
				city        string
				province    string
				postalCode  string
			}{
				{"Rumah", "Jl. Petir Utama No. 1", "Semarang", "Jawa Tengah", "50123"},
				{"Kantor", "Jl. Industri Modern I", "Jakarta", "DKI Jakarta", "10110"},
				{"Rumah Orang Tua", "Jl. Desa Sejahtera 44", "Solo", "Jawa Tengah", "57123"},
				{"Kost", "Jl. Tembalang No. 5", "Semarang", "Jawa Tengah", "50275"},
				{"Villa", "Jl. Puncak Indah 7", "Bogor", "Jawa Barat", "16750"},
			},
		},
		{
			name:  "Jane Smith",
			email: "siti@gmail.com",
			pass:  "123456",
			telp:  "082222222222",
			addresses: []struct {
				label       string
				addressLine string
				city        string
				province    string
				postalCode  string
			}{
				{"Rumah", "Jl. Mawar Merah No. 12", "Surabaya", "Jawa Timur", "60111"},
				{"Apartment", "Pakuwon City Lt. 15", "Surabaya", "Jawa Timur", "60112"},
				{"Kantor Surabaya", "Jl. HR Muhammad No. 1", "Surabaya", "Jawa Timur", "60226"},
				{"Rumah Nenek", "Jl. Pahlawan 3", "Malang", "Jawa Timur", "65111"},
				{"Workshop", "Jl. Gresik Industri 5", "Gresik", "Jawa Timur", "61121"},
			},
		},
		{
			name:  "Bob Johnson",
			email: "bob@gmail.com",
			pass:  "123456",
			telp:  "083333333333",
			addresses: []struct {
				label       string
				addressLine string
				city        string
				province    string
				postalCode  string
			}{
				{"Headquarters", "Jl. Gatot Subroto No. 50", "Jakarta", "DKI Jakarta", "12710"},
				{"Warehouse", "Jl. Marunda Center No. 8", "Bekasi", "Jawa Barat", "17111"},
				{"Private Studio", "Jl. Kemang Timur 12", "Jakarta", "DKI Jakarta", "12730"},
				{"Family Home", "Jl. Menteng No. 1", "Jakarta", "DKI Jakarta", "10310"},
				{"Guest House", "Jl. Braga No. 20", "Bandung", "Jawa Barat", "40111"},
			},
		},
	}

	for _, u := range users {
		var existingUser models.User
		if err := db.Where("email = ?", u.email).First(&existingUser).Error; err == nil {
			continue
		}

		hashedPassword, err := helper.HashPassword(u.pass)
		if err != nil {
			return err
		}

		user := models.User{
			Name:            u.name,
			Email:           u.email,
			Password:        hashedPassword,
			TelephoneNumber: u.telp,
			RoleID:          userRole.ID,
		}

		if err := db.Create(&user).Error; err != nil {
			return err
		}

		for i, addr := range u.addresses {
			uid, err := helper.NewGenerateAddressUID(db)
			if err != nil {
				return err
			}

			address := models.Address{
				UID:           uid,
				UserID:        user.ID,
				Label:         addr.label,
				RecipientName: u.name,
				Phone:         u.telp,
				AddressLine:   addr.addressLine,
				City:          addr.city,
				Province:      addr.province,
				PostalCode:    addr.postalCode,
				IsPrimary:     i == 0,
			}

			if err := db.Create(&address).Error; err != nil {
				return err
			}
		}
	}

	return nil
}

func SeedProductsOld(db *gorm.DB) error {

	// ================= AMBIL CATEGORY =================

	var laptop, smartphone, accessories, makanan, minuman, pertanian, mainan, lainnya models.Category

	db.Where("name = ?", "Laptop").First(&laptop)
	db.Where("name = ?", "Smartphone").First(&smartphone)
	db.Where("name = ?", "Accessories").First(&accessories)
	db.Where("name = ?", "Makanan").First(&makanan)
	db.Where("name = ?", "Minuman").First(&minuman)
	db.Where("name = ?", "Pertanian").First(&pertanian)
	db.Where("name = ?", "Mainan").First(&mainan)
	db.Where("name = ?", "Lainnya").First(&lainnya)

	// ================= DATA PRODUK =================

	products := []models.Product{

		// ===== LAPTOP =====
		{Name: "MacBook Air M2", Price: 18500000, Stock: 10, ImageURL: "https://images.unsplash.com/photo-1517336714731-489689fd1ca8", CategoryID: laptop.ID},
		{Name: "MacBook Pro M3", Price: 32000000, Stock: 10, ImageURL: "https://images.unsplash.com/photo-1515879218367-8466d910aaa4", CategoryID: laptop.ID},
		{Name: "ASUS ROG Zephyrus G14", Price: 25000000, Stock: 10, ImageURL: "https://images.unsplash.com/photo-1593642532400-2682810df593", CategoryID: laptop.ID},
		{Name: "Lenovo Legion 5", Price: 21000000, Stock: 10, ImageURL: "https://images.unsplash.com/photo-1603302576837-37561b2e2302", CategoryID: laptop.ID},
		{Name: "Acer Aspire 5", Price: 7500000, Stock: 10, ImageURL: "https://images.unsplash.com/photo-1588872657578-7efd1f1555ed", CategoryID: laptop.ID},
		{Name: "HP Pavilion 14", Price: 9000000, Stock: 10, ImageURL: "https://images.unsplash.com/photo-1496181133206-80ce9b88a853", CategoryID: laptop.ID},
		{Name: "Dell XPS 13", Price: 24000000, Stock: 10, ImageURL: "https://images.unsplash.com/photo-1527443224154-c4a3942d3acf", CategoryID: laptop.ID},
		{Name: "MSI Katana GF66", Price: 17000000, Stock: 10, ImageURL: "https://images.unsplash.com/photo-1611078489935-0cb964de46d6", CategoryID: laptop.ID},
		{Name: "ASUS Vivobook 15", Price: 8000000, Stock: 10, ImageURL: "https://images.unsplash.com/photo-1518770660439-4636190af475", CategoryID: laptop.ID},
		{Name: "Lenovo IdeaPad 3", Price: 6500000, Stock: 10, ImageURL: "https://images.unsplash.com/photo-1519389950473-47ba0277781c", CategoryID: laptop.ID},

		// ===== SMARTPHONE =====
		{Name: "iPhone 15 Pro", Price: 21000000, Stock: 10, ImageURL: "https://images.unsplash.com/photo-1695048133142-1a20484d2569", CategoryID: smartphone.ID},
		{Name: "Samsung S24 Ultra", Price: 22000000, Stock: 10, ImageURL: "https://images.unsplash.com/photo-1610945265064-0e34e5519bbf", CategoryID: smartphone.ID},
		{Name: "Xiaomi 14", Price: 11000000, Stock: 10, ImageURL: "https://images.unsplash.com/photo-1598327105666-5b89351aff97", CategoryID: smartphone.ID},
		{Name: "Oppo Find X6", Price: 15000000, Stock: 10, ImageURL: "https://images.unsplash.com/photo-1605236453806-6ff36851218e", CategoryID: smartphone.ID},
		{Name: "Vivo X100", Price: 14000000, Stock: 10, ImageURL: "https://images.unsplash.com/photo-1580910051074-3eb694886505", CategoryID: smartphone.ID},
		{Name: "Realme GT5", Price: 9000000, Stock: 10, ImageURL: "https://images.unsplash.com/photo-1556656793-08538906a9f8", CategoryID: smartphone.ID},
		{Name: "Samsung A54", Price: 6000000, Stock: 10, ImageURL: "https://images.unsplash.com/photo-1567581935884-3349723552ca", CategoryID: smartphone.ID},
		{Name: "iPhone 13", Price: 11000000, Stock: 10, ImageURL: "https://images.unsplash.com/photo-1632661674596-df8be070a5c5", CategoryID: smartphone.ID},
		{Name: "Redmi Note 13", Price: 3500000, Stock: 10, ImageURL: "https://images.unsplash.com/photo-1592750475338-74b7b21085ab", CategoryID: smartphone.ID},
		{Name: "Infinix Zero Ultra", Price: 5000000, Stock: 10, ImageURL: "https://images.unsplash.com/photo-1601784551446-20c9e07cdbdb", CategoryID: smartphone.ID},

		// ===== ACCESSORIES =====
		{Name: "Logitech MX Master 3", Price: 1500000, Stock: 10, ImageURL: "https://images.unsplash.com/photo-1587825140708-dfaf72ae4b04", CategoryID: accessories.ID},
		{Name: "Razer DeathAdder", Price: 800000, Stock: 10, ImageURL: "https://images.unsplash.com/photo-1615663245857-ac93bb7c39e7", CategoryID: accessories.ID},
		{Name: "Keychron K2", Price: 1400000, Stock: 10, ImageURL: "https://images.unsplash.com/photo-1618384887929-16ec33fab9ef", CategoryID: accessories.ID},
		{Name: "Sony WH1000XM5", Price: 5500000, Stock: 10, ImageURL: "https://images.unsplash.com/photo-1518444028785-8fbcd101ebb9", CategoryID: accessories.ID},
		{Name: "AirPods Pro 2", Price: 3800000, Stock: 10, ImageURL: "https://images.unsplash.com/photo-1600294037681-c80b4cb5b434", CategoryID: accessories.ID},
		{Name: "Samsung SSD T7", Price: 1900000, Stock: 10, ImageURL: "https://images.unsplash.com/photo-1591488320449-011701bb6704", CategoryID: accessories.ID},
		{Name: "Sandisk Flashdisk 128GB", Price: 150000, Stock: 10, ImageURL: "https://images.unsplash.com/photo-1587033411391-5d9e51cce126", CategoryID: accessories.ID},
		{Name: "Anker Powerbank", Price: 600000, Stock: 10, ImageURL: "https://images.unsplash.com/photo-1609091839311-d5365f9ff1c5", CategoryID: accessories.ID},
		{Name: "UGREEN USB Hub", Price: 300000, Stock: 10, ImageURL: "https://images.unsplash.com/photo-1625842268584-8f3296236761", CategoryID: accessories.ID},
		{Name: "HyperX Cloud II", Price: 1200000, Stock: 10, ImageURL: "https://images.unsplash.com/photo-1585298723682-7115561c51b7", CategoryID: accessories.ID},

		// ===== MAKANAN =====
		{Name: "Indomie Goreng", Price: 3500, Stock: 100, ImageURL: "https://images.unsplash.com/photo-1604908176997-125f25cc6f3d", CategoryID: makanan.ID},
		{Name: "Mie Sedaap", Price: 3200, Stock: 100, ImageURL: "https://images.unsplash.com/photo-1612929633738-8fe44f7ec841", CategoryID: makanan.ID},
		{Name: "Beras Ramos 5kg", Price: 75000, Stock: 50, ImageURL: "https://images.unsplash.com/photo-1586201375761-83865001e31c", CategoryID: makanan.ID},
		{Name: "Chitato", Price: 10000, Stock: 80, ImageURL: "https://images.unsplash.com/photo-1566478989037-eec170784d0b", CategoryID: makanan.ID},
		{Name: "SilverQueen", Price: 15000, Stock: 80, ImageURL: "https://images.unsplash.com/photo-1582176604856-e824b4736522", CategoryID: makanan.ID},
		{Name: "Tango Wafer", Price: 12000, Stock: 80, ImageURL: "https://images.unsplash.com/photo-1558961363-fa8fdf82db35", CategoryID: makanan.ID},
		{Name: "Roma Biscuit", Price: 8000, Stock: 80, ImageURL: "https://images.unsplash.com/photo-1509440159596-0249088772ff", CategoryID: makanan.ID},
		{Name: "Sarden ABC", Price: 12000, Stock: 80, ImageURL: "https://images.unsplash.com/photo-1596040033229-a9821ebd058d", CategoryID: makanan.ID},
		{Name: "Kornet Pronas", Price: 25000, Stock: 80, ImageURL: "https://images.unsplash.com/photo-1604908554027-5c7c3f9d6e2f", CategoryID: makanan.ID},
		{Name: "Nugget Fiesta", Price: 45000, Stock: 80, ImageURL: "https://images.unsplash.com/photo-1606755962773-0c6d1e5a0b6e", CategoryID: makanan.ID},

		// ===== MINUMAN =====
		{Name: "Aqua 600ml", Price: 4000, Stock: 100, ImageURL: "https://images.unsplash.com/photo-1564415315949-7a0c4c73e6c4", CategoryID: minuman.ID},
		{Name: "Teh Botol Sosro", Price: 5000, Stock: 100, ImageURL: "https://images.unsplash.com/photo-1558640476-437a2b9438a2", CategoryID: minuman.ID},
		{Name: "Coca Cola", Price: 6000, Stock: 100, ImageURL: "https://images.unsplash.com/photo-1622484212850-eb596d769edc", CategoryID: minuman.ID},
		{Name: "Sprite", Price: 6000, Stock: 100, ImageURL: "https://images.unsplash.com/photo-1624517452488-04869289c4ca", CategoryID: minuman.ID},
		{Name: "Fanta", Price: 6000, Stock: 100, ImageURL: "https://images.unsplash.com/photo-1581006852262-e4307cf6283a", CategoryID: minuman.ID},
		{Name: "Pocari Sweat", Price: 7000, Stock: 100, ImageURL: "https://images.unsplash.com/photo-1590080877777-95a7f3cddbb1", CategoryID: minuman.ID},
		{Name: "Ultra Milk", Price: 7000, Stock: 100, ImageURL: "https://images.unsplash.com/photo-1582719471384-894fbb16e074", CategoryID: minuman.ID},
		{Name: "Good Day Coffee", Price: 4000, Stock: 100, ImageURL: "https://images.unsplash.com/photo-1509042239860-f550ce710b93", CategoryID: minuman.ID},
		{Name: "Yakult", Price: 9000, Stock: 100, ImageURL: "https://images.unsplash.com/photo-1615486363973-1d8c2c5be147", CategoryID: minuman.ID},
		{Name: "Floridina", Price: 3500, Stock: 100, ImageURL: "https://images.unsplash.com/photo-1622597467836-f3285f2131b5", CategoryID: minuman.ID},

		// ===== PERTANIAN =====
		{Name: "Bibit Padi", Price: 50000, Stock: 50, ImageURL: "https://images.unsplash.com/photo-1592982537447-7440770cbfc9", CategoryID: pertanian.ID},
		{Name: "Bibit Jagung", Price: 40000, Stock: 50, ImageURL: "https://images.unsplash.com/photo-1500595046743-cd271d694d30", CategoryID: pertanian.ID},
		{Name: "Pupuk Urea", Price: 120000, Stock: 50, ImageURL: "https://images.unsplash.com/photo-1625246333195-78d9c38ad449", CategoryID: pertanian.ID},
		{Name: "Pupuk Kompos", Price: 60000, Stock: 50, ImageURL: "https://images.unsplash.com/photo-1589923188900-85dae523342b", CategoryID: pertanian.ID},
		{Name: "Cangkul", Price: 80000, Stock: 50, ImageURL: "https://images.unsplash.com/photo-1598514982901-9d6c3e4f6a45", CategoryID: pertanian.ID},
		{Name: "Sekop", Price: 70000, Stock: 50, ImageURL: "https://images.unsplash.com/photo-1574263867128-a3d5c1b5a8e5", CategoryID: pertanian.ID},
		{Name: "Sprayer", Price: 150000, Stock: 50, ImageURL: "https://images.unsplash.com/photo-1628352081506-83c43123ed6d", CategoryID: pertanian.ID},
		{Name: "Polybag", Price: 20000, Stock: 50, ImageURL: "https://images.unsplash.com/photo-1585320806297-9794b3e4eeae", CategoryID: pertanian.ID},
		{Name: "Bibit Cabai", Price: 25000, Stock: 50, ImageURL: "https://images.unsplash.com/photo-1592928301664-5b9c7a509a68", CategoryID: pertanian.ID},
		{Name: "Bibit Tomat", Price: 25000, Stock: 50, ImageURL: "https://images.unsplash.com/photo-1592921870789-04563d55041c", CategoryID: pertanian.ID},

		// ===== MAINAN =====
		{Name: "Lego Classic", Price: 350000, Stock: 30, ImageURL: "https://images.unsplash.com/photo-1587654780291-39c9404d746b", CategoryID: mainan.ID},
		{Name: "Hot Wheels", Price: 50000, Stock: 30, ImageURL: "https://images.unsplash.com/photo-1583267746897-2cf415887172", CategoryID: mainan.ID},
		{Name: "Rubik 3x3", Price: 40000, Stock: 30, ImageURL: "https://images.unsplash.com/photo-1586953208448-b95a79798f07", CategoryID: mainan.ID},
		{Name: "Boneka Teddy", Price: 80000, Stock: 30, ImageURL: "https://images.unsplash.com/photo-1582721478779-0ae163c05a60", CategoryID: mainan.ID},
		{Name: "RC Car", Price: 250000, Stock: 30, ImageURL: "https://images.unsplash.com/photo-1608889175123-8ee362201f81", CategoryID: mainan.ID},
		{Name: "Puzzle", Price: 60000, Stock: 30, ImageURL: "https://images.unsplash.com/photo-1604881987924-1a4b60e0c7c5", CategoryID: mainan.ID},
		{Name: "UNO Card", Price: 30000, Stock: 30, ImageURL: "https://images.unsplash.com/photo-1610890716171-6b1bb98ffd09", CategoryID: mainan.ID},
		{Name: "Monopoly", Price: 200000, Stock: 30, ImageURL: "https://images.unsplash.com/photo-1610890716171-6b1bb98ffd09", CategoryID: mainan.ID},
		{Name: "Action Figure Naruto", Price: 150000, Stock: 30, ImageURL: "https://images.unsplash.com/photo-1611605698335-8b1569810432", CategoryID: mainan.ID},
		{Name: "Drone Mini", Price: 500000, Stock: 30, ImageURL: "https://images.unsplash.com/photo-1508614589041-895b88991e3e", CategoryID: mainan.ID},

		// ===== LAINNYA =====
		{Name: "Kursi Plastik", Price: 50000, Stock: 30, ImageURL: "https://images.unsplash.com/photo-1582582494700-5c0b9d9c4f4c", CategoryID: lainnya.ID},
		{Name: "Meja Lipat", Price: 150000, Stock: 30, ImageURL: "https://images.unsplash.com/photo-1555041469-a586c61ea9bc", CategoryID: lainnya.ID},
		{Name: "Lampu LED", Price: 40000, Stock: 30, ImageURL: "https://images.unsplash.com/photo-1513506003901-1e6a229e2d15", CategoryID: lainnya.ID},
		{Name: "Kipas Angin", Price: 200000, Stock: 30, ImageURL: "https://images.unsplash.com/photo-1582719478250-c89cae4dc85b", CategoryID: lainnya.ID},
		{Name: "Jam Dinding", Price: 70000, Stock: 30, ImageURL: "https://images.unsplash.com/photo-1501139083538-0139583c060f", CategoryID: lainnya.ID},
		{Name: "Botol Minum", Price: 60000, Stock: 30, ImageURL: "https://images.unsplash.com/photo-1526401485004-2fda9f6c0f1c", CategoryID: lainnya.ID},
		{Name: "Payung", Price: 50000, Stock: 30, ImageURL: "https://images.unsplash.com/photo-1500375592092-40eb2168fd21", CategoryID: lainnya.ID},
		{Name: "Tas Ransel", Price: 200000, Stock: 30, ImageURL: "https://images.unsplash.com/photo-1509762774605-f07235a08f1f", CategoryID: lainnya.ID},
		{Name: "Karpet", Price: 150000, Stock: 30, ImageURL: "https://images.unsplash.com/photo-1582582621959-48d27397dc69", CategoryID: lainnya.ID},
		{Name: "Bantal", Price: 80000, Stock: 30, ImageURL: "https://images.unsplash.com/photo-1584100936595-c0654b55a2e2", CategoryID: lainnya.ID},
	}

	// ================= INSERT =================

	for _, p := range products {

		if p.CategoryID == 0 {
			log.Println("skip product, category not found:", p.Name)
			continue
		}

		var existing models.Product
		slug := helper.GenerateSlug(p.Name)

		if err := db.Where("slug = ?", slug).First(&existing).Error; err == nil {
			continue
		}

		p.Slug = slug
		p.ImagePublicID = "seed"

		if err := db.Create(&p).Error; err != nil {
			return err
		}
	}

	return nil
}

func SeedProducts(db *gorm.DB) error {

	defaultImg := "https://picsum.photos/seed/picsum/400/400"

	var laptop, smartphone, accessories, makanan, minuman, pertanian, mainan, lainnya models.Category

	db.Where("name = ?", "Laptop").First(&laptop)
	db.Where("name = ?", "Smartphone").First(&smartphone)
	db.Where("name = ?", "Accessories").First(&accessories)
	db.Where("name = ?", "Makanan").First(&makanan)
	db.Where("name = ?", "Minuman").First(&minuman)
	db.Where("name = ?", "Pertanian").First(&pertanian)
	db.Where("name = ?", "Mainan").First(&mainan)
	db.Where("name = ?", "Lainnya").First(&lainnya)

	products := []models.Product{

		// ================= LAPTOP =================

		{Name: "MacBook Air M2", Price: 18500000, Stock: 10, ImageURL: defaultImg, CategoryID: laptop.ID},
		{Name: "MacBook Pro M3", Price: 32000000, Stock: 10, ImageURL: defaultImg, CategoryID: laptop.ID},
		{Name: "ASUS ROG Zephyrus G14", Price: 25000000, Stock: 10, ImageURL: defaultImg, CategoryID: laptop.ID},
		{Name: "Lenovo Legion 5", Price: 21000000, Stock: 10, ImageURL: defaultImg, CategoryID: laptop.ID},
		{Name: "Acer Aspire 5", Price: 7500000, Stock: 10, ImageURL: defaultImg, CategoryID: laptop.ID},
		{Name: "HP Pavilion 14", Price: 9000000, Stock: 10, ImageURL: defaultImg, CategoryID: laptop.ID},
		{Name: "Dell XPS 13", Price: 24000000, Stock: 10, ImageURL: defaultImg, CategoryID: laptop.ID},
		{Name: "MSI Katana GF66", Price: 17000000, Stock: 10, ImageURL: defaultImg, CategoryID: laptop.ID},
		{Name: "ASUS Vivobook 15", Price: 8000000, Stock: 10, ImageURL: defaultImg, CategoryID: laptop.ID},
		{Name: "Lenovo IdeaPad 3", Price: 6500000, Stock: 10, ImageURL: defaultImg, CategoryID: laptop.ID},

		// ================= SMARTPHONE =================

		{Name: "iPhone 15 Pro", Price: 21000000, Stock: 10, ImageURL: defaultImg, CategoryID: smartphone.ID},
		{Name: "Samsung S24 Ultra", Price: 22000000, Stock: 10, ImageURL: defaultImg, CategoryID: smartphone.ID},
		{Name: "Xiaomi 14", Price: 11000000, Stock: 10, ImageURL: defaultImg, CategoryID: smartphone.ID},
		{Name: "Oppo Find X6", Price: 15000000, Stock: 10, ImageURL: defaultImg, CategoryID: smartphone.ID},
		{Name: "Vivo X100", Price: 14000000, Stock: 10, ImageURL: defaultImg, CategoryID: smartphone.ID},
		{Name: "Realme GT5", Price: 9000000, Stock: 10, ImageURL: defaultImg, CategoryID: smartphone.ID},
		{Name: "Samsung A54", Price: 6000000, Stock: 10, ImageURL: defaultImg, CategoryID: smartphone.ID},
		{Name: "iPhone 13", Price: 11000000, Stock: 10, ImageURL: defaultImg, CategoryID: smartphone.ID},
		{Name: "Redmi Note 13", Price: 3500000, Stock: 10, ImageURL: defaultImg, CategoryID: smartphone.ID},
		{Name: "Infinix Zero Ultra", Price: 5000000, Stock: 10, ImageURL: defaultImg, CategoryID: smartphone.ID},

		// ================= ACCESSORIES =================

		{Name: "Logitech MX Master 3", Price: 1500000, Stock: 10, ImageURL: defaultImg, CategoryID: accessories.ID},
		{Name: "Razer DeathAdder", Price: 800000, Stock: 10, ImageURL: defaultImg, CategoryID: accessories.ID},
		{Name: "Keychron K2", Price: 1400000, Stock: 10, ImageURL: defaultImg, CategoryID: accessories.ID},
		{Name: "Sony WH1000XM5", Price: 5500000, Stock: 10, ImageURL: defaultImg, CategoryID: accessories.ID},
		{Name: "AirPods Pro 2", Price: 3800000, Stock: 10, ImageURL: defaultImg, CategoryID: accessories.ID},
		{Name: "Samsung SSD T7", Price: 1900000, Stock: 10, ImageURL: defaultImg, CategoryID: accessories.ID},
		{Name: "Sandisk Flashdisk 128GB", Price: 150000, Stock: 10, ImageURL: defaultImg, CategoryID: accessories.ID},
		{Name: "Anker Powerbank", Price: 600000, Stock: 10, ImageURL: defaultImg, CategoryID: accessories.ID},
		{Name: "UGREEN USB Hub", Price: 300000, Stock: 10, ImageURL: defaultImg, CategoryID: accessories.ID},
		{Name: "HyperX Cloud II", Price: 1200000, Stock: 10, ImageURL: defaultImg, CategoryID: accessories.ID},

		// ================= MAKANAN =================

		{Name: "Indomie Goreng", Price: 3500, Stock: 100, ImageURL: defaultImg, CategoryID: makanan.ID},
		{Name: "Mie Sedaap", Price: 3200, Stock: 100, ImageURL: defaultImg, CategoryID: makanan.ID},
		{Name: "Beras Ramos 5kg", Price: 75000, Stock: 50, ImageURL: defaultImg, CategoryID: makanan.ID},
		{Name: "Chitato", Price: 10000, Stock: 80, ImageURL: defaultImg, CategoryID: makanan.ID},
		{Name: "SilverQueen", Price: 15000, Stock: 80, ImageURL: defaultImg, CategoryID: makanan.ID},
		{Name: "Tango Wafer", Price: 12000, Stock: 80, ImageURL: defaultImg, CategoryID: makanan.ID},
		{Name: "Roma Biscuit", Price: 8000, Stock: 80, ImageURL: defaultImg, CategoryID: makanan.ID},
		{Name: "Sarden ABC", Price: 12000, Stock: 80, ImageURL: defaultImg, CategoryID: makanan.ID},
		{Name: "Kornet Pronas", Price: 25000, Stock: 80, ImageURL: defaultImg, CategoryID: makanan.ID},
		{Name: "Nugget Fiesta", Price: 45000, Stock: 80, ImageURL: defaultImg, CategoryID: makanan.ID},

		// ================= MINUMAN =================

		{Name: "Aqua 600ml", Price: 4000, Stock: 100, ImageURL: defaultImg, CategoryID: minuman.ID},
		{Name: "Teh Botol Sosro", Price: 5000, Stock: 100, ImageURL: defaultImg, CategoryID: minuman.ID},
		{Name: "Coca Cola", Price: 6000, Stock: 100, ImageURL: defaultImg, CategoryID: minuman.ID},
		{Name: "Sprite", Price: 6000, Stock: 100, ImageURL: defaultImg, CategoryID: minuman.ID},
		{Name: "Fanta", Price: 6000, Stock: 100, ImageURL: defaultImg, CategoryID: minuman.ID},
		{Name: "Pocari Sweat", Price: 7000, Stock: 100, ImageURL: defaultImg, CategoryID: minuman.ID},
		{Name: "Ultra Milk", Price: 7000, Stock: 100, ImageURL: defaultImg, CategoryID: minuman.ID},
		{Name: "Good Day Coffee", Price: 4000, Stock: 100, ImageURL: defaultImg, CategoryID: minuman.ID},
		{Name: "Yakult", Price: 9000, Stock: 100, ImageURL: defaultImg, CategoryID: minuman.ID},
		{Name: "Floridina", Price: 3500, Stock: 100, ImageURL: defaultImg, CategoryID: minuman.ID},

		// ================= PERTANIAN =================

		{Name: "Bibit Padi", Price: 50000, Stock: 50, ImageURL: defaultImg, CategoryID: pertanian.ID},
		{Name: "Bibit Jagung", Price: 40000, Stock: 50, ImageURL: defaultImg, CategoryID: pertanian.ID},
		{Name: "Pupuk Urea", Price: 120000, Stock: 50, ImageURL: defaultImg, CategoryID: pertanian.ID},
		{Name: "Pupuk Kompos", Price: 60000, Stock: 50, ImageURL: defaultImg, CategoryID: pertanian.ID},
		{Name: "Cangkul", Price: 80000, Stock: 50, ImageURL: defaultImg, CategoryID: pertanian.ID},
		{Name: "Sekop", Price: 70000, Stock: 50, ImageURL: defaultImg, CategoryID: pertanian.ID},
		{Name: "Sprayer", Price: 150000, Stock: 50, ImageURL: defaultImg, CategoryID: pertanian.ID},
		{Name: "Polybag", Price: 20000, Stock: 50, ImageURL: defaultImg, CategoryID: pertanian.ID},
		{Name: "Bibit Cabai", Price: 25000, Stock: 50, ImageURL: defaultImg, CategoryID: pertanian.ID},
		{Name: "Bibit Tomat", Price: 25000, Stock: 50, ImageURL: defaultImg, CategoryID: pertanian.ID},

		// ================= MAINAN =================

		{Name: "Lego Classic", Price: 350000, Stock: 30, ImageURL: defaultImg, CategoryID: mainan.ID},
		{Name: "Hot Wheels", Price: 50000, Stock: 30, ImageURL: defaultImg, CategoryID: mainan.ID},
		{Name: "Rubik 3x3", Price: 40000, Stock: 30, ImageURL: defaultImg, CategoryID: mainan.ID},
		{Name: "Boneka Teddy", Price: 80000, Stock: 30, ImageURL: defaultImg, CategoryID: mainan.ID},
		{Name: "RC Car", Price: 250000, Stock: 30, ImageURL: defaultImg, CategoryID: mainan.ID},
		{Name: "Puzzle", Price: 60000, Stock: 30, ImageURL: defaultImg, CategoryID: mainan.ID},
		{Name: "UNO Card", Price: 30000, Stock: 30, ImageURL: defaultImg, CategoryID: mainan.ID},
		{Name: "Monopoly", Price: 200000, Stock: 30, ImageURL: defaultImg, CategoryID: mainan.ID},
		{Name: "Action Figure Naruto", Price: 150000, Stock: 30, ImageURL: defaultImg, CategoryID: mainan.ID},
		{Name: "Drone Mini", Price: 500000, Stock: 30, ImageURL: defaultImg, CategoryID: mainan.ID},

		// ================= LAINNYA =================

		{Name: "Kursi Plastik", Price: 50000, Stock: 30, ImageURL: defaultImg, CategoryID: lainnya.ID},
		{Name: "Meja Lipat", Price: 150000, Stock: 30, ImageURL: defaultImg, CategoryID: lainnya.ID},
		{Name: "Lampu LED", Price: 40000, Stock: 30, ImageURL: defaultImg, CategoryID: lainnya.ID},
		{Name: "Kipas Angin", Price: 200000, Stock: 30, ImageURL: defaultImg, CategoryID: lainnya.ID},
		{Name: "Jam Dinding", Price: 70000, Stock: 30, ImageURL: defaultImg, CategoryID: lainnya.ID},
		{Name: "Botol Minum", Price: 60000, Stock: 30, ImageURL: defaultImg, CategoryID: lainnya.ID},
		{Name: "Payung", Price: 50000, Stock: 30, ImageURL: defaultImg, CategoryID: lainnya.ID},
		{Name: "Tas Ransel", Price: 200000, Stock: 30, ImageURL: defaultImg, CategoryID: lainnya.ID},
		{Name: "Karpet", Price: 150000, Stock: 30, ImageURL: defaultImg, CategoryID: lainnya.ID},
		{Name: "Bantal", Price: 80000, Stock: 30, ImageURL: defaultImg, CategoryID: lainnya.ID},
	}

	for _, p := range products {

		var existing models.Product

		slug := helper.GenerateSlug(p.Name)

		if err := db.Where("slug = ?", slug).First(&existing).Error; err == nil {
			continue
		}

		p.Slug = slug
		p.ImagePublicID = "seed"

		db.Create(&p)

	}

	return nil
}

func SeedProductsFromAssets(db *gorm.DB) error {

	assetDir := "AssetPrivate"

	// ambil semua category dari database
	var categories []models.Category
	if err := db.Find(&categories).Error; err != nil {
		return err
	}

	// buat map: lowercase(name) -> category
	categoryMap := make(map[string]models.Category)
	for _, c := range categories {
		categoryMap[strings.ToLower(c.Name)] = c
	}

	// baca folder utama
	dirs, err := os.ReadDir(assetDir)
	if err != nil {
		return err
	}

	for _, dir := range dirs {

		if !dir.IsDir() {
			continue
		}

		folderName := strings.ToLower(dir.Name())

		category, ok := categoryMap[folderName]
		if !ok {
			fmt.Println("Category tidak ditemukan di DB:", folderName)
			continue
		}

		categoryPath := filepath.Join(assetDir, dir.Name())

		files, err := os.ReadDir(categoryPath)
		if err != nil {
			fmt.Println("Gagal baca folder:", err)
			continue
		}

		for _, file := range files {

			if file.IsDir() {
				continue
			}

			filePath := filepath.Join(categoryPath, file.Name())

			name := strings.TrimSuffix(file.Name(), filepath.Ext(file.Name()))

			var count int64
			db.Model(&models.Product{}).
				Where("name = ?", name).
				Count(&count)

			if count > 0 {
				continue
			}

			uploadResult, err := helper.UploadFile(filePath, "products")
			if err != nil {
				fmt.Println("Upload gagal:", err)
				continue
			}

			product := models.Product{
				Name:          name,
				Slug:          helper.GenerateSlug(name),
				Description:   name,
				ImageURL:      uploadResult.SecureURL,
				ImagePublicID: uploadResult.PublicID,
				Price:         1000000,
				Stock:         100,
				CategoryID:    category.ID,
			}

			if err := db.Create(&product).Error; err != nil {
				fmt.Println("Insert gagal:", err)
				continue
			}

			fmt.Println("Seeded:", name, "| Category:", category.Name)
		}
	}

	return nil
}

func SyncAssetProductsWithDefaultSeed(db *gorm.DB) error {

	// default data dari SeedProducts (tanpa ImageURL & CategoryID)
	defaultProducts := []models.Product{

		// ================= LAPTOP =================
		{Name: "Apple MacBook Air M2", Price: 18500000, Stock: 10,
			Description: "MacBook Air dengan chip Apple M2, RAM 8GB, SSD 256GB, layar Liquid Retina 13.6 inci, baterai hingga 18 jam."},

		{Name: "Apple MacBook Pro M3 14-inch", Price: 32000000, Stock: 10,
			Description: "MacBook Pro dengan chip Apple M3, RAM 16GB, SSD 512GB, layar 14 inci Liquid Retina XDR."},

		{Name: "ASUS ROG Zephyrus G14", Price: 25000000, Stock: 10,
			Description: "Laptop gaming Ryzen 9 dan RTX 4060 dengan layar 165Hz."},

		{Name: "ASUS TUF Gaming F15", Price: 18000000, Stock: 10,
			Description: "Laptop gaming Intel i7 dan RTX series dengan layar 144Hz."},

		{Name: "Lenovo Legion 5 Pro", Price: 21000000, Stock: 10,
			Description: "Ryzen 7 dengan RTX 3060 dan sistem pendingin optimal."},

		{Name: "Lenovo IdeaPad Slim 5", Price: 9500000, Stock: 10,
			Description: "Laptop tipis Ryzen 5, RAM 8GB, SSD 512GB untuk produktivitas."},

		{Name: "Dell XPS 13", Price: 24000000, Stock: 10,
			Description: "Ultrabook premium Intel i7, layar 4K InfinityEdge."},

		{Name: "Dell Inspiron 14", Price: 10000000, Stock: 10,
			Description: "Laptop Intel i5 untuk kerja dan kuliah."},

		{Name: "Acer Aspire 5", Price: 7500000, Stock: 10,
			Description: "Laptop Intel i5 dengan SSD 512GB cocok untuk kebutuhan harian."},

		{Name: "Acer Predator Helios 300", Price: 22000000, Stock: 10,
			Description: "Laptop gaming RTX series dengan performa tinggi."},

		// ================= SMARTPHONE =================
		{Name: "Apple iPhone 15 Pro", Price: 21000000, Stock: 10,
			Description: "Chip A17 Pro, layar 120Hz, kamera 48MP, storage 256GB."},

		{Name: "Apple iPhone 14", Price: 16000000, Stock: 10,
			Description: "Chip A15 Bionic, layar OLED 6.1 inci, kamera dual 12MP."},

		{Name: "Samsung Galaxy S24 Ultra", Price: 22000000, Stock: 10,
			Description: "Snapdragon 8 Gen 3, kamera 200MP, layar AMOLED 120Hz."},

		{Name: "Samsung Galaxy S23 FE", Price: 9000000, Stock: 10,
			Description: "Exynos flagship, layar AMOLED 120Hz, RAM 8GB."},

		{Name: "Xiaomi 14", Price: 11000000, Stock: 10,
			Description: "Snapdragon 8 Gen 3, kamera Leica, layar AMOLED 120Hz."},

		{Name: "Xiaomi Redmi Note 13 Pro", Price: 4500000, Stock: 10,
			Description: "Snapdragon series, kamera 200MP, layar AMOLED 120Hz."},

		{Name: "Vivo X100 Pro", Price: 15000000, Stock: 10,
			Description: "Dimensity 9300, kamera ZEISS flagship."},

		{Name: "Realme GT 5 Pro", Price: 10000000, Stock: 10,
			Description: "Snapdragon 8 Gen series dengan fast charging 150W."},

		{Name: "OPPO Reno11", Price: 7000000, Stock: 10,
			Description: "Smartphone stylish dengan kamera portrait unggulan."},

		{Name: "OPPO Find X6 Pro", Price: 17000000, Stock: 10,
			Description: "Flagship kamera Hasselblad, layar AMOLED premium."},

		// ================= ACCESSORIES =================
		{Name: "Logitech MX Master 3", Price: 1500000, Stock: 10,
			Description: "Mouse premium dengan sensor 4000 DPI, koneksi Bluetooth, baterai tahan lama hingga 70 hari."},

		{Name: "Razer DeathAdder", Price: 800000, Stock: 10,
			Description: "Mouse gaming dengan sensor 20.000 DPI dan desain ergonomis."},

		{Name: "Keychron K2", Price: 1400000, Stock: 10,
			Description: "Mechanical keyboard wireless dengan switch hot-swappable dan RGB backlight."},

		{Name: "Sony WH1000XM5", Price: 5500000, Stock: 10,
			Description: "Headphone dengan Active Noise Cancelling terbaik dan baterai 30 jam."},

		{Name: "AirPods Pro 2", Price: 3800000, Stock: 10,
			Description: "TWS dengan chip H2, ANC, dan spatial audio."},

		{Name: "Samsung SSD T7", Price: 1900000, Stock: 10,
			Description: "External SSD 1TB dengan kecepatan hingga 1050MB/s."},

		{Name: "Sandisk Flashdisk 128GB", Price: 150000, Stock: 10,
			Description: "Flashdisk USB 3.0 kapasitas 128GB untuk penyimpanan cepat dan praktis."},

		{Name: "Anker Powerbank", Price: 600000, Stock: 10,
			Description: "Powerbank 20.000mAh dengan fast charging dan proteksi keamanan."},

		{Name: "UGREEN USB Hub", Price: 300000, Stock: 10,
			Description: "USB Hub multiport dengan HDMI dan USB 3.0."},

		{Name: "HyperX Cloud II", Price: 1200000, Stock: 10,
			Description: "Headset gaming dengan surround sound 7.1 dan mic noise cancelling."},

		// ================= MAKANAN =================
		{Name: "Indomie Goreng", Price: 3500, Stock: 100,
			Description: "Mi instan goreng dengan rasa khas Indonesia, praktis dan lezat."},

		{Name: "Mie Sedaap", Price: 3200, Stock: 100,
			Description: "Mi instan dengan bumbu gurih dan tekstur kenyal."},

		{Name: "Beras Ramos 5kg", Price: 75000, Stock: 50,
			Description: "Beras premium pulen dan wangi untuk kebutuhan rumah tangga."},

		{Name: "Chitato", Price: 10000, Stock: 80,
			Description: "Keripik kentang renyah dengan berbagai varian rasa."},

		{Name: "SilverQueen", Price: 15000, Stock: 80,
			Description: "Coklat premium dengan kacang mete berkualitas."},

		{Name: "Tango Wafer", Price: 12000, Stock: 80,
			Description: "Wafer renyah dengan lapisan krim tebal dan rasa nikmat."},

		{Name: "Roma Biscuit", Price: 8000, Stock: 80,
			Description: "Biskuit renyah cocok untuk camilan keluarga."},

		{Name: "Sarden ABC", Price: 12000, Stock: 80,
			Description: "Sarden kaleng dengan saus tomat lezat dan bergizi."},

		{Name: "Kornet Pronas", Price: 25000, Stock: 80,
			Description: "Kornet sapi berkualitas untuk berbagai olahan makanan."},

		{Name: "Nugget Fiesta", Price: 45000, Stock: 80,
			Description: "Nugget ayam praktis dan lezat untuk keluarga."},

		// ================= MINUMAN =================
		{Name: "Aqua 600ml", Price: 4000, Stock: 100,
			Description: "Air mineral higienis dan menyegarkan."},

		{Name: "Teh Botol Sosro", Price: 5000, Stock: 100,
			Description: "Teh melati asli dengan rasa manis segar."},

		{Name: "Coca Cola", Price: 6000, Stock: 100,
			Description: "Minuman soda berkarbonasi dengan rasa khas."},

		{Name: "Sprite", Price: 6000, Stock: 100,
			Description: "Minuman soda rasa lemon-lime menyegarkan."},

		{Name: "Fanta", Price: 6000, Stock: 100,
			Description: "Minuman soda dengan rasa buah segar."},

		{Name: "Pocari Sweat", Price: 7000, Stock: 100,
			Description: "Minuman isotonik pengganti cairan tubuh."},

		{Name: "Ultra Milk", Price: 7000, Stock: 100,
			Description: "Susu UHT bernutrisi untuk kebutuhan harian."},

		{Name: "Good Day Coffee", Price: 4000, Stock: 100,
			Description: "Minuman kopi instan dengan rasa creamy."},

		{Name: "Yakult", Price: 9000, Stock: 100,
			Description: "Minuman probiotik untuk kesehatan pencernaan."},

		{Name: "Floridina", Price: 5000, Stock: 100,
			Description: "Minuman jeruk dengan bulir asli yang segar."},

		// ================= PERTANIAN =================
		{Name: "Bibit Padi", Price: 50000, Stock: 50,
			Description: "Bibit padi unggul dengan daya tumbuh tinggi."},

		{Name: "Bibit Jagung", Price: 40000, Stock: 50,
			Description: "Bibit jagung berkualitas dengan hasil panen maksimal."},

		{Name: "Pupuk Urea", Price: 120000, Stock: 50,
			Description: "Pupuk nitrogen untuk pertumbuhan optimal tanaman."},

		{Name: "Pupuk Kompos", Price: 60000, Stock: 50,
			Description: "Pupuk organik ramah lingkungan."},

		{Name: "Cangkul", Price: 80000, Stock: 50,
			Description: "Alat pertanian berbahan baja kuat dan tahan lama."},

		{Name: "Sekop", Price: 70000, Stock: 50,
			Description: "Sekop baja ergonomis untuk berbagai kebutuhan."},

		{Name: "Sprayer", Price: 150000, Stock: 50,
			Description: "Sprayer manual untuk penyemprotan pupuk dan pestisida."},

		{Name: "Polybag", Price: 20000, Stock: 50,
			Description: "Polybag kuat untuk media tanam."},

		{Name: "Bibit Cabai", Price: 25000, Stock: 50,
			Description: "Bibit cabai unggul dengan hasil pedas maksimal."},

		{Name: "Bibit Tomat", Price: 25000, Stock: 50,
			Description: "Bibit tomat segar dengan pertumbuhan cepat."},

		// ================= MAINAN =================
		{Name: "Lego Classic", Price: 350000, Stock: 30,
			Description: "Balok Lego kreatif untuk melatih imajinasi anak."},

		{Name: "Hot Wheels", Price: 50000, Stock: 30,
			Description: "Mobil mini koleksi dengan desain menarik."},

		{Name: "Rubik 3x3", Price: 40000, Stock: 30,
			Description: "Puzzle klasik untuk melatih logika dan konsentrasi."},

		{Name: "Boneka Teddy", Price: 80000, Stock: 30,
			Description: "Boneka lembut dan nyaman untuk anak-anak."},

		{Name: "RC Car", Price: 250000, Stock: 30,
			Description: "Mobil remote control dengan baterai rechargeable."},

		{Name: "Puzzle", Price: 60000, Stock: 30,
			Description: "Puzzle edukatif untuk melatih fokus anak."},

		{Name: "UNO Card", Price: 30000, Stock: 30,
			Description: "Permainan kartu seru untuk keluarga."},

		{Name: "Monopoly", Price: 200000, Stock: 30,
			Description: "Permainan papan strategi jual beli properti."},

		{Name: "Action Figure Naruto", Price: 150000, Stock: 30,
			Description: "Figure karakter Naruto dengan detail presisi tinggi."},

		{Name: "Drone Mini", Price: 500000, Stock: 30,
			Description: "Drone mini dengan kamera HD dan kontrol stabil."},

		// ================= LAINNYA =================
		{Name: "Kursi Plastik", Price: 50000, Stock: 30,
			Description: "Kursi plastik kuat dan ringan untuk berbagai kebutuhan."},

		{Name: "Meja Lipat", Price: 150000, Stock: 30,
			Description: "Meja lipat praktis dan mudah disimpan."},

		{Name: "Lampu LED", Price: 40000, Stock: 30,
			Description: "Lampu LED hemat energi dengan cahaya terang."},

		{Name: "Kipas Angin", Price: 200000, Stock: 30,
			Description: "Kipas angin listrik 3 kecepatan dengan konsumsi daya rendah."},

		{Name: "Jam Dinding", Price: 70000, Stock: 30,
			Description: "Jam dinding minimalis dengan mesin presisi."},

		{Name: "Botol Minum", Price: 60000, Stock: 30,
			Description: "Botol minum tahan panas dan dingin, cocok untuk aktivitas harian."},

		{Name: "Payung", Price: 50000, Stock: 30,
			Description: "Payung kuat anti angin dan tahan air."},

		{Name: "Tas Ransel", Price: 200000, Stock: 30,
			Description: "Tas ransel multifungsi dengan bahan tahan air ringan."},

		{Name: "Karpet", Price: 150000, Stock: 30,
			Description: "Karpet lembut dan nyaman untuk ruang tamu atau kamar."},

		{Name: "Bantal", Price: 80000, Stock: 30,
			Description: "Bantal empuk dengan bahan berkualitas untuk kenyamanan tidur."},
	}

	// buat map name -> default product
	defaultMap := make(map[string]models.Product)
	for _, p := range defaultProducts {
		defaultMap[p.Name] = p
	}

	// ambil semua product di DB
	var products []models.Product
	if err := db.Find(&products).Error; err != nil {
		return err
	}

	for _, product := range products {

		defaultData, exists := defaultMap[product.Name]
		if !exists {
			continue
		}

		updates := map[string]interface{}{
			"description": defaultData.Description,
			"price":       defaultData.Price,
			"stock":       defaultData.Stock,
		}

		if err := db.Model(&product).Updates(updates).Error; err != nil {
			fmt.Println("Gagal update:", product.Name)
			continue
		}

		fmt.Println("Updated:", product.Name)
	}

	return nil
}

func SeedCoupons(db *gorm.DB) error {

	now := time.Now()

	exp30 := now.AddDate(0, 0, 30)
	exp60 := now.AddDate(0, 0, 60)
	exp3 := now.AddDate(0, 0, 3)
	exp1 := now.AddDate(0, 0, 1)

	coupons := []models.Coupon{

		// percentage
		{
			Code:            "DISC10",
			Type:            "percentage",
			Value:           10,
			MinimumPurchase: 50000, // Contoh: Min. belanja 50rb
			Quota:           100,
			UsedCount:       0,
			ExpiresAt:       &exp30,
		},
		{
			Code:            "DISC20",
			Type:            "percentage",
			Value:           20,
			MinimumPurchase: 100000,
			Quota:           50,
			UsedCount:       0,
			ExpiresAt:       &exp30,
		},
		{
			Code:            "DISC30",
			Type:            "percentage",
			Value:           30,
			MinimumPurchase: 150000,
			Quota:           30,
			UsedCount:       0,
			ExpiresAt:       &exp30,
		},
		{
			Code:            "WELCOME",
			Type:            "percentage",
			Value:           15,
			MinimumPurchase: 0,
			Quota:           200,
			UsedCount:       0,
			ExpiresAt:       &exp30,
		},
		{
			Code:            "VOCA",
			Type:            "percentage",
			Value:           75,
			MinimumPurchase: 0,
			Quota:           10000,
			UsedCount:       0,
			ExpiresAt:       &exp60,
		},
		{
			Code:            "FLASH50",
			Type:            "percentage",
			Value:           50,
			MinimumPurchase: 200000,
			Quota:           10,
			UsedCount:       0,
			ExpiresAt:       &exp3,
		},

		// fixed
		{
			Code:            "FIXED10K",
			Type:            "fixed",
			Value:           10000,
			MinimumPurchase: 50000,
			Quota:           100,
			UsedCount:       0,
			ExpiresAt:       &exp30,
		},
		{
			Code:            "FIXED25K",
			Type:            "fixed",
			Value:           25000,
			MinimumPurchase: 100000,
			Quota:           100,
			UsedCount:       0,
			ExpiresAt:       &exp30,
		},
		{
			Code:            "FIXED50K",
			Type:            "fixed",
			Value:           50000,
			MinimumPurchase: 250000,
			Quota:           50,
			UsedCount:       0,
			ExpiresAt:       &exp30,
		},
		{
			Code:            "FIXED100K",
			Type:            "fixed",
			Value:           100000,
			MinimumPurchase: 500000,
			Quota:           20,
			UsedCount:       0,
			ExpiresAt:       &exp30,
		},
		{
			Code:            "BIGSALE",
			Type:            "fixed",
			Value:           150000,
			MinimumPurchase: 750000,
			Quota:           10,
			UsedCount:       0,
			ExpiresAt:       &exp3,
		},
		{
			Code:            "LIMITED",
			Type:            "fixed",
			Value:           200000,
			MinimumPurchase: 1000000,
			Quota:           5,
			UsedCount:       0,
			ExpiresAt:       &exp1,
		},
		{
			Code:            "ONGKIR",
			Type:            "fixed",
			Value:           15000,
			MinimumPurchase: 30000,
			Quota:           300,
			UsedCount:       0,
			ExpiresAt:       &exp60,
		},
	}

	return db.
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "code"}},
			DoNothing: true,
		}).
		Create(&coupons).Error
}
