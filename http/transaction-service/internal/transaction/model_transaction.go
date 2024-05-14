package transaction

import "database/sql"

type Transaction struct {
	ID      string         `json:"id"`
	Type    string         `json:"type"`
	Amount  int            `json:"amount"`
	Comment sql.NullString `json:"comment"`
}
