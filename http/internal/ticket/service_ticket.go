package ticket

import (
	"context"
	"event_service/internal/handlers"
	"event_service/pkg/logging"
	"github.com/labstack/echo/v4"
	"net/http"
)

type handlerTicket struct {
	logger     *logging.Logger
	repository Repository
}

func NewTicketHandler(repository Repository, logger *logging.Logger) handlers.Handler {
	return &handlerTicket{
		repository: repository,
		logger:     logger,
	}
}

func (h *handlerTicket) Register(router *echo.Echo) {
	router.GET("/ticket", h.GetListTicket)
	router.GET("/ticket/:id", h.GetTicketById)
	router.POST("/ticket", h.CreateTicket)
}

func (h *handlerTicket) GetListTicket(c echo.Context) error {
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

func (h *handlerTicket) GetTicketById(c echo.Context) error {
	id := c.Param("id")

	one, err := h.repository.FindOne(context.TODO(), id)
	if err != nil {
		h.logger.Fatalf("%v", err)
		return c.JSON(http.StatusNotFound, err)
	}

	return c.JSON(http.StatusOK, one)
}

func (h *handlerTicket) CreateTicket(c echo.Context) error {
	newTicket := new(Ticket)
	if err := c.Bind(newTicket); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	createdTicket, err := h.repository.Create(context.TODO(), newTicket)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusCreated, createdTicket)
}
