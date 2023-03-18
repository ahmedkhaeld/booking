package handlers

import (
	"fmt"
	"github.com/ahmedkhaeld/booking/data"
	"github.com/ahmedkhaeld/jazz/forms"
	"github.com/ahmedkhaeld/jazz/mailer"
	"github.com/ahmedkhaeld/jazz/render"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
	"time"
)

func (h *Handlers) Availability(w http.ResponseWriter, r *http.Request) {
	defer h.LoadTime(time.Now())
	err := h.Render.Page(w, r, "search-availability.page.tmpl", nil, nil)
	if err != nil {
		h.ErrorLog.Println("error rendering:", err)
	}
}

// PostAvailability is used to process the form submission from the search-availability page
// It will parse the form, check the start date &end date if available then send the user to the make-reservation page.
func (h *Handlers) PostAvailability(w http.ResponseWriter, r *http.Request) {
	start := r.Form.Get("start")
	end := r.Form.Get("end")

	layout := "2006-01-02"
	startDate, err := time.Parse(layout, start)
	if err != nil {
		h.ErrorStatus(w, http.StatusInternalServerError)
		return
	}
	endDate, err := time.Parse(layout, end)
	if err != nil {
		h.ErrorStatus(w, http.StatusInternalServerError)
		return
	}

	rooms, err := h.Models.Rooms.GetAnyAvailable(startDate, endDate)
	if err != nil {
		h.ErrorStatus(w, http.StatusInternalServerError)
		return
	}

	if len(rooms) == 0 {
		h.Session.Put(r.Context(), "error", "No rooms available")
		http.Redirect(w, r, "/check/rooms", http.StatusSeeOther)
		return
	}

	// start build with the reservation [arrival, departure] in the session
	reservation := data.Reservation{
		StartDate: startDate,
		EndDate:   endDate,
	}
	h.Session.Put(r.Context(), "reservation", reservation)

	// pass the rooms to the choose-room template
	d := make(map[string]interface{})
	d["rooms"] = rooms
	td := &render.TemplateData{
		Data: d,
	}

	err = h.Render.Page(w, r, "available-rooms.page.tmpl", nil, td)

}

func (h *Handlers) ChooseRoom(w http.ResponseWriter, r *http.Request) {

	id := chi.URLParam(r, "id")
	roomID, err := strconv.Atoi(id)
	if err != nil {
		h.ErrorStatus(w, http.StatusInternalServerError)
		return
	}

	// get the reservation from the session
	reservation, ok := h.Session.Get(r.Context(), "reservation").(data.Reservation)
	if !ok {
		h.ErrorStatus(w, http.StatusInternalServerError)
		return
	}
	//get the room name by roomID and put it in the reservation
	room, err := h.Models.Rooms.GetById(roomID)
	if err != nil {
		h.ErrorStatus(w, http.StatusInternalServerError)
		return
	}

	//update the reservation with the room name and room id
	reservation.Room.Name = room.Name
	reservation.RoomID = roomID

	// put the reservation in the session
	h.Session.Put(r.Context(), "reservation", reservation)

	http.Redirect(w, r, "/bookings/reservation", http.StatusSeeOther)
}

///-----------------Reservation Processing-----------------///

func (h *Handlers) Reservation(w http.ResponseWriter, r *http.Request) {
	//pull the reservation out of the session
	reservation, ok := h.Session.Get(r.Context(), "reservation").(data.Reservation)
	if !ok {
		h.Session.Put(r.Context(), "error", "can't get reservation from session")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	//get the room from the database
	room, err := h.Models.Rooms.GetById(reservation.RoomID)
	if err != nil {
		h.ErrorLog.Println("error getting room by id:", err)
		return
	}
	reservation.Room.Name = room.Name

	//update the reservation in the session
	h.Session.Put(r.Context(), "reservation", reservation)

	//type cast the start date and end date from time.Time to string and pass them to td
	//this is because the template is expecting a string
	sd := reservation.StartDate.Format("2006-01-02")
	ed := reservation.EndDate.Format("2006-01-02")

	stringData := make(map[string]string)
	stringData["start_date"] = sd
	stringData["end_date"] = ed

	d := make(map[string]interface{})
	d["reservation"] = reservation
	td := &render.TemplateData{
		Form:       forms.New(nil),
		StringData: stringData,
		Data:       d,
	}

	err = h.Render.Page(w, r, "make-reservation.page.tmpl", nil, td)
	if err != nil {
		h.ErrorLog.Println("error rendering:", err)
	}
}

func (h *Handlers) PostReservation(w http.ResponseWriter, r *http.Request) {

	//pull the reservation out of the session
	reservation, ok := h.Session.Get(r.Context(), "reservation").(data.Reservation)
	if !ok {
		h.ErrorLog.Println("error getting reservation from session")
		return
	}
	err := r.ParseForm()
	if err != nil {
		h.ErrorLog.Println("error parsing form:", err)
		return
	}

	reservation.FirstName = r.Form.Get("first_name")
	reservation.LastName = r.Form.Get("last_name")
	reservation.Email = r.Form.Get("email")
	reservation.Phone = r.Form.Get("phone")

	//validate the user's input
	form := forms.New(r.PostForm)
	reservation.Validate(form)
	if !form.Valid() {
		d := make(map[string]interface{})
		d["reservation"] = reservation
		h.Render.Page(w, r, "make-reservation.page.tmpl", nil, &render.TemplateData{
			Form: form,
			Data: d,
		})
		return
	}

	//begin transaction

	//insert the reservation into the database
	newResID, err := h.Models.Reservations.Create(reservation)
	if err != nil {
		h.ErrorLog.Println("error inserting reservation:", err)
		return
	}

	reservation.ID = newResID
	// build a restriction related to the added reservation
	restriction := data.Restriction{
		StartDate:     reservation.StartDate,
		EndDate:       reservation.EndDate,
		RoomID:        reservation.RoomID,
		ReservationID: newResID,
	}
	//insert the restriction into the database
	err = h.Models.Restrictions.Create(restriction)
	if err != nil {
		h.ErrorLog.Println("error inserting restriction:", err)
		return
	}

	//send notification to the guest and the owner
	var content struct {
		Name string
		Body string
	}
	content.Name = reservation.FirstName + " " + reservation.LastName
	content.Body = fmt.Sprintf("This is confirm your resrvation from %s to %s. At %s",
		reservation.StartDate.Format("2006-01-02"), reservation.EndDate.Format("2006-01-02"), reservation.Room.Name)
	msg := mailer.Message{
		From:        "breadandbreakfast@booking.com",
		To:          reservation.Email,
		Subject:     "Reservation Confirmation",
		Template:    "mail",
		Attachments: nil,
		Data:        content,
	}

	h.Mailer.Jobs <- msg
	res := <-h.Mailer.Results
	if res.Error != nil {
		h.ErrorLog.Println(res.Error)
	}

	h.Session.Put(r.Context(), "reservation", reservation)
	//redirect to prevent the client to submit the form again
	//[good practice any time we are using post request]
	http.Redirect(w, r, "/booking/reservation-summary", http.StatusSeeOther)
}

// ReservationSummary displays the reservation summary page to the user
func (h *Handlers) ReservationSummary(w http.ResponseWriter, r *http.Request) {
	//pull the reservation out of the session
	reservation, ok := h.Session.Get(r.Context(), "reservation").(data.Reservation)
	if !ok {
		h.Session.Put(r.Context(), "error", "can't get reservation from session")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	h.Session.Remove(r.Context(), "reservation")

	//type cast back the start date and end date from time.Time to string and pass them to td
	//this is because the template is expecting a string
	sd := reservation.StartDate.Format("2006-01-02")
	ed := reservation.EndDate.Format("2006-01-02")

	stringData := make(map[string]string)
	stringData["start_date"] = sd
	stringData["end_date"] = ed

	d := make(map[string]interface{})
	d["reservation"] = reservation
	td := &render.TemplateData{
		Form:       forms.New(nil),
		StringData: stringData,
		Data:       d,
	}

	err := h.Render.Page(w, r, "reservation-summary.page.tmpl", nil, td)
	if err != nil {
		h.ErrorLog.Println("error rendering:", err)
	}
}
