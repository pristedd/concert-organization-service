package booking

import (
	"booking_service/internal/booking"
	"booking_service/pkg/client/postgresql"
	"booking_service/pkg/logging"
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"strings"
)

type repository struct {
	client postgresql.Client
	logger *logging.Logger
}

func formatQuery(q string) string {
	return strings.ReplaceAll(strings.ReplaceAll(q, "\t", ""), "\n", " ")
}

func (r repository) Create(ctx context.Context, booking *booking.Booking) (*booking.Booking, error) {
	q := `INSERT INTO public.booking (location_id, artist_id, date, comment) 
			VALUES ($1, $2, $3, $4) 
			RETURNING id`

	r.logger.Trace(fmt.Sprintf("SQL Query: %s", formatQuery(q)))
	queryRow := r.client.QueryRow(ctx, q, booking.LocationID, booking.ArtistID, booking.Date, booking.Comment)
	if err := queryRow.Scan(&booking.ID); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			newErr := fmt.Errorf(fmt.Sprintf("SQL Error: %s, Detail: %s, Where: %s, Code: %s, SQLState: %s", pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState()))
			r.logger.Error(newErr)
			return booking, newErr
		}
		return booking, err
	}

	return booking, nil
}

func (r repository) FindAll(ctx context.Context) (b []booking.Booking, err error) {
	q := `
		SELECT booking.id, booking.location_id, booking.artist_id, booking.date, booking.comment 
		FROM public.booking
		JOIN public.location ON location.id = booking.location_id
		JOIN public.artist ON artist.id = booking.artist_id;
	`

	r.logger.Trace(fmt.Sprintf("SQL Query: %s", formatQuery(q)))

	rows, err := r.client.Query(ctx, q)
	if err != nil {
		return nil, err
	}

	bookings := make([]booking.Booking, 0)

	for rows.Next() {
		var b booking.Booking

		err = rows.Scan(&b.ID, &b.LocationID, &b.ArtistID, &b.Date, &b.Comment)
		if err != nil {
			return nil, err
		}

		bookings = append(bookings, b)

		if err = rows.Err(); err != nil {
			return nil, err
		}
	}

	return bookings, nil
}

func (r repository) FindOne(ctx context.Context, id string) (booking.Booking, error) {
	q := `
		SELECT booking.id, booking.location_id, booking.artist_id, booking.date, booking.comment 
		FROM public.booking
		JOIN public.location ON location.id = booking.location_id
		JOIN public.artist ON artist.id = booking.artist_id
		WHERE booking.id = $1;
	`

	r.logger.Trace(fmt.Sprintf("SQL Query: %s", formatQuery(q)))

	var b booking.Booking
	queryRow := r.client.QueryRow(ctx, q, id)
	err := queryRow.Scan(&b.ID, &b.LocationID, &b.ArtistID, &b.Date, &b.Comment)
	if err != nil {
		return booking.Booking{}, err
	}
	return b, nil
}

func (r repository) Update(ctx context.Context, booking *booking.Booking, id string) error {
	q := `
	UPDATE public.booking 
	SET location_id = $1, artist_id = $2, date = $3, comment = $4
	Where id = $5;
	`

	r.logger.Trace(fmt.Sprintf("SQL Query: %s", formatQuery(q)))

	_, err := r.client.Exec(ctx, q, booking.LocationID, booking.ArtistID, booking.Date, booking.Comment, id)
	if err != nil {
		return err
	}
	return nil
}

func (r repository) Delete(ctx context.Context, id string) error {
	q := `
	DELETE FROM public.booking WHERE id=$1
	`

	r.logger.Trace(fmt.Sprintf("SQL Query: %s", formatQuery(q)))

	_, err := r.client.Exec(ctx, q, id)
	if err != nil {
		return err
	}

	return nil
}

func NewRepository(client *pgxpool.Pool, logger *logging.Logger) booking.Repository {
	return &repository{
		client: client,
		logger: logger,
	}
}
