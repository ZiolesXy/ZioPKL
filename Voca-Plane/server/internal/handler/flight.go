package handler

import (
	"net/http"
	"strconv"
	"time"
	"voca-plane/internal/service"
	"voca-plane/pkg/response"

	"github.com/gin-gonic/gin"
)

type FlightHandler struct {
	service *service.FlightService
}

func NewFlightHandler(s *service.FlightService) *FlightHandler {
	return &FlightHandler{service: s}
}

func (h *FlightHandler) Search(c *gin.Context) {
	origin := c.Query("origin")
	destination := c.Query("destination")
	date := c.Query("date")
	classType := c.Query("class_type")

	// VALIDASI QUERY
	if origin == "" || destination == "" || date == "" {
		response.Error(c, http.StatusBadRequest, "origin, destination, and date are required")
		return
	}

	// VALIDASI FORMAT DATE
	_, err := time.Parse("2006-01-02", date)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid date format (YYYY-MM-DD)")
		return
	}

	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		response.Error(c, http.StatusBadRequest, "invalid page")
		return
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil || limit < 1 {
		response.Error(c, http.StatusBadRequest, "invalid limit")
		return
	}

	flights, total, err := h.service.SearchFlight(
		c.Request.Context(),
		origin,
		destination,
		date,
		classType,
		page,
		limit,
	)

	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	meta := gin.H{
		"total": total,
		"page":  page,
		"limit": limit,
	}

	response.SuccessWithMeta(c, http.StatusOK, "flights found", flights, meta)
}


func (h *FlightHandler) GetAll(c *gin.Context) {
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

func (h *FlightHandler) GetAllFull(c *gin.Context) {
	flights, err := h.service.GetAllFlightsFull(c.Request.Context())
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "all flights retrieved", flights)
}

func (h *FlightHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	flightID, err := strconv.Atoi(id)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid flight id")
		return
	}

	flight, err := h.service.GetFlightByID(c.Request.Context(), uint(flightID))
	if err != nil {
		response.Error(c, http.StatusNotFound, "flight not found")
		return
	}

	response.Success(c, http.StatusOK, "flight found", flight)
}

// GetSeats returns available seats for a flight, optionally filtered by class_type.
func (h *FlightHandler) GetSeats(c *gin.Context) {
	id := c.Param("id")
	flightID, err := strconv.Atoi(id)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid flight id")
		return
	}

	classType := c.Query("class_type")

	seats, err := h.service.GetAvailableSeats(c.Request.Context(), uint(flightID), classType)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "available seats retrieved", seats)
}