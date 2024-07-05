package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/lucasvictor3/bookingsbackend/driver"
	"github.com/lucasvictor3/bookingsbackend/internal/config"
	"github.com/lucasvictor3/bookingsbackend/internal/forms"
	"github.com/lucasvictor3/bookingsbackend/internal/helpers"
	"github.com/lucasvictor3/bookingsbackend/internal/models"
	"github.com/lucasvictor3/bookingsbackend/internal/repository"
	"github.com/lucasvictor3/bookingsbackend/internal/repository/dbrepo"
	"github.com/lucasvictor3/bookingsbackend/internal/utils"
)

// Repo the repository used by the handlers
var Repo *Repository

type Repository struct {
	App *config.AppConfig
	DB  repository.DatabaseRepo
}

// NewRepo creates a new repository
func NewRepo(a *config.AppConfig, db *driver.DB) *Repository {
	return &Repository{
		App: a,
		DB:  dbrepo.NewPostgresRepo(db.SQL, a),
	}
}

// NewTestRepo creates a new repository for tests
func NewTestRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
		DB:  dbrepo.NewTestingsRepo(a),
	}
}

// NewHandlers sets the repository for the handlers
func NewHandlers(r *Repository) {
	Repo = r
}

// Home is the home page handler
func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	remoteIP := r.RemoteAddr

	// create a cookie
	m.App.Session.Put(r.Context(), "remoteIP", remoteIP)

	utils.RenderTemplate(w, r, "home.page.tmpl", &models.TemplateData{})
}

// About is the about page handler
func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	// perform some logic
	stringMap := map[string]string{}
	stringMap["test"] = "Hello, again"

	// retrieve cookie remoteIP info
	remoteIP := m.App.Session.GetString(r.Context(), "remoteIP")
	stringMap["remoteIP"] = remoteIP

	// send the data to the template
	utils.RenderTemplate(w, r, "about.page.tmpl", &models.TemplateData{
		StringMap: stringMap,
	})
}

// Reservation is the reservation page handler
func (m *Repository) Reservation(w http.ResponseWriter, r *http.Request) {
	res, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		m.App.Session.Put(r.Context(), "error", "can't get reservation from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	sd := res.StartDate.Format("2006-01-02")
	ed := res.EndDate.Format("2006-01-02")

	room, err := m.DB.GetRoomById(res.RoomID)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't find room!")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	res.Room = room

	m.App.Session.Put(r.Context(), "reservation", res)

	stringMap := make(map[string]string)
	stringMap["start_date"] = sd
	stringMap["end_date"] = ed

	data := make(map[string]interface{})
	data["reservation"] = res

	utils.RenderTemplate(w, r, "make-reservation.page.tmpl", &models.TemplateData{
		Form:      forms.New(nil),
		Data:      data,
		StringMap: stringMap,
	})
}

// PostReservation is the post reservation form
func (m *Repository) PostReservation(w http.ResponseWriter, r *http.Request) {
	reservation, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		m.App.Session.Put(r.Context(), "error", "can't get res from session!")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	err := r.ParseForm()
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't parse form")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	reservation.FirstName = r.Form.Get("first_name")
	reservation.LastName = r.Form.Get("last_name")
	reservation.Email = r.Form.Get("email")
	reservation.Phone = r.Form.Get("phone")

	form := forms.New(r.PostForm)

	form.Required("first_name", "last_name", "email")
	form.MinLength("first_name", 4)
	form.IsEmail("email")

	if !form.Valid() {
		data := make(map[string]interface{})
		data["reservation"] = reservation

		http.Error(w, "invalid form data", http.StatusSeeOther)
		utils.RenderTemplate(w, r, "make-reservation.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}

	newReservationId, err := m.DB.InsertReservation(reservation)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't insert reservation into DB!")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	m.App.Session.Put(r.Context(), "reservation", reservation)

	restriction := models.RoomRestriction{
		StartDate:     reservation.StartDate,
		EndDate:       reservation.EndDate,
		RoomID:        reservation.RoomID,
		ReservationID: newReservationId,
		RestrictionID: 1,
	}

	err = m.DB.InsertRoomRestriction(restriction)

	// msg := models.MailData{
	// 	To:      "john@do.ca",
	// 	From:    "me@here.com",
	// 	Subject: "Some subject",
	// 	Content: "Hello, <strong>world</strong>!",
	//  Template: ""
	// }

	// m.App.MailChan <- msg

	if err != nil {
		helpers.ServerError(w, err)
		m.App.Session.Put(r.Context(), "error", "can't insert restriction into DB!")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	http.Redirect(w, r, "/reservation-summary", http.StatusSeeOther)

}

// Generals is the reservation page handler
func (m *Repository) Generals(w http.ResponseWriter, r *http.Request) {
	utils.RenderTemplate(w, r, "generals.page.tmpl", &models.TemplateData{})
}

// Majors is the reservation page handler
func (m *Repository) Majors(w http.ResponseWriter, r *http.Request) {
	utils.RenderTemplate(w, r, "majors.page.tmpl", &models.TemplateData{})
}

// Availability is the reservation page handler
func (m *Repository) Availability(w http.ResponseWriter, r *http.Request) {
	utils.RenderTemplate(w, r, "search-availability.page.tmpl", &models.TemplateData{})
}

// Contact is the contact page handler
func (m *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	utils.RenderTemplate(w, r, "contact.page.tmpl", &models.TemplateData{})
}

// Login is the login page handler
func (m *Repository) Login(w http.ResponseWriter, r *http.Request) {
	form := forms.New(r.PostForm)

	utils.RenderTemplate(w, r, "login.page.tmpl", &models.TemplateData{
		Form: form,
	})
}

// Logout is the logout get handler
func (m *Repository) Logout(w http.ResponseWriter, r *http.Request) {
	_ = m.App.Session.Destroy(r.Context())
	_ = m.App.Session.RenewToken(r.Context())

	http.Redirect(w, r, "/user/login", http.StatusSeeOther)

}

// PostLogin is the POST login handler
func (m *Repository) PostLogin(w http.ResponseWriter, r *http.Request) {
	_ = m.App.Session.RenewToken(r.Context()) // good practice - make sure to renew token in login or logout

	err := r.ParseForm()
	if err != nil {
		log.Println(err)
	}

	email := r.Form.Get("email")
	password := r.Form.Get("password")
	form := forms.New(r.PostForm)
	form.Required("email", "password")
	form.IsEmail("email")

	if !form.Valid() {
		utils.RenderTemplate(w, r, "login.page.tmpl", &models.TemplateData{
			Form: form,
		})

		return
	}

	id, _, err := m.DB.Authenticate(email, password)
	if err != nil {
		log.Println(err)

		m.App.Session.Put(r.Context(), "error", "Invalid login credentials")
		http.Redirect(w, r, "/user/login", http.StatusSeeOther)
		return
	}

	m.App.Session.Put(r.Context(), "user_id", id)
	m.App.Session.Put(r.Context(), "flash", "Logged in successfully!")
	http.Redirect(w, r, "/admin/dashboard", http.StatusSeeOther)

}

// PostAvailability is the post to search availability form action
func (m *Repository) PostAvailability(w http.ResponseWriter, r *http.Request) {
	start := r.Form.Get("start")
	end := r.Form.Get("end")

	startDate, endDate, err := formatStartAndDateToTime(start, end)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	rooms, err := m.DB.SearchAvailabilityForAllRooms(startDate, endDate)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	if len(rooms) == 0 {
		m.App.Session.Put(r.Context(), "error", "No availability")
		http.Redirect(w, r, "/search-availability", http.StatusSeeOther)
		return
	}

	data := make(map[string]interface{})
	data["rooms"] = rooms

	res := models.Reservation{
		StartDate: startDate,
		EndDate:   endDate,
	}

	m.App.Session.Put(r.Context(), "reservation", res)

	utils.RenderTemplate(w, r, "choose-room.page.tmpl", &models.TemplateData{
		Data: data,
	})

}

type jsonResponse struct {
	OK        bool   `json:"ok"`
	Message   string `json:"message"`
	RoomID    string `json:"room_id"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

// AvailabilityJSON handle the request
func (m *Repository) AvailabilityJSON(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		resp := jsonResponse{
			OK:      false,
			Message: "Internal server error",
		}

		out, _ := json.MarshalIndent(resp, "", "        ")
		w.Header().Set("Content-Type", "application/json")
		w.Write(out)
		return
	}

	sd := r.Form.Get("start")
	ed := r.Form.Get("end")

	start_date, end_date, err := formatStartAndDateToTime(sd, ed)
	if err != nil {
		resp := jsonResponse{
			OK:      false,
			Message: "Dates in invalid format",
		}

		out, _ := json.MarshalIndent(resp, "", "        ")
		w.Header().Set("Content-Type", "application/json")
		w.Write(out)
		return
	}

	roomID, _ := strconv.Atoi(r.Form.Get("room_id"))

	available, err := m.DB.SearchAvailabilityByDatesByRoomId(start_date, end_date, roomID)
	if err != nil {
		resp := jsonResponse{
			OK:      false,
			Message: "Error connecting to database",
		}

		out, _ := json.MarshalIndent(resp, "", "      ")
		w.Header().Set("Content-Type", "application/json")
		w.Write(out)
		return
	}

	resp := jsonResponse{
		OK:        available,
		Message:   "",
		StartDate: sd,
		EndDate:   ed,
		RoomID:    strconv.Itoa(roomID),
	}

	out, _ := json.MarshalIndent(resp, "", "    ")

	log.Println(string(out))
	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

// Contact is the reservation page handler
func (m *Repository) ReservationSummary(w http.ResponseWriter, r *http.Request) {
	reservation, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		m.App.ErrorLog.Println("Can't get error from session")
		m.App.Session.Put(r.Context(), "error", "Cant get reservation from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	sd := reservation.StartDate.Format("2006-01-02")
	ed := reservation.EndDate.Format("2006-01-02")

	stringMap := make(map[string]string)
	stringMap["start_date"] = sd
	stringMap["end_date"] = ed

	m.App.Session.Remove(r.Context(), "reservation")
	data := make(map[string]interface{})
	data["reservation"] = reservation

	utils.RenderTemplate(w, r, "summary.page.tmpl", &models.TemplateData{
		Data:      data,
		StringMap: stringMap,
	})
}

// ChooseRoom is the redirect to reservation with selected user data
func (m *Repository) ChooseRoom(w http.ResponseWriter, r *http.Request) {
	roomID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	m.App.Session.Get(r.Context(), "reservation")

	res, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		helpers.ServerError(w, err)
		return
	}

	res.RoomID = roomID

	m.App.Session.Put(r.Context(), "reservation", res)

	http.Redirect(w, r, "/make-reservation", http.StatusSeeOther)
}

// BookRoom takes URL parameters, builds a session variable and takes user to make reservation.
func (m *Repository) BookRoom(w http.ResponseWriter, r *http.Request) {
	roomID := getIntUrlParam("id", w, r)
	startDate := r.URL.Query().Get("s")
	endDate := r.URL.Query().Get("e")

	sd, ed, err := formatStartAndDateToTime(startDate, endDate)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	var res models.Reservation

	res.RoomID = roomID
	res.StartDate = sd
	res.EndDate = ed

	room, err := m.DB.GetRoomById(res.RoomID)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	res.Room = room

	m.App.Session.Put(r.Context(), "reservation", res)

	http.Redirect(w, r, "/make-reservation", http.StatusSeeOther)

}

// - ADMIN

func (m *Repository) AdminDashboard(w http.ResponseWriter, r *http.Request) {
	utils.RenderTemplate(w, r, "admin-dashboard.page.tmpl", &models.TemplateData{})
}

// AdminNewReservations shows all new reservations in admin tool
func (m *Repository) AdminNewReservations(w http.ResponseWriter, r *http.Request) {
	reservations, err := m.DB.AllNewReservations()

	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	data := make(map[string]interface{})
	data["reservations"] = reservations

	utils.RenderTemplate(w, r, "admin-new-reservations.page.tmpl", &models.TemplateData{
		Data: data,
	})
}

// AdminNewReservations shows all reservations in admin tool
func (m *Repository) AdminAllReservations(w http.ResponseWriter, r *http.Request) {
	reservations, err := m.DB.AllReservations()

	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	data := make(map[string]interface{})
	data["reservations"] = reservations

	utils.RenderTemplate(w, r, "admin-all-reservations.page.tmpl", &models.TemplateData{
		Data: data,
	})
}

func (m *Repository) AdminShowReservation(w http.ResponseWriter, r *http.Request) {
	exploded := strings.Split(r.RequestURI, "/")
	id, err := strconv.Atoi(exploded[4])
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	src := exploded[3]

	stringMap := make(map[string]string)
	stringMap["src"] = src

	reservation, err := m.DB.GetReservationById(id)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	data := make(map[string]interface{})
	data["reservation"] = reservation

	utils.RenderTemplate(w, r, "admin-show-reservation.page.tmpl", &models.TemplateData{
		StringMap: stringMap,
		Data:      data,
		Form:      forms.New(nil),
	})
}

func (m *Repository) AdminPostShowReservation(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	exploded := strings.Split(r.RequestURI, "/")
	id, err := strconv.Atoi(exploded[4])
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	src := exploded[3]

	stringMap := make(map[string]string)
	stringMap["src"] = src

	res, err := m.DB.GetReservationById(id)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	res.FirstName = r.Form.Get("first_name")
	res.LastName = r.Form.Get("last_name")
	res.Email = r.Form.Get("email")
	res.Phone = r.Form.Get("phone")

	err = m.DB.UpdateReservation(res)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	m.App.Session.Put(r.Context(), "flash", "Changes saved.")
	http.Redirect(w, r, fmt.Sprintf("/admin/reservations-%s", src), http.StatusSeeOther)
}

func (m *Repository) AdminReservationsCalendar(w http.ResponseWriter, r *http.Request) {
	now := time.Now()

	if r.URL.Query().Get("y") != "" {
		year, _ := strconv.Atoi(r.URL.Query().Get("y"))
		month, _ := strconv.Atoi(r.URL.Query().Get("m"))
		now = time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	}

	next := now.AddDate(0, 1, 0)
	last := now.AddDate(0, -1, 0)

	nextMonth := next.Format("01")
	nextMonthYear := next.Format("2006")

	lastMonth := last.Format("01")
	lastMonthYear := last.Format("2006")

	stringMap := make(map[string]string)

	stringMap["next_month"] = nextMonth
	stringMap["next_month_year"] = nextMonthYear
	stringMap["last_month"] = lastMonth
	stringMap["last_month_year"] = lastMonthYear

	stringMap["this_month"] = now.Format("01")
	stringMap["this_month_year"] = now.Format("2006")

	utils.RenderTemplate(w, r, "admin-reservations-calendar.page.tmpl", &models.TemplateData{
		StringMap: stringMap,
	})
}

// AdminProcessReservation marks a reservation as processed
func (m *Repository) AdminProcessReservation(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	src := chi.URLParam(r, "src")

	log.Println(src)

	_ = m.DB.UpdateProcessedForReservation(id, 1)

	m.App.Session.Put(r.Context(), "flash", "Reservation marked as processed")
	http.Redirect(w, r, fmt.Sprintf("/admin/reservations-%s", src), http.StatusSeeOther)

}

// AdminDeleteReservation deletes a reservation
func (m *Repository) AdminDeleteReservation(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	src := chi.URLParam(r, "src")

	log.Println(src)

	_ = m.DB.DeleteReservation(id)

	m.App.Session.Put(r.Context(), "flash", "Reservation deleted")
	http.Redirect(w, r, fmt.Sprintf("/admin/reservations-%s", src), http.StatusSeeOther)

}

// - HELPERS

func getIntUrlParam(param string, w http.ResponseWriter, r *http.Request) int {
	selectedParam, err := strconv.Atoi(r.URL.Query().Get(param))
	if err != nil {
		helpers.ServerError(w, err)
		return -1
	}
	return selectedParam
}

func formatStartAndDateToTime(start, end string) (time.Time, time.Time, error) {
	// 2020-01-01 ---> 01/02 03:04:05PM '06 -0700

	layout := "2006-01-02"
	startDate, err := time.Parse(layout, start)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}

	endDate, err := time.Parse(layout, end)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}

	return startDate, endDate, nil
}
