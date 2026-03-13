package handlers

import (
	"fmt"
	"net/http"
	"os"
	"voca-store/helper"
	"voca-store/models"
	"voca-store/response"

	"github.com/gin-gonic/gin"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/coreapi"
	"gorm.io/gorm"
)

func MidtransWebhookManual(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var payload struct {
			OrderID string `json:"order_id"`
		}

		if err := c.ShouldBindJSON(&payload); err != nil {
			response.ErrorResponse(c, http.StatusBadRequest, "invalid body request")
			return
		}

		serverKey := os.Getenv("MIDTRANS_SERVER_KEY")
		coreClient := coreapi.Client{}
		coreClient.New(serverKey, midtrans.Sandbox)

		resp, err := coreClient.CheckTransaction(payload.OrderID)
		if err != nil {
			response.ErrorResponse(c, http.StatusInternalServerError, "failed check transaction")
			return
		}

		tx := db.Begin()
		var checkout models.Checkout
		if err := tx.
			Preload("Items").
			Where("midtrans_order_id = ?", payload.OrderID).
			First(&checkout).Error; err != nil {
			tx.Rollback()
			response.ErrorResponse(c, http.StatusNotFound, "order not found")
			return
		}

		if checkout.PaymentStatus == "paid" ||
			checkout.PaymentStatus == "failed" ||
			checkout.PaymentStatus == "expired" {
			tx.Rollback()
			response.SuccessResponse(c, "already processed", nil)
			return
		}

		switch resp.TransactionStatus {

		case "capture", "settlement":
			checkout.PaymentStatus = "paid"
			checkout.Status = "approved"

		case "expire":
			checkout.PaymentStatus = "expired"
			checkout.Status = "rejected"

			for _, item := range checkout.Items {
				err := tx.Model(&models.Product{}).
					Where("id = ?", item.ProductID).
					Update("stock", gorm.Expr("stock + ?", item.Quantity)).
					Error
				if err != nil {
					tx.Rollback()
					response.ErrorResponse(c, http.StatusInternalServerError, "failed restore stock")
					return
				}
			}

		case "cancel", "deny":
			checkout.PaymentStatus = "failed"
			checkout.Status = "rejected"

			for _, item := range checkout.Items {
				err := tx.Model(&models.Product{}).
					Where("id = ?", item.ProductID).
					Update("stock", gorm.Expr("stock + ?", item.Quantity)).
					Error
				if err != nil {
					tx.Rollback()
					response.ErrorResponse(c, http.StatusInternalServerError, "failed restore stock")
					return
				}
			}
		}

		if err := tx.Save(&checkout).Error; err != nil {
			response.ErrorResponse(c, http.StatusInternalServerError, "failed to update status")
			return
		}

		if err := tx.Commit().Error; err != nil {
			response.ErrorResponse(c, http.StatusInternalServerError, "commit failed")
			return
		}

		response.SuccessResponse(c, "midtrans refreshed", nil)
	}
}

func MidtransWebhook(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("WEBHOOK HIT")
		var notification struct {
			OrderID           string `json:"order_id"`
			TransactionStatus string `json:"transaction_status"`
			FraudStatus       string `json:"fraud_status"`
			SignatureKey      string `json:"signature_key"`
			StatusCode        string `json:"status_code"`
			GrossAmount       string `json:"gross_amount"`
		}

		if err := c.ShouldBindJSON(&notification); err != nil {
			response.ErrorResponse(c, http.StatusBadRequest, "invalid body")
			return
		}

		serverKey := os.Getenv("MIDTRANS_SERVER_KEY")
		expectedSignature := helper.GenerateMidtransSignature(
			notification.OrderID,
			notification.StatusCode,
			notification.GrossAmount,
			serverKey,
		)

		if notification.SignatureKey != expectedSignature {
			response.ErrorResponse(c, http.StatusUnauthorized, "invalid signature")
			return
		}

		tx := db.Begin()

		var checkout models.Checkout
		if err := tx.Preload("Items").Where("midtrans_order_id = ?", notification.OrderID).First(&checkout).Error; err != nil {
			fmt.Println("Order not found:", notification.OrderID)
			response.SuccessResponse(c, "ignored - order not found", nil)
		}

		if checkout.PaymentStatus == "paid" ||
			checkout.PaymentStatus == "failed" ||
			checkout.PaymentStatus == "expired" {
			tx.Rollback()
			response.SuccessResponse(c, "already processed", nil)
			return
		}

		switch notification.TransactionStatus {

		case "capture", "settlement":
			checkout.PaymentStatus = "paid"
			checkout.Status = "approved"

		case "expire":
			checkout.PaymentStatus = "expired"
			checkout.Status = "rejected"

			for _, item := range checkout.Items {
				err := tx.Model(&models.Product{}).
					Where("id = ?", item.ProductID).
					Update("stock", gorm.Expr("stock + ?", item.Quantity)).
					Error
				if err != nil {
					tx.Rollback()
					response.ErrorResponse(c, http.StatusInternalServerError, "failed restore stock")
					return
				}
			}

		case "cancel", "deny":
			checkout.PaymentStatus = "failed"
			checkout.Status = "rejected"

			for _, item := range checkout.Items {
				err := tx.Model(&models.Product{}).
					Where("id = ?", item.ProductID).
					Update("stock", gorm.Expr("stock + ?", item.Quantity)).
					Error
				if err != nil {
					tx.Rollback()
					response.ErrorResponse(c, http.StatusInternalServerError, "failed restore stock")
					return
				}
			}
		}

		if err := tx.Save(&checkout).Error; err != nil {
			response.ErrorResponse(c, http.StatusInternalServerError, "failed to update status")
			return
		}

		if err := tx.Commit().Error; err != nil {
			response.ErrorResponse(c, http.StatusInternalServerError, "commit failed")
			return
		}

		response.SuccessResponse(c, "midtrans refreshed", nil)
	}
}
