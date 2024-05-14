package artist

import "context"

type Repository interface {
	Create(ctx context.Context, artist *Artist) (*Artist, error)
	FindAll(ctx context.Context) (l []Artist, err error)
	FindOne(ctx context.Context, id string) (Artist, error)
	Update(ctx context.Context, artist *Artist, id string) error
	Delete(ctx context.Context, id string) error
}
