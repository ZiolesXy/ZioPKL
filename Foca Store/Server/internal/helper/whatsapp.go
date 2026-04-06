package helper

import (
	"fmt"
	"net/url"
	"os"

	"golang.org/x/text/language"
	"golang.org/x/text/message"

	"voca-store/internal/domain/models"
)

func formatRupiah(value float64) string {
	p := message.NewPrinter(language.Indonesian)
	return p.Sprintf("%d", int(value))
}

func GenerateCheckoutWhatsappURL(checkout models.Checkout) string {

	adminNumber := os.Getenv("WHATSAPP_ADMIN_NUMBER")
	if adminNumber == "" {
		adminNumber = "6281234567890"
	}

	userName := "Unknown"
	userEmail := "Unknown"
	userPhone := "Unknown"
	if checkout.User != nil {
		userName = checkout.User.Name
		userEmail = checkout.User.Email
		userPhone = checkout.User.TelephoneNumber
	}

	message := fmt.Sprintf(
		"Halo Admin, saya ingin konfirmasi checkout:\n\n"+
			"UID: %s\n"+
			"Status: %s\n"+
			"Nama: %s\n"+
			"Email: %s\n"+
			"No HP: %s\n\n",
		checkout.UID,
		checkout.Status,
		userName,
		userEmail,
		userPhone,
	)

	message += "List Items:\n"
	for _, item := range checkout.Items {
		message += fmt.Sprintf(
			"- %s (Qty: %d) - Rp %s\n",
			item.Product.Name,
			item.Quantity,
			formatRupiah(item.Price),
		)
	}

	message += fmt.Sprintf("\nSubtotal: Rp %s\n", formatRupiah(checkout.Subtotal))
	message += fmt.Sprintf("Discount: Rp %s\n", formatRupiah(checkout.DiscountAmount))
	message += fmt.Sprintf("Total Price: Rp %s\n\n", formatRupiah(checkout.TotalPrice))

	if checkout.Address != nil {
		message += fmt.Sprintf("Alamat Lengkap:\n%s, %s, %s, %s, %s\n",
			checkout.Address.RecipientName,
			checkout.Address.Phone,
			checkout.Address.Province,
			checkout.Address.City,
			checkout.Address.AddressLine,
		)
	} else {
		message += "Alamat Lengkap:\n(Tidak tersedia)\n"
	}

	encodedMessage := url.QueryEscape(message)

	return fmt.Sprintf(
		"https://api.whatsapp.com/send?phone=%s&text=%s",
		adminNumber,
		encodedMessage,
	)
}
