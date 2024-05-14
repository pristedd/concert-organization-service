package location

import (
	"booking_service/internal/location"
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

func (r repository) Create(ctx context.Context, location *location.Location) (*location.Location, error) {
	q := `INSERT INTO public.location (name, address, seats_num, comment) 
			VALUES ($1, $2, $3, $4) 
			RETURNING id`

	r.logger.Trace(fmt.Sprintf("SQL Query: %s", formatQuery(q)))
	queryRow := r.client.QueryRow(ctx, q, location.Name, location.Address, location.SeatsNum, location.Comment)
	if err := queryRow.Scan(&location.ID); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			newErr := fmt.Errorf(fmt.Sprintf("SQL Error: %s, Detail: %s, Where: %s, Code: %s, SQLState: %s", pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState()))
			r.logger.Error(newErr)
			return location, newErr
		}
		return location, err
	}

	return location, nil
}

func (r repository) FindAll(ctx context.Context) (l []location.Location, err error) {
	q := `
		SELECT id, name, address, seats_num, comment FROM public.location;
	`

	r.logger.Trace(fmt.Sprintf("SQL Query: %s", formatQuery(q)))

	rows, err := r.client.Query(ctx, q)
	if err != nil {
		return nil, err
	}

	locations := make([]location.Location, 0)

	for rows.Next() {
		var l location.Location

		err = rows.Scan(&l.ID, &l.Name, &l.Address, &l.SeatsNum, &l.Comment)
		if err != nil {
			return nil, err
		}

		locations = append(locations, l)

		if err = rows.Err(); err != nil {
			return nil, err
		}
	}

	return locations, nil
}

func (r repository) FindOne(ctx context.Context, id string) (location.Location, error) {
	q := `
		SELECT id, name, address, seats_num, comment FROM public.location WHERE id = $1;
	`

	r.logger.Trace(fmt.Sprintf("SQL Query: %s", formatQuery(q)))

	var l location.Location
	queryRow := r.client.QueryRow(ctx, q, id)
	err := queryRow.Scan(&l.ID, &l.Name, &l.Address, &l.SeatsNum, &l.Comment)
	if err != nil {
		return location.Location{}, err
	}
	return l, nil
}

func (r repository) Update(ctx context.Context, location *location.Location, id string) error {
	q := `
	UPDATE public.location  
	SET name = $1, address = $2, seats_num = $3, comment = $4
	Where id = $5;
	`

	r.logger.Trace(fmt.Sprintf("SQL Query: %s", formatQuery(q)))

	_, err := r.client.Exec(ctx, q, location.Name, location.Address, location.SeatsNum, location.Comment, id)
	if err != nil {
		return err
	}
	return nil
}

func (r repository) Delete(ctx context.Context, id string) error {
	q := `
	DELETE FROM public.location WHERE id=$1
	`

	r.logger.Trace(fmt.Sprintf("SQL Query: %s", formatQuery(q)))

	_, err := r.client.Exec(ctx, q, id)
	if err != nil {
		return err
	}

	return nil
}

func NewRepository(client *pgxpool.Pool, logger *logging.Logger) location.Repository {
	return &repository{
		client: client,
		logger: logger,
	}
}
