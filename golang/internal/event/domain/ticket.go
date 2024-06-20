package domain

import (
	"errors"

	"github.com/google/uuid"
)

var (
	ErrTicketPriceZero   = errors.New("ticket price must be greater than zero")
	ErrInvalidTicketType = errors.New("invalid ticket type")
)

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

func NewTicket(event *Event, spot *Spot, ticketType TicketTipe) (*Ticket, error) {
	if !IsValidTicketType(ticketType) {
		return nil, ErrInvalidTicketType
	}

	ticket := &Ticket{
		ID:         uuid.New().String(),
		EventID:    event.ID,
		Spot:       spot,
		TicketTipe: ticketType,
		Price:      event.Price,
	}

	ticket.CalculatePrice()
	if err := ticket.Validate(); err != nil {
		return nil, err
	}

	return ticket, nil
}

func IsValidTicketType(ticketType TicketTipe) bool {
	return ticketType == TicketTypeFull || ticketType == TicketTypeHalf
}

func (t *Ticket) CalculatePrice() {
	if t.TicketTipe == TicketTypeHalf {
		t.Price /= 2
	}
}

func (t Ticket) Validate() error {
	if t.Price <= 0 {
		return ErrTicketPriceZero
	}

	return nil
}
