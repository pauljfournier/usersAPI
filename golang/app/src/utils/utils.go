package utils

import (
	"context"
	"github.com/go-chi/render"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	errors2 "test/errors"
	"time"
)

//Collect a date from query
func DateFromQuery(dateQueryKey string, query url.Values) (time.Time, error) {
	dateQuery, ok := query[dateQueryKey]
	if ok && len(dateQuery[0]) > 0 {
		return time.Parse(time.RFC3339, dateQuery[0])
	} else {
		return time.Time{}, nil
	}
}

//Collect a string from query then escape it for safety and comparaison to db
func StringFromQuery(queryKey string, query url.Values) string {
	stringQuery, ok := query[queryKey]
	if ok && len(stringQuery[0]) > 0 {
		tmp := stringQuery[0]

		tmp = strings.ReplaceAll(url.QueryEscape(tmp), "+", "%20")
		return strings.ReplaceAll(tmp, ".", "\\.")
	} else {
		return ""
	}
}

//Collect a int from query
func Int64FromQuery(queryKey string, query url.Values) (int64, error) {
	intQuery, ok := query[queryKey]
	if ok && len(intQuery[0]) > 0 {
		intRes, err := strconv.Atoi(intQuery[0])
		if err != nil {
			return 0, err
		}
		return int64(intRes), nil
	}
	return 0, nil
}

func Render(w http.ResponseWriter, r *http.Request, err error) {
	log.Println(err)
	_ = render.Render(w, r, errors2.ErrRender(err))
}

// DbConnection represents the connection to pass around
type DbConnection struct {
	Client   *mongo.Client
	Database *mongo.Database
	Ctx      context.Context
}
