package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"voca-plane/internal/domain/dto/request"
	"voca-plane/internal/service"
	"voca-plane/pkg/response"

	"github.com/gin-gonic/gin"
)

type TransactionHandler struct {
	service *service.TransactionService
}

func NewTransactionHandler(s *service.TransactionService) *TransactionHandler {
	return &TransactionHandler{service: s}
}

func (h *TransactionHandler) Create(c *gin.Context) {
	var req request.CreateTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	tx, err := h.service.CreateTransaction(c.Request.Context(), userID.(uint), req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(c, http.StatusCreated, "transaction created", tx)
}

func (h *TransactionHandler) Pay(c *gin.Context) {
	code := c.Param("code")

	err := h.service.PayTransaction(c.Request.Context(), code)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "payment successful", nil)
}

func (h *TransactionHandler) GetList(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	transaction, total, err := h.service.GetUserTransactions(c.Request.Context(), userID.(uint), page, limit)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	meta := gin.H{"total": total, "page": page, "limit": limit}
	response.SuccessWithMeta(c, http.StatusOK, "transaction retrieved", transaction, meta)
}

func (h *TransactionHandler) GetListAll(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	transaction, err := h.service.GetUserTransactionsAll(c.Request.Context(), userID.(uint))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "transaction retrieved", transaction)
}

func(h *TransactionHandler) GetByCode(c *gin.Context) {
	code := c.Param("code")

	transaction, err := h.service.GetTransactionByCode(c.Request.Context(), code)
	if err != nil {
		response.Error(c, http.StatusNotFound, "transaction not found")
		return
	}

	response.Success(c, http.StatusOK, "transaction found", transaction)
}

func (h *TransactionHandler) Cancel(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	code := c.Param("code")

	err := h.service.CancelTransaction(c.Request.Context(), userID.(uint), code)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "transaction cancelled", nil)
}

func (h *TransactionHandler) MidtransCallback(c *gin.Context) {

	var payload map[string]interface{}

	if err := c.ShouldBindJSON(&payload); err != nil {
		fmt.Println("BIND ERROR:", err)
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	fmt.Println("MIDTRANS CALLBACK RAW:", payload)

	orderID, ok := payload["order_id"].(string)
	if !ok {
		fmt.Println("order_id missing")
		response.Error(c, http.StatusBadRequest, "invalid order_id")
		return
	}

	status, ok := payload["transaction_status"].(string)
	if !ok {
		fmt.Println("transaction_status missing")
		response.Error(c, http.StatusBadRequest, "invalid transaction_status")
		return
	}

	fmt.Println("ORDER:", orderID)
	fmt.Println("STATUS:", status)

	switch status{
	case "settlement", "capture":
		err := h.service.PayTransaction(
			c.Request.Context(),
			orderID,
		)

		if err != nil {
			response.Error(c, http.StatusInternalServerError, err.Error())
			return
		}
	
	case "expire", "cancel", "deny":
		err := h.service.ExpireTransaction(
			c.Request.Context(),
			orderID,
		)

		if err != nil {
			response.Error(c, http.StatusInternalServerError, err.Error())
		}
	}

	response.Success(c, http.StatusOK, "callback received", nil)
}