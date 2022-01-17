package api

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"net/http"
	"test/user"
	"test/utils"
	"time"
)

// API provides application resources and handlers.
type API struct {
	Resource *user.UsersResource
}

// NewAPI configures and returns application API.
func NewAPI(dbConnection utils.DbConnection) (*API, error) {
	usersStore, err := user.NewUsersStore(dbConnection.Database, dbConnection.Ctx)
	if err != nil {
		return nil, err
	}
	resource := user.NewUsersResource(*usersStore)

	Api := &API{
		Resource: resource,
	}
	return Api, nil
}

// Router provides application routes.
func (a *API) Router() *chi.Mux {
	r := chi.NewRouter()

	r.Mount("/users", a.Resource.Router())

	return r
}

// New configures application resources and routes.
func NewApp(dbConnection utils.DbConnection) (*chi.Mux, error) {

	api, err := NewAPI(dbConnection)
	if err != nil {
		return nil, err
	}

	//Init the routers
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Timeout(15 * time.Second))
	r.Use(render.SetContentType(render.ContentTypeJSON))

	r.Mount("/", api.Router())

	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})

	return r, nil
}
