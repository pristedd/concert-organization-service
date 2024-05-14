package main

import (
	"context"
	"github.com/labstack/echo/v4"
	"transaction_service/internal/config"
	transaction2 "transaction_service/internal/transaction"
	transaction "transaction_service/internal/transaction/db"
	"transaction_service/pkg/client/postgresql"
	"transaction_service/pkg/logging"
)

func main() {
	logger := logging.GetLogger()
	logger.Info("create router...")
	router := echo.New()

	cfg := config.GetConfig()
	postgreSQLClient, err := postgresql.NewClient(context.TODO(), 3, cfg.Storage)
	if err != nil {
		logger.Fatalf("%v", err)
	}

	transactionRepository := transaction.NewRepository(postgreSQLClient, logger)

	logger.Info("register transaction handlers...")
	transactionHandler := transaction2.NewHandler(transactionRepository, logger)
	transactionHandler.Register(router)

	logger.Info("start server...")
	addr := cfg.Listen.BindIP + ":" + cfg.Listen.Port
	logger.Fatal(router.Start(addr))
}
