package address

import (
	"net/http"
	"voca-store/internal/domain/dto/request"
	"voca-store/internal/domain/dto/response"
	"voca-store/internal/service"

	"github.com/gin-gonic/gin"
)

type AddressHandler struct {
	addressService service.AddressService
}

func NewAddressHandler(addressService service.AddressService) *AddressHandler {
	return &AddressHandler{addressService}
}

func (h *AddressHandler) CreateAddress(c *gin.Context) {
	userIDRaw, _ := c.Get("user_id")
	userID := userIDRaw.(uint)

	var req request.CreateAddressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	addressRes, err := h.addressService.Create(userID, req)
	if err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.SuccessResponse(c, "address created", addressRes)
}

func (h *AddressHandler) GetMyAddresses(c *gin.Context) {
	userIDRaw, _ := c.Get("user_id")
	userID := userIDRaw.(uint)

	addresses, err := h.addressService.GetMyAddresses(userID)
	if err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, "failed to fetch addresses")
		return
	}

	response.SuccessResponse(c, "address list fetched", addresses)
}

func (h *AddressHandler) GetAddressByUID(c *gin.Context) {
	userIDRaw, _ := c.Get("user_id")
	userID := userIDRaw.(uint)
	uid := c.Param("uid")

	addressRes, err := h.addressService.GetByUID(uid, userID)
	if err != nil {
		response.ErrorResponse(c, http.StatusNotFound, err.Error())
		return
	}

	response.SuccessResponse(c, "address fetched", addressRes)
}

func (h *AddressHandler) UpdateAddress(c *gin.Context) {
	userIDRaw, _ := c.Get("user_id")
	userID := userIDRaw.(uint)
	uid := c.Param("uid")

	var req request.UpdateAddressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	addressRes, err := h.addressService.Update(uid, userID, req)
	if err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	response.SuccessResponse(c, "address updated", addressRes)
}

func (h *AddressHandler) DeleteAddress(c *gin.Context) {
	userIDRaw, _ := c.Get("user_id")
	userID := userIDRaw.(uint)
	uid := c.Param("uid")

	if err := h.addressService.Delete(uid, userID); err != nil {
		response.ErrorResponse(c, http.StatusNotFound, err.Error())
		return
	}

	response.SuccessResponse(c, "address deleted", nil)
}
