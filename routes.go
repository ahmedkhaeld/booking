package main

import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

func (a *application) routes() *chi.Mux {
	// middleware must come before any routes

	// add routes here
	a.Get("/", a.Handlers.Home)
	a.Get("/about", a.Handlers.About)
	a.Get("/contact", a.Handlers.Contact)

	a.Get("/rooms", a.Handlers.Rooms)
	a.Get("/rooms/generals-quarters", a.Handlers.Generals)
	a.Get("/rooms/majors-suite", a.Handlers.Majors)
	a.Post("/room/check-json", a.Handlers.AvailabilityJSON)
	a.Get("/bookings/room", a.Handlers.BookRoom)

	a.Get("/check/rooms", a.Handlers.Availability)
	a.Post("/check/rooms", a.Handlers.PostAvailability)
	a.Get("/check/rooms/{id}", a.Handlers.ChooseRoom)

	a.Get("/bookings/reservation", a.Handlers.Reservation)
	a.Post("/bookings/reservation", a.Handlers.PostReservation)

	a.Get("/booking/reservation-summary", a.Handlers.ReservationSummary)

	// static routes
	fileServer := http.FileServer(http.Dir("./public"))
	a.Routes.Handle("/public/*", http.StripPrefix("/public", fileServer))

	return a.Routes
}
