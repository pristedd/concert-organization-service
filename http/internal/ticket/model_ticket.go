package ticket

type Ticket struct {
	Id        string  `json:"id"`
	TypeID    string  `json:"typeID"`
	EventID   string  `json:"eventID"`
	Price     float32 `json:"price"`
	Purchased bool    `json:"purchased"`
}

type Type struct {
	Id       string `json:"id"`
	TypeName string `json:"statusName"`
}
