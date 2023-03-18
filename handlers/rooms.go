package handlers

import (
	"encoding/json"
	"github.com/ahmedkhaeld/booking/data"
	"net/http"
	"strconv"
	"time"
)

type response struct {
	Ok        bool   `json:"ok"`
	Message   string `json:"message"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
	RoomID    string `json:"room_id"`
}

// AvailabilityJSON handles request for availability from client side [Check availability button]
// takes start and end date, process them, search the database, and send the response back to the client
func (h *Handlers) AvailabilityJSON(w http.ResponseWriter, r *http.Request) {
	//  parse request body
	err := r.ParseForm()
	if err != nil {
		// can't parse form, so return appropriate json
		resp := response{
			Ok:      false,
			Message: "Internal server error",
		}

		out, _ := json.MarshalIndent(resp, "", "     ")
		w.Header().Set("Content-Type", "application/json")
		w.Write(out)
		return
	}

	sd := r.Form.Get("start")
	ed := r.Form.Get("end")
	layout := "2006-01-02"
	startDate, _ := time.Parse(layout, sd)
	endDate, _ := time.Parse(layout, ed)

	roomID, _ := strconv.Atoi(r.Form.Get("room_id"))

	available, err := h.Models.Rooms.IsAvailable(roomID, startDate, endDate)
	if err != nil {
		// got a database error, so return appropriate json
		resp := response{
			Ok:      false,
			Message: "Error querying database",
		}

		out, _ := json.MarshalIndent(resp, "", "     ")
		w.Header().Set("Content-Type", "application/json")
		w.Write(out)
		return
	}
	resp := response{
		Ok:        available,
		Message:   "",
		StartDate: sd,
		EndDate:   ed,
		RoomID:    strconv.Itoa(roomID),
	}

	out, _ := json.MarshalIndent(resp, "", "    ")

	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

func (h *Handlers) BookRoom(w http.ResponseWriter, r *http.Request) {
	roomID, _ := strconv.Atoi(r.URL.Query().Get("id"))
	sd := r.URL.Query().Get("s")
	ed := r.URL.Query().Get("e")

	layout := "2006-01-02"
	startDate, _ := time.Parse(layout, sd)
	endDate, _ := time.Parse(layout, ed)

	room, err := h.Models.Rooms.GetById(roomID)
	if err != nil {
		h.ErrorStatus(w, http.StatusInternalServerError)
		return
	}
	var reservation data.Reservation

	reservation.Room.Name = room.Name
	reservation.RoomID = roomID
	reservation.StartDate = startDate
	reservation.EndDate = endDate

	h.Session.Put(r.Context(), "reservation", reservation)

	http.Redirect(w, r, "/bookings/reservation", http.StatusSeeOther)
}
