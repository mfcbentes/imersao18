package domain

type TicketTipe string

const (
	TicketTypeFull TicketTipe = "full"
	TicketTypeHalf TicketTipe = "half"
)

type Ticket struct {
	ID         string
	EventID    string
	Spot       *Spot
	TicketTipe TicketTipe
	Price      float64
}
