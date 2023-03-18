package handlers

import (
	"github.com/ahmedkhaeld/booking/data"
	"github.com/ahmedkhaeld/jazz"
	"github.com/ahmedkhaeld/jazz/render"
	"net/http"
	"time"
)

type Handlers struct {
	*jazz.Jazz
	data.Models
}

func (h *Handlers) Home(w http.ResponseWriter, r *http.Request) {
	defer h.LoadTime(time.Now())
	err := h.Render.Page(w, r, "home.page.tmpl", nil, nil)
	if err != nil {
		h.ErrorLog.Println("error rendering:", err)
	}
}

func (h *Handlers) About(w http.ResponseWriter, r *http.Request) {
	defer h.LoadTime(time.Now())
	err := h.Render.Page(w, r, "about.page.tmpl", nil, nil)
	if err != nil {
		h.ErrorLog.Println("error rendering:", err)
	}
}

func (h *Handlers) Contact(w http.ResponseWriter, r *http.Request) {
	defer h.LoadTime(time.Now())
	err := h.Render.Page(w, r, "contact.page.tmpl", nil, nil)
	if err != nil {
		h.ErrorLog.Println("error rendering:", err)
	}
}

func (h *Handlers) Rooms(w http.ResponseWriter, r *http.Request) {

	//get all rooms
	rooms, err := h.Models.Rooms.GetAll()
	if err != nil {
		h.ErrorLog.Println("error getting rooms:", err)
		return
	}

	if len(rooms) == 0 {
		h.Session.Put(r.Context(), "error", "No rooms available")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	d := make(map[string]interface{})
	d["rooms"] = rooms
	td := &render.TemplateData{
		Data: d,
	}
	defer h.LoadTime(time.Now())
	err = h.Render.Page(w, r, "rooms.page.tmpl", nil, td)
	if err != nil {
		h.ErrorLog.Println("error rendering:", err)
	}
}

func (h *Handlers) Generals(w http.ResponseWriter, r *http.Request) {
	defer h.LoadTime(time.Now())
	err := h.Render.Page(w, r, "generals.page.tmpl", nil, nil)
	if err != nil {
		h.ErrorLog.Println("error rendering:", err)
	}
}

func (h *Handlers) Majors(w http.ResponseWriter, r *http.Request) {
	defer h.LoadTime(time.Now())
	err := h.Render.Page(w, r, "majors.page.tmpl", nil, nil)
	if err != nil {
		h.ErrorLog.Println("error rendering:", err)
	}
}
