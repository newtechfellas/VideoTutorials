package tutorials

import (
	"net/http"
	"google.golang.org/appengine"
	"time"
	"errors"
	"github.com/gorilla/mux"
	"fmt"
	"log"
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
	r.HandleFunc("/refreshCache", RefreshCache).Methods("POST")
	r.HandleFunc("/search", Search).Methods("GET")
	r.HandleFunc("/NewCourse", CreateCourse).Methods("POST")
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello, world!")
}

func Search(w http.ResponseWriter, r *http.Request) {
	var course []Course
	ctx := appengine.NewContext(r)
	course = GetCoursesFromCache(ctx)
	if course == nil || len(course) == 0 {
		ErrorResponse(w, errors.New("Courses not found"), http.StatusNotFound)
		return
	}
	JsonResponse(w, course, nil, http.StatusOK)
}

func CreateCourse(w http.ResponseWriter, r *http.Request) {
	var c Course
	if err := DecodeAndValidate(w, r, &c); err != nil {
		return //http response is already handled by DecodeAndValidate
	}
	if c.Date.IsZero() {
		c.Date = time.Now()
	} else {
	}
	ctx := appengine.NewContext(r)
	if err := CreateOrUpdate(ctx, &c, "Course", c.Date.Unix()); err != nil {
		ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}
	AddToCache(c) //update the cache
	JsonResponse(w, nil, nil, http.StatusAccepted)
}


func RefreshCache(w http.ResponseWriter, r *http.Request) {
	PurgeCache()
	ctx := appengine.NewContext(r)
	LoadCoursesToCache(ctx)
	log.Println("Cache reloaded")
}