package event

import (
	"context"
	"errors"
	"event_service/internal/event"
	"event_service/pkg/client/postgresql"
	"event_service/pkg/logging"
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

func (r repository) Create(ctx context.Context, event *event.Event) (*event.Event, error) {
	q := `INSERT INTO public.event (name)
			VALUES ($1) 
			RETURNING id`

	r.logger.Trace(fmt.Sprintf("SQL Query: %s", formatQuery(q)))
	queryRow := r.client.QueryRow(ctx, q, event.Name)
	if err := queryRow.Scan(&event.ID); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			newErr := fmt.Errorf(fmt.Sprintf("SQL Error: %s, Detail: %s, Where: %s, Code: %s, SQLState: %s", pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState()))
			r.logger.Error(newErr)
			return event, newErr
		}
		return event, err
	}

	return event, nil
}

func (r repository) FindAll(ctx context.Context) (e []event.Event, err error) {
	q := `
		SELECT id, name FROM public.event;
	`

	r.logger.Trace(fmt.Sprintf("SQL Query: %s", formatQuery(q)))

	rows, err := r.client.Query(ctx, q)
	if err != nil {
		return nil, err
	}

	events := make([]event.Event, 0)

	for rows.Next() {
		var e event.Event

		err = rows.Scan(&e.ID, &e.Name)
		if err != nil {
			return nil, err
		}

		events = append(events, e)

		if err = rows.Err(); err != nil {
			return nil, err
		}
	}

	return events, nil
}

func (r repository) FindOne(ctx context.Context, id string) (event.Event, error) {
	q := `
		SELECT id, name FROM public.event WHERE id = $1;
	`

	r.logger.Trace(fmt.Sprintf("SQL Query: %s", formatQuery(q)))

	var e event.Event
	queryRow := r.client.QueryRow(ctx, q, id)
	err := queryRow.Scan(&e.ID, &e.Name)
	if err != nil {
		return event.Event{}, err
	}
	return e, nil
}

func (r repository) Update(ctx context.Context, event *event.Event, id string) error {
	q := `
	UPDATE public.event   
	SET name = $1
	Where id = $2;
	`

	r.logger.Trace(fmt.Sprintf("SQL Query: %s", formatQuery(q)))

	_, err := r.client.Exec(ctx, q, event.Name, id)
	if err != nil {
		return err
	}
	return nil
}

func (r repository) Delete(ctx context.Context, id string) error {
	q := `
	DELETE FROM public.event WHERE id=$1
	`

	r.logger.Trace(fmt.Sprintf("SQL Query: %s", formatQuery(q)))

	_, err := r.client.Exec(ctx, q, id)
	if err != nil {
		return err
	}

	return nil
}

func NewRepository(client *pgxpool.Pool, logger *logging.Logger) event.Repository {
	return &repository{
		client: client,
		logger: logger,
	}
}
