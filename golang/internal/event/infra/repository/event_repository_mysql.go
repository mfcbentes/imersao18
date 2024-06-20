package repository

import (
	"database/sql"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/mfcbentes/imersao18/golang/internal/event/domain"
)

type mysqlEventRepository struct {
	db *sql.DB
}

func NewMysqlEventRepository(db *sql.DB) (domain.EventRepository, error) {
	return &mysqlEventRepository{db: db}, nil
}

// CreateSpot implements domain.EventRepository.
func (m *mysqlEventRepository) CreateSpot(spot *domain.Spot) error {
	query := "INSERT INTO spots (id, event_id, name, status, ticket_id)	VALUES (?, ?, ?, ?, ?)"

	_, err := m.db.Exec(query, spot.ID, spot.EventID, spot.Name, spot.Status, spot.TicketID)
	if err != nil {
		return err
	}

	return nil
}

// CreateTicket implements domain.EventRepository.
func (m *mysqlEventRepository) CreateTicket(ticket *domain.Ticket) error {
	query := "INSERT INTO tickets (id, event_id, spot_id, ticket_type, price) VALUES (?, ?, ?, ?, ?)"

	_, err := m.db.Exec(query, ticket.ID, ticket.EventID, ticket.Spot.ID, ticket.TicketTipe, ticket.Price)
	if err != nil {
		return err
	}

	return nil
}

// FindEventByID implements domain.EventRepository.
func (r *mysqlEventRepository) FindEventByID(eventId string) (*domain.Event, error) {
	query := "SELECT id, name, location, organization, rating, date, image_url, capacity, price, partner_id FROM events WHERE id = ?"

	row, err := r.db.Query(query, eventId)
	if err != nil {
		return nil, err
	}
	defer row.Close()

	var event *domain.Event

	err = row.Scan(
		&event.ID,
		&event.Name,
		&event.Location,
		&event.Organization,
		&event.Rating,
		&event.Date,
		&event.ImageURL,
		&event.Capacity,
		&event.Price,
		&event.PartnerID,
	)
	if err != nil {
		return nil, err
	}

	return event, nil

}

// FindSpotByName implements domain.EventRepository.
func (r *mysqlEventRepository) FindSpotByName(eventID string, name string) (*domain.Spot, error) {
	query := `
		SELECT
			s.id, s.event_id, s.name, s.status, s.ticket_id
			t.id, t.event_id, t.spot_id, t.ticket_type, t.price
		FROM spots s
		LEFT JOIN tickets t ON s.id = t.spot_id
		WHERE s.event_id = ? AND s.name = ?
	`

	row := r.db.QueryRow(query, eventID, name)

	var spot domain.Spot
	var ticket domain.Ticket
	var ticketID, ticketEventID, ticketSpotID, ticketType sql.NullString
	var ticketPrice sql.NullFloat64

	err := row.Scan(
		&spot.ID,
		&spot.EventID,
		&spot.Name,
		&spot.Status,
		&spot.TicketID,
		&ticketID,
		&ticketEventID,
		&ticketSpotID,
		&ticketType,
		&ticketPrice,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrSpotNotFound
		}
		return nil, err
	}

	if ticketID.Valid {
		ticket.ID = ticketID.String
		ticket.EventID = ticketEventID.String
		ticket.Spot = &spot
		ticket.TicketTipe = domain.TicketTipe(ticketType.String)
		ticket.Price = ticketPrice.Float64
		spot.TicketID = ticket.ID
	}

	return &spot, nil

}

// FindSpotsByEventID implements domain.EventRepository.
func (m *mysqlEventRepository) FindSpotsByEventID(eventID string) ([]*domain.Spot, error) {
	query := "SELECT id, event_id, name, status, ticket_id FROM spots WHERE event_id = ?"

	rows, err := m.db.Query(query, eventID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var spots []*domain.Spot
	for rows.Next() {
		var spot domain.Spot

		if err := rows.Scan(&spot.ID, &spot.EventID, &spot.Name, &spot.Status, &spot.TicketID); err != nil {
			return nil, err
		}

		spots = append(spots, &spot)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return spots, nil
}

// ListEvents implements domain.EventRepository.
func (r *mysqlEventRepository) ListEvents() ([]domain.Event, error) {
	query := `
		SELECT
			e.id, e.name, e.location, e.organization, e.rating, e.date, e.image_url, e.capacity, e.price, e.partner_id,
			s.id, s.event_id, s.name, s.status, s.ticket_id,
			t.id, t.event_id, t.spot_id, t.ticket_type, t.price
		FROM events e
		LEFT JOIN spots s ON e.id = s.event_id
		LEFT JOIN tickets t ON s.id = t.spot_id
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	eventMap := make(map[string]*domain.Event)
	spotMap := make(map[string]*domain.Spot)

	for rows.Next() {
		var eventID, eventName, eventLocation, eventOrganization, eventRating, eventImageURL, spotID, spotEventID, spotName, spotStatus, spotTicketID, ticketID, ticketEventID, ticketSpotID, ticketType sql.NullString
		var eventDate sql.NullString
		var eventCapacity int
		var eventPrice, ticketPrice sql.NullFloat64
		var partnerID sql.NullInt32

		err := rows.Scan(
			&eventID,
			&eventName,
			&eventLocation,
			&eventOrganization,
			&eventRating,
			&eventDate,
			&eventImageURL,
			&eventCapacity,
			&eventPrice,
			&partnerID,
			&spotID,
			&spotEventID,
			&spotName,
			&spotStatus,
			&spotTicketID,
			&ticketID,
			&ticketEventID,
			&ticketSpotID,
			&ticketType,
			&ticketPrice,
		)
		if err != nil {
			return nil, err
		}

		if !eventID.Valid || !eventName.Valid || !eventLocation.Valid || !eventOrganization.Valid || !eventRating.Valid || !eventDate.Valid || !eventImageURL.Valid || !eventPrice.Valid || !partnerID.Valid {
			continue
		}

		event, exists := eventMap[eventID.String]
		if !exists {
			eventDateParsed, err := time.Parse("2006-01-02 15:04:05", eventDate.String)
			if err != nil {
				return nil, err
			}

			event = &domain.Event{
				ID:           eventID.String,
				Name:         eventName.String,
				Location:     eventLocation.String,
				Organization: eventOrganization.String,
				Rating:       domain.Rating(eventRating.String),
				Date:         eventDateParsed,
				ImageURL:     eventImageURL.String,
				Capacity:     eventCapacity,
				Price:        eventPrice.Float64,
				PartnerID:    int(partnerID.Int32),
				Spots:        []domain.Spot{},
				Tickets:      []domain.Ticket{},
			}

			eventMap[eventID.String] = event
		}

		if spotID.Valid {
			spot, spotExists := spotMap[spotID.String]
			if !spotExists {
				spot = &domain.Spot{
					ID:       spotID.String,
					EventID:  spotEventID.String,
					Name:     spotName.String,
					Status:   domain.SpotStatus(spotStatus.String),
					TicketID: spotTicketID.String,
				}

				event.Spots = append(event.Spots, *spot)
				spotMap[spotID.String] = spot
			}

			if ticketID.Valid {
				ticket := domain.Ticket{
					ID:         ticketID.String,
					EventID:    ticketEventID.String,
					Spot:       spot,
					TicketTipe: domain.TicketTipe(ticketType.String),
					Price:      ticketPrice.Float64,
				}

				event.Tickets = append(event.Tickets, ticket)
			}
		}
	}

	events := make([]domain.Event, 0, len(eventMap))
	for _, event := range eventMap {
		events = append(events, *event)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return events, nil
}

// ReserveSpot implements domain.EventRepository.
func (m *mysqlEventRepository) ReserveSpot(spotID string, ticketID string) error {
	query := "UPDATE spots SET status = ?, ticket_id = ? WHERE id = ?"

	_, err := m.db.Exec(query, domain.SpotStatusSold, ticketID, spotID)
	if err != nil {
		return err
	}

	return nil
}
