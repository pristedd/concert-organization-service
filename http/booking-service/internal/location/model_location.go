package location

import "database/sql"

type Location struct {
	ID       string         `json:"id"`
	Name     string         `json:"name"`
	Address  string         `json:"address"`
	SeatsNum int            `json:"seats_num"`
	Comment  sql.NullString `json:"comment"`
}
