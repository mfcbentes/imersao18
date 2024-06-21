package main

import (
	"database/sql"
	"net/http"

	httpHandler "github.com/mfcbentes/imersao18/golang/internal/event/infra/http"
	"github.com/mfcbentes/imersao18/golang/internal/event/infra/repository"
	"github.com/mfcbentes/imersao18/golang/internal/event/infra/service"
	"github.com/mfcbentes/imersao18/golang/internal/event/usecase"
)

func main() {
	db, err := sql.Open("mysql", "root:root@tcp(localhost:3306)/imersao18")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	eventRepo, err := repository.NewMysqlEventRepository(db)
	if err != nil {
		panic(err)
	}

	partnerBaseURLs := map[int]string{
		1: "http://localhost:8081/api1",
		2: "http://localhost:8082/api2",
	}

	partnerFactory := service.NewPartnerFactory(partnerBaseURLs)

	listEventsUseCase := usecase.NewListEventsUseCase(eventRepo)
	getEventUseCase := usecase.NewGetEventUseCase(eventRepo)
	listSpotsUseCase := usecase.NewListSpotsUseCase(eventRepo)
	buyTicketsUseCase := usecase.NewBuyTicketsUseCase(eventRepo, partnerFactory)

	eventsHandler := httpHandler.NewEventsHandler(
		listEventsUseCase,
		listSpotsUseCase,
		getEventUseCase,
		buyTicketsUseCase,
	)

	r := http.NewServeMux()
	r.HandleFunc("GET /events", eventsHandler.ListEvents)
	r.HandleFunc("GET /events/{eventID}", eventsHandler.GetEvent)
	r.HandleFunc("GET /events/{eventID}/spots", eventsHandler.ListSpots)
	r.HandleFunc("POST /checkout", eventsHandler.BuyTickets)

	http.ListenAndServe(":8080", r)

}
