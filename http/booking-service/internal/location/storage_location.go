package location

import "context"

type Repository interface {
	Create(ctx context.Context, location *Location) (*Location, error)
	FindAll(ctx context.Context) (l []Location, err error)
	FindOne(ctx context.Context, id string) (Location, error)
	Update(ctx context.Context, location *Location, id string) error
	Delete(ctx context.Context, id string) error
}
