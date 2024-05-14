package ticket

import (
	"context"
	"errors"
	"event_service/internal/ticket"
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

func (r repository) Create(ctx context.Context, ticket *ticket.Ticket) (*ticket.Ticket, error) {
	q := `INSERT INTO public.ticket (typeid, eventid, price, purchased)
			VALUES ($1, $2, $3, $4) 
			RETURNING id`

	r.logger.Trace(fmt.Sprintf("SQL Query: %s", formatQuery(q)))
	queryRow := r.client.QueryRow(ctx, q, ticket.TypeID, ticket.EventID, ticket.Price, ticket.Purchased)
	if err := queryRow.Scan(&ticket.Id); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			newErr := fmt.Errorf(fmt.Sprintf("SQL Error: %s, Detail: %s, Where: %s, Code: %s, SQLState: %s", pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState()))
			r.logger.Error(newErr)
			return ticket, newErr
		}
		return ticket, err
	}

	return ticket, nil
}

func (r repository) FindAll(ctx context.Context) (a []ticket.Ticket, err error) {
	q := `
		SELECT ticket.id, type.typename, event.name, price, purchased 
			FROM public.ticket
			join public.event on ticket.eventid = event.id
			join public.type on ticket.typeid = type.id;
	`

	r.logger.Trace(fmt.Sprintf("SQL Query: %s", formatQuery(q)))

	rows, err := r.client.Query(ctx, q)
	if err != nil {
		return nil, err
	}

	tickets := make([]ticket.Ticket, 0)

	for rows.Next() {
		var t ticket.Ticket

		err = rows.Scan(&t.Id, &t.TypeID, &t.EventID, &t.Price, &t.Purchased)
		if err != nil {
			return nil, err
		}

		tickets = append(tickets, t)

		if err = rows.Err(); err != nil {
			return nil, err
		}
	}

	return tickets, nil
}

func (r repository) FindOne(ctx context.Context, id string) (ticket.Ticket, error) {
	q := `
		SELECT ticket.id, type.typename, event.name, price, purchased 
			FROM public.ticket
			join public.event on ticket.eventid = event.id
			join public.type on ticket.typeid = type.id
		WHERE ticket.id = $1;
	`

	r.logger.Trace(fmt.Sprintf("SQL Query: %s", formatQuery(q)))

	var t ticket.Ticket
	queryRow := r.client.QueryRow(ctx, q, id)
	err := queryRow.Scan(&t.Id, &t.TypeID, &t.EventID, &t.Price, &t.Purchased)
	if err != nil {
		return ticket.Ticket{}, err
	}
	return t, nil
}

func NewRepository(client *pgxpool.Pool, logger *logging.Logger) ticket.Repository {
	return &repository{
		client: client,
		logger: logger,
	}
}
