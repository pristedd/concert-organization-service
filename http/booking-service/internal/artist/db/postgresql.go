package artist

import (
	"booking_service/internal/artist"
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

func (r repository) Create(ctx context.Context, artist *artist.Artist) (*artist.Artist, error) {
	q := `INSERT INTO public.artist (name, email, phone)
			VALUES ($1, $2, $3) 
			RETURNING id`

	r.logger.Trace(fmt.Sprintf("SQL Query: %s", formatQuery(q)))
	queryRow := r.client.QueryRow(ctx, q, artist.Name, artist.Email, artist.Phone)
	//
	if err := queryRow.Scan(&artist.ID); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			newErr := fmt.Errorf(fmt.Sprintf("SQL Error: %s, Detail: %s, Where: %s, Code: %s, SQLState: %s", pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState()))
			r.logger.Error(newErr)
			return artist, newErr
		}
		return artist, err
	}

	return artist, nil
}

func (r repository) FindAll(ctx context.Context) (a []artist.Artist, err error) {
	q := `
		SELECT id, name, email, phone FROM public.artist;
	`

	r.logger.Trace(fmt.Sprintf("SQL Query: %s", formatQuery(q)))

	rows, err := r.client.Query(ctx, q)
	if err != nil {
		return nil, err
	}

	artists := make([]artist.Artist, 0)

	for rows.Next() {
		var a artist.Artist
		err = rows.Scan(&a.ID, &a.Name, &a.Email, &a.Phone)
		if err != nil {
			return nil, err
		}

		artists = append(artists, a)

		if err = rows.Err(); err != nil {
			return nil, err
		}
	}

	return artists, nil
}

func (r repository) FindOne(ctx context.Context, id string) (artist.Artist, error) {
	q := `
		SELECT id, name, email, phone FROM public.artist WHERE id = $1;
	`

	r.logger.Trace(fmt.Sprintf("SQL Query: %s", formatQuery(q)))

	var a artist.Artist
	queryRow := r.client.QueryRow(ctx, q, id)
	err := queryRow.Scan(&a.ID, &a.Name, &a.Email, &a.Phone)
	if err != nil {
		return artist.Artist{}, err
	}
	return a, nil
}

func (r repository) Update(ctx context.Context, artist *artist.Artist, id string) error {
	q := `
	UPDATE public.artist   
	SET name = $1, email = $2, phone = $3
	Where id = $4;
	`

	r.logger.Trace(fmt.Sprintf("SQL Query: %s", formatQuery(q)))

	_, err := r.client.Exec(ctx, q, artist.Name, artist.Email, artist.Phone, id)
	if err != nil {
		return err
	}
	return nil
}

func (r repository) Delete(ctx context.Context, id string) error {
	q := `
	DELETE FROM public.artist WHERE id=$1
	`

	r.logger.Trace(fmt.Sprintf("SQL Query: %s", formatQuery(q)))

	_, err := r.client.Exec(ctx, q, id)
	if err != nil {
		return err
	}

	return nil
}

func NewRepository(client *pgxpool.Pool, logger *logging.Logger) artist.Repository {
	return &repository{
		client: client,
		logger: logger,
	}
}
