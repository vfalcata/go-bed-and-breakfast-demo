package handlers

// repository pattern, allofehwr us to swap componnets out of our application with minimal changes required to our code base
// we want to ling all the function for handlers together with the repository so that ALL the handlers have access to the repository
// Repository Pattern allows us to share information that parts of our application need easily. If we need to add something we just do so to the repo, in this case our app config struct. This will make it immediatly available to every other part of the application that has access to it.
import (
	"encoding/json"
	"errors"
	"log"
	"myapp/internal/config"
	"myapp/internal/driver"
	"myapp/internal/forms"
	"myapp/internal/helpers"
	"myapp/internal/models"
	"myapp/internal/render"
	"myapp/internal/repository"
	"myapp/internal/repository/dbrepo"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
)

// repository that will be used by the handler, public var
var Repo *Repository

// The repository type
type Repository struct {
	App *config.AppConfig
	DB  repository.DatabaseRepo
}

// Allows us to create a new repo
func NewRepo(a *config.AppConfig, db *driver.DB) *Repository {
	// db is the database connection pool

	return &Repository{
		App: a,
		DB:  dbrepo.NewPostgresRepo(db.SQL, a),
	}
}

// allows us to set the repository for the handlers
func NewHandlers(r *Repository) {
	Repo = r
}

// In order for a function to respond to a request from a web browser it has to handle TWO PARAMETERS
func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	remoteIP := r.RemoteAddr
	m.App.Session.Put(r.Context(), "remote_ip", remoteIP) // the request context

	// adding this receiver in this way allows it to have access to everything inside repository, such as our AppConfig
	render.Template(w, r, "home.page.tmpl", &models.TemplateData{})
}

func (m *Repository) About(w http.ResponseWriter, r *http.Request) {

	/* When we were just experimenting with session
	stringMap := make(map[string]string)
	stringMap["test"] = "Why Hellooo Againnn"

	remoteIP := m.App.Session.GetString(r.Context(), "remote_ip") // will be empty if there is nothing in the session named remoteIP
	stringMap["remote_ip"] = remoteIP
	*/
	render.Template(w, r, "about.page.tmpl", &models.TemplateData{})

}

// Reservation renders the make a reservation page and displays form
func (m *Repository) Reservation(w http.ResponseWriter, r *http.Request) {

	// first we create empty form data
	res, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		helpers.ServerError(w, errors.New("cannot get reservation from session"))
		return
	}

	room, err := m.DB.GetRoomByID(res.RoomID) // get the room name via id
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	res.Room.RoomName = room.RoomName                  // now we have the room name
	m.App.Session.Put(r.Context(), "reservation", res) // store the reservation in the session
	sd := res.StartDate.Format("2006-01-02")
	ed := res.StartDate.Format("2006-01-02")

	stringMap := make(map[string]string)
	stringMap["start_date"] = sd
	stringMap["end_date"] = ed

	data := make(map[string]interface{})
	data["reservation"] = res

	// then we pass that form data in to the render template, along with a new empty form
	// When we render the reservation, we need to give it a form object, via the template
	render.Template(w, r, "make-reservation.page.tmpl", &models.TemplateData{
		Form:      forms.New(nil), // The very first time the form is renderd on the page, will then give access to that forms object
		Data:      data,
		StringMap: stringMap,
	})
}

// post the reservation form
func (m *Repository) PostReservation(w http.ResponseWriter, r *http.Request) {
	reservation, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		helpers.ServerError(w, errors.New("can't get from session"))
		return
	}

	// We need to check the form before it is submitted, but we need to display it before it gets submitted, we need to add a field to template data for the form.
	err := r.ParseForm() // first parse the form
	if err != nil {
		/* Previously we just printed out the error, but now we will replace this with our error package helper
		log.Println(err)
		*/
		helpers.ServerError(w, err) // no we get detailed error messages
		return
	}

	/* // we already populated this information, so it is no longer needed here
	sd := r.Form.Get("start_date")
	ed := r.Form.Get("end_date")

	layout := "2006-01-02"
	startDate, err := time.Parse(layout, sd)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	endDate, err := time.Parse(layout, ed)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	roomID, err := strconv.Atoi(r.Form.Get("room_id")) // A to i, for alpha to integer
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	*/

	// pull data from request form (which had been parsed above), and store that data
	reservation.FirstName = r.Form.Get("first_name")
	reservation.LastName = r.Form.Get("last_name")
	reservation.Phone = r.Form.Get("phone")
	reservation.Email = r.Form.Get("email")
	/* //replaced by above block
	reservation := models.Reservation{
		FirstName: r.Form.Get("first_name"),
		LastName:  r.Form.Get("last_name"),
		Phone:     r.Form.Get("phone"),
		Email:     r.Form.Get("email"),
		// now add ther rest of the values for our model
		StartDate: startDate,
		EndDate:   endDate,
		RoomID:    roomID,
	}
	*/
	// how do we now check our form, first we create a new form object ,and send back a pointer to it then we can check things
	form := forms.New(r.PostForm)

	// Now with the form object check if it is valid
	// form.Has("first_name", r) // does this form has the value first name and is more than just an empty string. If it does have an error attach attach
	// this replaces the above
	// errors are stored here fifo, so if multiple errors proc, the lowest index one triggers first

	form.Required("first_name", "last_name", "email")
	form.MinLength("first_name", 3)
	form.IsEmail("email")

	if !form.Valid() {
		data := make(map[string]interface{})
		data["reservation"] = reservation // add reservation to data
		// If a field is empty we cannot simply return it to the resevation handler because that is creating an empty form. So instead we can render it here
		render.Template(w, r, "make-reservation.page.tmpl", &models.TemplateData{
			Form: form, // now we have the errors if any
			Data: data, // we also have the form data
		})
		return
	}

	// now write form to database
	newReservationID, err := m.DB.InsertReservation(reservation)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	restriction := models.RoomRestriction{
		StartDate:     reservation.StartDate,
		EndDate:       reservation.EndDate,
		RoomID:        reservation.RoomID,
		ReservationID: newReservationID,
		RestrictionID: 1,
	}

	err = m.DB.InsertRoomRestriction(restriction)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	// Now we save the reservation to the session
	m.App.Session.Put(r.Context(), "reservation", reservation)
	// now we want to do a redirect to the session summary and we DO NOT want people to submit the form twice

	http.Redirect(w, r, "/reservation-summary", http.StatusSeeOther)
	// http.StatusSeeOther, is code 303 which is appropriate for a redirect
}

// Generals renders the room page
func (m *Repository) Generals(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "generals.page.tmpl", &models.TemplateData{})
}

// Majors renders the room page
func (m *Repository) Majors(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "majors.page.tmpl", &models.TemplateData{})
}

// Availability renders the search availability page
func (m *Repository) Availability(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "search-availability.page.tmpl", &models.TemplateData{})
}

// handles post req for search availability
func (m *Repository) PostAvailability(w http.ResponseWriter, r *http.Request) {
	start := r.Form.Get("start") // grab the post data from the request specifically the field named "start" the value is generally a string
	end := r.Form.Get("end")

	layout := "2006-01-02"
	startDate, err := time.Parse(layout, start)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	endDate, err := time.Parse(layout, end)
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

	// pull data of rooms that are not reserved
	data := make(map[string]interface{})
	data["rooms"] = rooms

	// store start and end date to session, so we can make a reservation
	res := models.Reservation{
		StartDate: startDate,
		EndDate:   endDate,
	}

	m.App.Session.Put(r.Context(), "reservation", res)

	render.Template(w, r, "choose-room.page.tmpl", &models.TemplateData{
		Data: data,
	})

	/* This was only for testing, no longer needed
	w.Write([]byte(fmt.Sprintf("Start date is %s and end date is %s", start, end))) // cast string to byte
	*/
}

// we will never use this type outside of this file
// IMPORTANT RULES...If you want to export a struct to JSON, the member names MUST START WITH A CAPITAL LETTER, (they have to be exported in other words)
type jsonResponse struct {
	OK        bool   `json:"ok"` // explicitly tell go what you want to use as the values in the JSON response. In this case Anytime we are populating the value of, OK, we want it to return in JSON so we need the backticks to tell go what we want the field to be called in JSON. Go can kind of make assumptions if we do not specify these back ticks but you should not always trust it.
	Message   string `json:"message"`
	RoomID    string `json:"room_id"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

// Json request for availability
func (m *Repository) AvailabilityJSON(w http.ResponseWriter, r *http.Request) {

	sd := r.Form.Get("start")
	ed := r.Form.Get("end")

	layout := "2006-01-02"
	startDate, _ := time.Parse(layout, sd)
	endDate, _ := time.Parse(layout, ed)

	roomID, _ := strconv.Atoi(r.Form.Get("room_id"))

	available, _ := m.DB.SearchAvailabilityByDatesByRoomID(startDate, endDate, roomID)

	resp := jsonResponse{
		OK:        available,
		Message:   "",
		StartDate: sd,
		EndDate:   ed,
		RoomID:    strconv.Itoa(roomID),
	}

	out, err := json.MarshalIndent(resp, "", "     ") // will marshall our JSON and indent it so it is pretty
	if err != nil {
		// log.Println(err)
		helpers.ServerError(w, err)
		return
	}
	// We need to tell the web browser that is receiving the response, what kind of resoehponse we are sending it. Whenever we call a webpage we are automatically getting a header sent back to you by the server that says this is of type text HTML, mehbut we no want to send a type of application JSON, which is the standard header for JSON files
	w.Header().Set("Content-Type", "application/json") // these arguments have to be exact (they are a key value pair)
	w.Write(out)
}

// Contact renders the contact page
func (m *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "contact.page.tmpl", &models.TemplateData{})
}

// After user makes a reservation, the summary is generated
// we need to pass post reservation info from reservation page to this summary page. We can do this via the session
func (m *Repository) ReservationSummary(w http.ResponseWriter, r *http.Request) {
	reservation, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	// now the value of the variable "reservation" has no idea what type it is. we NEED to type assert to "models.Reservation"
	// Now if this works  it implies it managed to find something called a reservation in the session and it manages to assert it to type "models.Reservation" then the value of "ok" will be true else false
	// otherwise we want to do something with the successful "session" so we pass it as template data
	if !ok {
		// log.Println("Cannot get item from the session!!")
		m.App.ErrorLog.Println("Can't get error from session")
		m.App.Session.Put(r.Context(), "error", "Can't get reservation from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect) // as this is for error only
		return
	}

	// once we get the reservation data from the session we no longer need it so we can simply remove it
	m.App.Session.Remove(r.Context(), "reservation")

	data := make(map[string]interface{})
	data["reservation"] = reservation

	sd := reservation.StartDate.Format("2006-01-02")
	ed := reservation.EndDate.Format("2006-01-02")

	stringMap := make(map[string]string)
	stringMap["start_date"] = sd
	stringMap["end_date"] = ed
	render.Template(w, r, "reservation-summary.page.tmpl", &models.TemplateData{
		Data:      data,
		StringMap: stringMap,
	})
}

// get list of available rooms
func (m *Repository) ChooseRoom(w http.ResponseWriter, r *http.Request) {
	roomID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	res, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		helpers.ServerError(w, err)
		return
	}
	res.RoomID = roomID
	m.App.Session.Put(r.Context(), "reservation", res)
	http.Redirect(w, r, "/make-reservation", http.StatusSeeOther)
}

// on the generals page after inputting a date it will make you get a url
// the url has parameters, where we will extract here and store in the session
func (m *Repository) BookRoom(w http.ResponseWriter, r *http.Request) {
	// id, s and e are the url params
	roomID, _ := strconv.Atoi(r.URL.Query().Get("id"))
	sd := r.URL.Query().Get("s")
	ed := r.URL.Query().Get("e")

	layout := "2006-01-02"
	startDate, _ := time.Parse(layout, sd)
	endDate, _ := time.Parse(layout, ed)

	var res models.Reservation

	room, err := m.DB.GetRoomByID(roomID) // get the room name via id
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	res.Room.RoomName = room.RoomName
	res.RoomID = roomID
	res.StartDate = startDate
	res.EndDate = endDate

	log.Println(roomID, sd, ed)

	m.App.Session.Put(r.Context(), "reservation", res)
	http.Redirect(w, r, "/make-reservation", http.StatusSeeOther)

}
