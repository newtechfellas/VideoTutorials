package tutorials

import (
	"net/http"
	"google.golang.org/appengine"
	"time"
	"errors"
	"github.com/gorilla/mux"
	"fmt"
)

var rootRouter *mux.Router = nil
func Router() *mux.Router {
	if rootRouter == nil {
		rootRouter = mux.NewRouter()
		http.Handle("/", rootRouter)
	}
	return rootRouter
}

func init() {
	r := Router()
	r.HandleFunc("/", rootHandler)
	r.HandleFunc("/videos", VideoHandler).Methods("GET")
	r.HandleFunc("/NewVideo", Create).Methods("POST")
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello, world!")
}

func VideoHandler(w http.ResponseWriter, r *http.Request) {
	var course []Course
	var err error
	ctx := appengine.NewContext(r)
	if err = GetAllCourses(ctx, &course); err != nil {
		ErrorResponse(w, errors.New("Did not find any courses"), http.StatusNotFound)
		return
	}
	JsonResponse(w, course, nil, http.StatusOK)
}

func Create(w http.ResponseWriter, r *http.Request) {
	var c Course
	if err := DecodeAndValidate(w, r, &c); err != nil {
		return //http response is already handled by DecodeAndValidate
	}
	c.Date = time.Now()
	ctx := appengine.NewContext(r)

	if err := CreateOrUpdate(ctx, &c, "Course"); err != nil {
		ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}
	JsonResponse(w, nil, nil, http.StatusAccepted)
}
