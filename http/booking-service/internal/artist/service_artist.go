package artist

import (
	"booking_service/internal/handlers"
	"booking_service/pkg/logging"
	"context"
	"github.com/labstack/echo/v4"
	"net/http"
)

type handlerArtist struct {
	logger     *logging.Logger
	repository Repository
}

func NewArtistHandler(repository Repository, logger *logging.Logger) handlers.Handler {
	return &handlerArtist{
		repository: repository,
		logger:     logger,
	}
}

func (h *handlerArtist) Register(router *echo.Echo) {
	router.GET("/artist", h.GetListArtists)
	router.GET("/artist/:id", h.GetArtistById)
	router.POST("/artist", h.CreateArtist)
	router.PUT("/artist/:id", h.UpdateArtist)
	router.DELETE("artist/:id", h.DeleteArtist)
}

func (h *handlerArtist) GetListArtists(c echo.Context) error {
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

func (h *handlerArtist) GetArtistById(c echo.Context) error {
	id := c.Param("id")

	one, err := h.repository.FindOne(context.TODO(), id)
	if err != nil {
		h.logger.Fatalf("%v", err)
		return c.JSON(http.StatusNotFound, err)
	}

	return c.JSON(http.StatusOK, one)
}

func (h *handlerArtist) CreateArtist(c echo.Context) error {
	newArt := new(Artist)
	if err := c.Bind(newArt); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	createdArtist, err := h.repository.Create(context.TODO(), newArt)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusCreated, createdArtist)
}

func (h *handlerArtist) UpdateArtist(c echo.Context) error {
	ArtistID := c.Param("id")

	var art *Artist
	if err := c.Bind(&art); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	err := h.repository.Update(context.TODO(), art, ArtistID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, "artist is updated")
}

func (h *handlerArtist) DeleteArtist(c echo.Context) error {
	id := c.Param("id")
	err := h.repository.Delete(context.TODO(), id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusNoContent, "artist is deleted")
}
