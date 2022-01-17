package user

import (
	"errors"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"net/http"
	"net/url"
	"strings"
	"test/utils"
)

//Implements the User management handler

var (
	ErrParamDate = errors.New("Date format error")
)

// UsersStore implements User management handler.
type UsersResource struct {
	Store UsersStore
}

// NewUsersStore creates and returns a User resource.
func NewUsersResource(store UsersStore) *UsersResource {
	return &UsersResource{
		Store: store,
	}
}

// Router for the interaction with User
func (rs *UsersResource) Router() *chi.Mux {
	r := chi.NewRouter()
	r.Get("/", rs.list)
	r.Post("/", rs.create)
	r.Route("/{userID}", func(r chi.Router) {
		r.Put("/", rs.update)
		r.Delete("/", rs.delete)
	})
	return r
}

// Request expected
type userRequest struct {
	User
}

//Binding of the http request to the userRequest
func (ur *userRequest) Bind(r *http.Request) error {
	return nil
}

//Response model for one User
type userResponse struct {
	success bool
	*User
}

//Response model for multiple User
type userListResponse struct {
	success bool
	Users   []User `json:"users"`
	Count   int    `json:"count"`
}

func newUserResponse(u *User, success bool) *userResponse {
	resp := &userResponse{success: success, User: u}
	return resp
}

func newUsersListResponse(u []User, count int, success bool) *userListResponse {
	resp := &userListResponse{
		success: success,
		Users:   u,
		Count:   count,
	}
	return resp
}

// Adds a new User
func (rs *UsersResource) create(w http.ResponseWriter, r *http.Request) {
	//binds body request to User

	uR := &userRequest{}
	if err := render.Bind(r, uR); err != nil {
		utils.Render(w, r, err)
		return
	}
	u := uR.User
	u.Escape()
	//creates it
	err := rs.Store.Create(&u)
	if err != nil {
		utils.Render(w, r, err)
		return
	}
	SendNotification("Created", u)
	render.Respond(w, r, newUserResponse(&u, true))
}

// Update an already existing User
func (rs *UsersResource) update(w http.ResponseWriter, r *http.Request) {
	//gets User ID from URL Parameters
	id := chi.URLParam(r, "userID")

	//binds body request to User
	uR := &userRequest{}
	if err := render.Bind(r, uR); err != nil {

		utils.Render(w, r, err)
		return
	}
	u := uR.User
	u.Escape()

	//update it
	err := rs.Store.Update(id, &u)
	if err != nil {
		utils.Render(w, r, err)
		return
	}
	SendNotification("Updated", u)
	render.Respond(w, r, newUserResponse(&u, true))
}

// Deletes User
func (rs *UsersResource) delete(w http.ResponseWriter, r *http.Request) {
	//gets User ID from URL Parameters
	id := chi.URLParam(r, "userID")

	//delete it
	if err := rs.Store.Delete(id); err != nil {
		utils.Render(w, r, err)
		return
	}
	SendNotification("Deleted", User{ID: id})

	render.Respond(w, r, newUserResponse(&User{ID: id}, true))
}

// Returns filtered User in a list
func (rs *UsersResource) list(w http.ResponseWriter, r *http.Request) {
	//get the url query
	query := r.URL.Query()

	//parses parameters from query
	textS := utils.StringFromQuery("text", query)
	strings.ReplaceAll(url.QueryEscape(textS), "%40", "@") //for email
	idS := utils.StringFromQuery("id", query)
	firstNameS := utils.StringFromQuery("first_name", query)
	lastNameS := utils.StringFromQuery("last_name", query)
	nicknameS := utils.StringFromQuery("nickname", query)
	passwordS := utils.StringFromQuery("password", query)
	emailS := utils.StringFromQuery("email", query)
	strings.ReplaceAll(url.QueryEscape(emailS), "%40", "@") //for email validation
	countryS := utils.StringFromQuery("country", query)
	startDateCreated, err := utils.DateFromQuery("startdcreated", query)
	if err != nil {
		utils.Render(w, r, ErrParamDate)
		return
	}
	endDateCreated, err := utils.DateFromQuery("enddcreated", query)
	if err != nil {
		utils.Render(w, r, ErrParamDate)
		return
	}
	startDateUpdated, err := utils.DateFromQuery("startdupdated", query)
	if err != nil {
		utils.Render(w, r, ErrParamDate)
		return
	}
	endDateUpdated, err := utils.DateFromQuery("enddupdated", query)
	if err != nil {
		utils.Render(w, r, ErrParamDate)
		return
	}
	//get page (number) and page_size
	page, err := utils.Int64FromQuery("page", query)
	if err != nil {
		utils.Render(w, r, err)
		return
	}
	pageSize, err := utils.Int64FromQuery("page_size", query)
	if err != nil {
		utils.Render(w, r, err)
		return
	}

	//gets corresponding entries from db
	usersList, count, err := rs.Store.List(textS, idS, firstNameS, lastNameS, nicknameS, passwordS, emailS, countryS, startDateCreated, endDateCreated, startDateUpdated, endDateUpdated, page, pageSize)
	if err != nil {
		utils.Render(w, r, err)
		return
	}

	render.Respond(w, r, newUsersListResponse(usersList, count, true))
}

//EXample of SendNotification function
func SendNotification(message string, u User) {
	// send String and User to other potential services
}
