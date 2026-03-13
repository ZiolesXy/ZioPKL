package handler

import (
	"net/http"
	"strconv"
	"time"
	"voca-plane/internal/domain/dto/request"
	"voca-plane/internal/domain/models"
	"voca-plane/internal/service"
	"voca-plane/pkg/helper"
	"voca-plane/pkg/response"

	"github.com/gin-gonic/gin"
)

type AdminHandler struct {
	service *service.AdminService
}

func NewAdminHandler(s *service.AdminService) *AdminHandler {
	return &AdminHandler{service: s}
}

func (h *AdminHandler) GetDashboard(c *gin.Context) {
	stats, err := h.service.GetDashboardStats(c.Request.Context())
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, http.StatusOK, "dashboard stats retrieved", stats)
}

func (h *AdminHandler) GetUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	sortBy := c.Query("sort_by")
	order := c.Query("order")

	users, total, err := h.service.GetAllUsers(c.Request.Context(), page, limit, sortBy, order)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	meta := gin.H{"total": total, "page": page, "limit": limit}
	response.SuccessWithMeta(c, http.StatusOK, "users retrieved", users, meta)
}

func (h *AdminHandler) UpdateUserRole(c *gin.Context) {
	id := c.Param("id")
	userID, err := strconv.Atoi(id)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid user id")
		return
	}

	var req request.UpdateUserRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	err = h.service.UpdateUserRole(c.Request.Context(), uint(userID), req.Role)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "user role updated", nil)
}

func (h *AdminHandler) DeleteUser(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	err := h.service.DeleteUser(c.Request.Context(), uint(id))
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "user deleted", nil)
}

func (h *AdminHandler) RestoreUser(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	err := h.service.RestoreUser(c.Request.Context(), uint(id))
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "user restored", nil)
}

func (h *AdminHandler) BanUser(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var req request.BanUserRequest

	c.ShouldBindJSON(&req)

	err := h.service.BanUser(c.Request.Context(), uint(id), req.Reason)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "user banned", nil)
}

func (h *AdminHandler) UnbanUser(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	err := h.service.UnbanUser(c.Request.Context(), uint(id))
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "user unbanned", nil)
}

func (h *AdminHandler) GetTransactions(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	sortBy := c.Query("sort_by")
	order := c.Query("order")

	transactions, total, err := h.service.GetAllTransactions(c.Request.Context(), page, limit, sortBy, order)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	meta := gin.H{"total": total, "page": page, "limit": limit}
	response.SuccessWithMeta(c, http.StatusOK, "transactions retrieved", transactions, meta)
}

func (h *AdminHandler) GetFlights(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	sortBy := c.Query("sort_by")
	order := c.Query("order")

	flights, total, err := h.service.GetAllFlights(c.Request.Context(), page, limit, sortBy, order)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	meta := gin.H{"total": total, "page": page, "limit": limit}
	response.SuccessWithMeta(c, http.StatusOK, "flights retrieved", flights, meta)
}

func (h *AdminHandler) CreateFlight(c *gin.Context) {
	var req request.CreateFlightRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	departureTime, _ := time.Parse(time.RFC3339, req.DepartureTime)
	arrivalTime, _ := time.Parse(time.RFC3339, req.ArrivalTime)

	flight := &models.Flight{
		AirlineID:     req.AirlineID,
		OriginID:      req.OriginID,
		DestinationID: req.DestinationID,
		DepartureTime: departureTime,
		ArrivalTime:   arrivalTime,
		FlightNumber:  req.FlightNumber,
		TotalSeats:    req.TotalSeats,
		TotalRows:     req.TotalRows,
		TotalColumns:  req.TotalColumns,
	}

	flightResponse, err := h.service.CreateFlight(c.Request.Context(), flight, req.ClassCount, req.ClassPrices)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(c, http.StatusCreated, "flight created", flightResponse)
}

func (h *AdminHandler) UpdateFlight(c *gin.Context) {
	id := c.Param("id")
	flightID, err := strconv.Atoi(id)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid flight id")
		return
	}

	var req request.UpdateFlightRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	flight, err := h.service.GetFlightModelByID(c.Request.Context(), uint(flightID))
	if err != nil {
		response.Error(c, http.StatusNotFound, "flight not found")
		return
	}

	if req.AirlineID != nil {
		flight.AirlineID = *req.AirlineID
	}
	if req.OriginID != nil {
		flight.OriginID = *req.OriginID
	}
	if req.DestinationID != nil {
		flight.DestinationID = *req.DestinationID
	}
	if req.DepartureTime != nil {
		departureTime, _ := time.Parse(time.RFC3339, *req.DepartureTime)
		flight.DepartureTime = departureTime
	}
	if req.ArrivalTime != nil {
		arrivalTime, _ := time.Parse(time.RFC3339, *req.ArrivalTime)
		flight.ArrivalTime = arrivalTime
	}
	if req.FlightNumber != nil {
		flight.FlightNumber = *req.FlightNumber
	}
	if req.TotalSeats != nil {
		flight.TotalSeats = *req.TotalSeats
	}
	if req.TotalRows != nil {
		flight.TotalRows = *req.TotalRows
	}
	if req.TotalColumns != nil {
		flight.TotalColumns = *req.TotalColumns
	}

	flightResponse, err := h.service.UpdateFlight(c.Request.Context(), flight)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "flight updated", flightResponse)
}

func (h *AdminHandler) DeleteFlight(c *gin.Context) {
	id := c.Param("id")
	flightID, err := strconv.Atoi(id)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid flight id")
		return
	}

	err = h.service.DeleteFlight(c.Request.Context(), uint(flightID))
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "flight deleted", nil)
}

func (h *AdminHandler) GetAirlines(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	sortBy := c.Query("sort_by")
	order := c.Query("order")

	airlines, total, err := h.service.GetAllAirlines(c.Request.Context(), page, limit, sortBy, order)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	meta := gin.H{"total": total, "page": page, "limit": limit}
	response.SuccessWithMeta(c, http.StatusOK, "airlines retrieved", airlines, meta)
}

func (h *AdminHandler) CreateAirline(c *gin.Context) {
	var req request.CreateAirlineRequest

	if err := c.ShouldBind(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	var logoURL string
	var publicID string

	file, err := c.FormFile("logo")

	if err == nil {
		url, pid, err := helper.UploadImage(file)
		if err != nil {
			response.Error(c, http.StatusInternalServerError, err.Error())
			return
		}

		logoURL = url
		publicID = pid
	} else {
		logoURL = req.LogoURL
	}

	airline := &models.Airline{
		Name:         req.Name,
		Code:         req.Code,
		LogoURL:      logoURL,
		LogoPublicID: publicID,
	}

	res, err := h.service.CreateAirline(c.Request.Context(), airline)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(c, http.StatusCreated, "airline created", res)
}

func (h *AdminHandler) UpdateAirline(c *gin.Context) {
	id := c.Param("id")
	airlineID, err := strconv.Atoi(id)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid airline id")
		return
	}

	var req request.UpdateAirlineRequest
	if err := c.ShouldBind(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	airline, err := h.service.GetAirlineByID(c.Request.Context(), uint(airlineID))
	if err != nil {
		response.Error(c, http.StatusNotFound, "airline not found")
		return
	}
	if req.Name != nil {
		airline.Name = *req.Name
	}
	if req.Code != nil {
		airline.Code = *req.Code
	}
	if req.LogoURL != nil {
		airline.LogoURL = *req.LogoURL
	}

	file, err := c.FormFile("logo")
	if err == nil {

		url, pid, err := helper.UploadImage(file)
		if err != nil {
			response.Error(c, http.StatusInternalServerError, err.Error())
			return
		}

		if airline.LogoPublicID != "" {
			helper.DeleteImage(airline.LogoPublicID)
		}

		airline.LogoURL = url
		airline.LogoPublicID = pid
	}

	res, err := h.service.UpdateAirline(c.Request.Context(), airline)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "airline updated", res)
}

func (h *AdminHandler) DeleteAirline(c *gin.Context) {
	id := c.Param("id")
	airlineID, err := strconv.Atoi(id)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid airline id")
		return
	}

	err = h.service.DeleteAirline(c.Request.Context(), uint(airlineID))
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "airline deleted", nil)
}

func (h *AdminHandler) GetAirports(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	sortBy := c.Query("sort_by")
	order := c.Query("order")

	airports, total, err := h.service.GetAllAirports(c.Request.Context(), page, limit, sortBy, order)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	meta := gin.H{"total": total, "page": page, "limit": limit}
	response.SuccessWithMeta(c, http.StatusOK, "airports retrieved", airports, meta)
}

func (h *AdminHandler) CreateAirport(c *gin.Context) {
	var req request.CreateAirportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	airport := &models.Airport{
		Code: req.Code,
		Name: req.Name,
		City: req.City,
	}

	res, err := h.service.CreateAirport(c.Request.Context(), airport)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(c, http.StatusCreated, "airport created", res)
}

func (h *AdminHandler) UpdateAirport(c *gin.Context) {
	id := c.Param("id")
	airportID, err := strconv.Atoi(id)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid airport id")
		return
	}

	var req request.UpdateAirportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	airport, err := h.service.GetAirportByID(c.Request.Context(), uint(airportID))
	if err != nil {
		response.Error(c, http.StatusNotFound, "airport not found")
		return
	}
	if req.Code != nil {
		airport.Code = *req.Code
	}
	if req.Name != nil {
		airport.Name = *req.Name
	}
	if req.City != nil {
		airport.City = *req.City
	}

	res, err := h.service.UpdateAirport(c.Request.Context(), airport)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "airport updated", res)
}

func (h *AdminHandler) DeleteAirport(c *gin.Context) {
	id := c.Param("id")
	airportID, err := strconv.Atoi(id)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid airport id")
		return
	}

	err = h.service.DeleteAirport(c.Request.Context(), uint(airportID))
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "airport deleted", nil)
}

func (h *AdminHandler) GetPromos(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	sortBy := c.Query("sort_by")
	order := c.Query("order")

	promos, total, err := h.service.GetAllPromos(c.Request.Context(), page, limit, sortBy, order)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	meta := gin.H{"total": total, "page": page, "limit": limit}
	response.SuccessWithMeta(c, http.StatusOK, "promos retrieved", promos, meta)
}

func (h *AdminHandler) CreatePromo(c *gin.Context) {
	var req request.CreatePromoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	promo := &models.PromoCode{
		Code:     req.Code,
		Discount: req.Discount,
		IsActive: req.IsActive,
	}

	res, err := h.service.CreatePromo(c.Request.Context(), promo)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(c, http.StatusCreated, "promo created", res)
}

func (h *AdminHandler) UpdatePromo(c *gin.Context) {
	id := c.Param("id")
	promoID, err := strconv.Atoi(id)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid promo id")
		return
	}

	var req request.UpdatePromoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	promo, err := h.service.GetPromoByID(c.Request.Context(), uint(promoID))
	if err != nil {
		response.Error(c, http.StatusNotFound, "promo not found")
		return
	}
	if req.Code != nil {
		promo.Code = *req.Code
	}
	if req.Discount != nil {
		promo.Discount = *req.Discount
	}
	if req.IsActive != nil {
		promo.IsActive = *req.IsActive
	}

	res, err := h.service.UpdatePromo(c.Request.Context(), promo)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "promo updated", res)
}

func (h *AdminHandler) DeletePromo(c *gin.Context) {
	id := c.Param("id")
	promoID, err := strconv.Atoi(id)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid promo id")
		return
	}

	err = h.service.DeletePromo(c.Request.Context(), uint(promoID))
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "promo deleted", nil)
}
