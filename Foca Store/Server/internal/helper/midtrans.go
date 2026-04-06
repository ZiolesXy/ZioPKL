package helper

import (
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"os"
	"voca-store/internal/domain/models"

	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
)

func InitMidtrans() snap.Client {
	var s snap.Client

	serverKey := os.Getenv("MIDTRANS_SERVER_KEY")
	isProd := os.Getenv("MIDTRANS_IS_PRODUCTION") == "true"

	if isProd {
		s.New(serverKey, midtrans.Production)
	} else {
		s.New(serverKey, midtrans.Sandbox)
	}

	return s
}

type MidtransResponse struct {
	Token       string
	RedirectURL string
	OrderID     string
}

func CreateSnapTransaction(
	orderID string,
	amount int64,
	discountAmount int64,
	name, email, phone string,
	address models.Address,
	items []models.CartItem,
) (*MidtransResponse, error) {

	snapClient := InitMidtrans()

	var itemDetails []midtrans.ItemDetails
	var calculatedTotal int64 = 0

	for _, item := range items {
		price := int64(item.Product.Price)
		qty := int32(item.Quantity)

		itemDetails = append(itemDetails, midtrans.ItemDetails{
			ID:    fmt.Sprintf("%d", item.ProductID),
			Name:  item.Product.Name,
			Price: price,
			Qty:   qty,
		})

		calculatedTotal += price * int64(qty)
	}

	if discountAmount > 0 {
		itemDetails = append(itemDetails, midtrans.ItemDetails{
			ID:    "DISCOUNT",
			Name:  "Coupon Discount",
			Price: -discountAmount,
			Qty:   1,
		})

		calculatedTotal -= discountAmount
		if calculatedTotal < 0 || calculatedTotal ==  0 {
			calculatedTotal = 1
		}
	}

	if calculatedTotal != amount {
		fmt.Println("CALCULATED:", calculatedTotal)
		fmt.Println("AMOUNT:", amount)
		fmt.Println("DISCOUNT:", discountAmount)
		return nil, fmt.Errorf("gross_amount mismatch with item_details")
	}

	req := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  orderID,
			GrossAmt: amount,
		},
		Items: &itemDetails,
		CustomerDetail: &midtrans.CustomerDetails{
			FName: name,
			Email: email,
			Phone: phone,
			ShipAddr: &midtrans.CustomerAddress{
				FName:       address.RecipientName,
				Phone:       address.Phone,
				Address:     address.AddressLine,
				City:        address.City,
				Postcode:    address.PostalCode,
				CountryCode: "IDN",
			},
			BillAddr: &midtrans.CustomerAddress{
				FName:       address.RecipientName,
				Phone:       address.Phone,
				Address:     address.AddressLine,
				City:        address.City,
				Postcode:    address.PostalCode,
				CountryCode: "IDN",
			},
		},
	}

	resp, err := snapClient.CreateTransaction(req)
	if err != nil {
		return nil, err
	}

	return &MidtransResponse{
		Token:       resp.Token,
		RedirectURL: resp.RedirectURL,
		OrderID:     orderID,
	}, nil
}

func GenerateMidtransSignature(
	orderID, statusCode, grossAmount, serverKey string,
) string {
	data := orderID + statusCode + grossAmount + serverKey
	hash := sha512.Sum512([]byte(data))
	return hex.EncodeToString(hash[:])
}
