package booking

import (
	"database/sql"
	"time"
)

type Booking struct {
	ID         string         `json:"id"`
	LocationID string         `json:"locationID"`
	ArtistID   string         `json:"artistID"`
	Date       time.Time      `json:"date"`
	Comment    sql.NullString `json:"comment"`
}
