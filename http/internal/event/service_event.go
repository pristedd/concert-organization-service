package event

import (
	"context"
	"event_service/internal/handlers"
	"event_service/pkg/logging"
	"github.com/labstack/echo/v4"
	"net/http"
)

type handlerEvent struct {
	logger     *logging.Logger
	repository Repository
}

func NewEventHandler(repository Repository, logger *logging.Logger) handlers.Handler {
	return &handlerEvent{
		repository: repository,
		logger:     logger,
	}
}

func (h *handlerEvent) Register(router *echo.Echo) {
	router.GET("/event", h.GetListEvent)
	router.GET("/event/:id", h.GetEventById)
	router.POST("/event", h.CreateEvent)
	router.PUT("/event/:id", h.UpdateEvent)
	router.DELETE("event/:id", h.DeleteEvent)
}

func (h *handlerEvent) GetListEvent(c echo.Context) error {
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

func (h *handlerEvent) GetEventById(c echo.Context) error {
	id := c.Param("id")

	one, err := h.repository.FindOne(context.TODO(), id)
	if err != nil {
		h.logger.Fatalf("%v", err)
		return c.JSON(http.StatusNotFound, err)
	}

	return c.JSON(http.StatusOK, one)
}

func (h *handlerEvent) CreateEvent(c echo.Context) error {
	newEvent := new(Event)
	if err := c.Bind(newEvent); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	createdEvent, err := h.repository.Create(context.TODO(), newEvent)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusCreated, createdEvent)
}

func (h *handlerEvent) UpdateEvent(c echo.Context) error {
	eventID := c.Param("id")

	var event *Event
	if err := c.Bind(&event); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	err := h.repository.Update(context.TODO(), event, eventID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, "event is updated")
}

func (h *handlerEvent) DeleteEvent(c echo.Context) error {
	id := c.Param("id")
	err := h.repository.Delete(context.TODO(), id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusNoContent, "event is deleted")
}
