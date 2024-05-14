package main

import (
	artist2 "booking_service/internal/artist"
	artist "booking_service/internal/artist/db"
	booking2 "booking_service/internal/booking"
	booking "booking_service/internal/booking/db"
	"booking_service/internal/config"
	location2 "booking_service/internal/location"
	location "booking_service/internal/location/db"
	"booking_service/pkg/client/postgresql"
	"booking_service/pkg/logging"
	"context"
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

	locationRepository := location.NewRepository(postgreSQLClient, logger)
	artistRepository := artist.NewRepository(postgreSQLClient, logger)
	bookingRepository := booking.NewRepository(postgreSQLClient, logger)

	logger.Info("register location handlers...")
	locationHandler := location2.NewHandler(locationRepository, logger)
	locationHandler.Register(router)

	logger.Info("register artist handlers...")
	artistHandler := artist2.NewArtistHandler(artistRepository, logger)
	artistHandler.Register(router)

	logger.Info("register booking handlers...")
	bookingHandler := booking2.NewBookingHandler(bookingRepository, logger)
	bookingHandler.Register(router)

	logger.Info("start server...")
	addr := cfg.Listen.BindIP + ":" + cfg.Listen.Port
	logger.Fatal(router.Start(addr))
}
