package booking

import (
	"booking_service/internal/handlers"
	"booking_service/pkg/logging"
	"context"
	"github.com/labstack/echo/v4"
	"net/http"
)

type handlerBooking struct {
	logger     *logging.Logger
	repository Repository
}

func NewBookingHandler(repository Repository, logger *logging.Logger) handlers.Handler {
	return &handlerBooking{
		repository: repository,
		logger:     logger,
	}
}

func (h *handlerBooking) Register(router *echo.Echo) {
	router.GET("/booking", h.GetListBooking)
	router.GET("/booking/:id", h.GetBookingById)
	router.POST("/booking", h.CreateBooking)
	router.PUT("/booking/:id", h.UpdateBooking)
	router.DELETE("booking/:id", h.DeleteBooking)
}

func (h *handlerBooking) GetListBooking(c echo.Context) error {
	all, err := h.repository.FindAll(context.TODO())
	if err != nil {
		h.logger.Fatalf("%v", err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	if all == nil {
		return c.JSON(http.StatusNotFound, "No locations found")
	}
	return c.JSON(http.StatusOK, all)
}

func (h *handlerBooking) GetBookingById(c echo.Context) error {
	id := c.Param("id")

	one, err := h.repository.FindOne(context.TODO(), id)
	if err != nil {
		h.logger.Fatalf("%v", err)
		return c.JSON(http.StatusNotFound, err)
	}

	return c.JSON(http.StatusOK, one)
}

func (h *handlerBooking) CreateBooking(c echo.Context) error {
	newBooking := new(Booking)
	if err := c.Bind(newBooking); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	createdBooking, err := h.repository.Create(context.TODO(), newBooking)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusCreated, createdBooking)
}

func (h *handlerBooking) UpdateBooking(c echo.Context) error {
	bookingID := c.Param("id")

	var booking *Booking
	if err := c.Bind(&booking); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	err := h.repository.Update(context.TODO(), booking, bookingID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, "booking is updated")
}

func (h *handlerBooking) DeleteBooking(c echo.Context) error {
	id := c.Param("id")
	err := h.repository.Delete(context.TODO(), id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusNoContent, "booking is deleted")
}
