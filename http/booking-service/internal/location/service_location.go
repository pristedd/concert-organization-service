package location

import (
	"booking_service/internal/handlers"
	"booking_service/pkg/logging"
	"context"
	"github.com/labstack/echo/v4"
	"net/http"
)

type handlerLocation struct {
	logger     *logging.Logger
	repository Repository
}

func NewHandler(repository Repository, logger *logging.Logger) handlers.Handler {
	return &handlerLocation{
		repository: repository,
		logger:     logger,
	}
}

func (h *handlerLocation) Register(router *echo.Echo) {
	router.GET("/location", h.GetListLocations)
	router.GET("/location/:id", h.GetLocationById)
	router.POST("/location", h.CreateLocation)
	router.PUT("/location/:id", h.UpdateLocation)
	router.DELETE("location/:id", h.DeleteLocation)
}

func (h *handlerLocation) GetListLocations(c echo.Context) error {
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

func (h *handlerLocation) GetLocationById(c echo.Context) error {
	id := c.Param("id")

	one, err := h.repository.FindOne(context.TODO(), id)
	if err != nil {
		h.logger.Fatalf("%v", err)
		return c.JSON(http.StatusNotFound, err)
	}

	return c.JSON(http.StatusOK, one)
}

func (h *handlerLocation) CreateLocation(c echo.Context) error {
	newLoc := new(Location)
	if err := c.Bind(newLoc); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	createdLocation, err := h.repository.Create(context.TODO(), newLoc)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusCreated, createdLocation)
}

func (h *handlerLocation) UpdateLocation(c echo.Context) error {
	locationID := c.Param("id")

	var loc *Location
	if err := c.Bind(&loc); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	err := h.repository.Update(context.TODO(), loc, locationID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, "location is updated")
}

func (h *handlerLocation) DeleteLocation(c echo.Context) error {
	id := c.Param("id")
	err := h.repository.Delete(context.TODO(), id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusNoContent, "location is deleted")
}
