package domain

type EventRepository interface {
	ListEvents() ([]Event, error)
	FindEventByID(id string) (*Event, error)
	FindSpotsByEventID(eventID string) ([]*Spot, error)
	FindSpotByName(eventID, name string) (*Spot, error)
	//CreateEvent(event *Event) error
	CreateSpot(spot *Spot) error
	CreateTicket(ticket *Ticket) error
	ReserveSpot(spotID, ticketID string) error
}
