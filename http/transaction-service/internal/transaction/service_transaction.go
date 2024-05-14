package transaction

import (
	"context"
	"github.com/labstack/echo/v4"
	"net/http"
	"transaction_service/internal/handlers"
	"transaction_service/pkg/logging"
)

type handlerTransaction struct {
	logger     *logging.Logger
	repository Repository
}

func NewHandler(repository Repository, logger *logging.Logger) handlers.Handler {
	return &handlerTransaction{
		repository: repository,
		logger:     logger,
	}
}

func (h *handlerTransaction) Register(router *echo.Echo) {
	router.GET("/transaction", h.GetListTransactions)
	router.GET("/transaction/:id", h.GetTransactionById)
	router.POST("/transaction", h.CreateTransaction)
}

func (h *handlerTransaction) GetListTransactions(c echo.Context) error {
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

func (h *handlerTransaction) GetTransactionById(c echo.Context) error {
	id := c.Param("id")

	one, err := h.repository.FindOne(context.TODO(), id)
	if err != nil {
		h.logger.Fatalf("%v", err)
		return c.JSON(http.StatusNotFound, err)
	}

	return c.JSON(http.StatusOK, one)
}

func (h *handlerTransaction) CreateTransaction(c echo.Context) error {
	newTransact := new(Transaction)
	if err := c.Bind(newTransact); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	createdTransaction, err := h.repository.Create(context.TODO(), newTransact)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusCreated, createdTransaction)
}
