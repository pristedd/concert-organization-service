package main

import (
	"context"
	"event_service/internal/config"
	event2 "event_service/internal/event"
	event "event_service/internal/event/db"
	ticket2 "event_service/internal/ticket"
	ticket "event_service/internal/ticket/db"
	"event_service/pkg/client/postgresql"
	"event_service/pkg/logging"
	"github.com/labstack/echo/v4"
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

	eventRepository := event.NewRepository(postgreSQLClient, logger)
	ticketRepository := ticket.NewRepository(postgreSQLClient, logger)

	logger.Info("register event handlers...")
	eventHandler := event2.NewEventHandler(eventRepository, logger)
	eventHandler.Register(router)

	logger.Info("register ticket handlers...")
	ticketHandler := ticket2.NewTicketHandler(ticketRepository, logger)
	ticketHandler.Register(router)

	logger.Info("start server...")
	addr := cfg.Listen.BindIP + ":" + cfg.Listen.Port
	logger.Fatal(router.Start(addr))
}
