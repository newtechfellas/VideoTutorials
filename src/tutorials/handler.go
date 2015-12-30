package tutorials

import (
	"net/http"
	"google.golang.org/appengine"
	"time"
	"errors"
	"github.com/gorilla/mux"
	"log"
	"html/template"
)

var homeTemplate *template.Template
func init() {
	log.Println("Inside handler's init")
	r := mux.NewRouter()
	http.Handle("/", r)
	r.HandleFunc("/", rootHandler)
	r.HandleFunc("/refreshCache", RefreshCache).Methods("POST")
	r.HandleFunc("/search", Search).Methods("GET")
	r.HandleFunc("/newCourse", CreateCourse).Methods("POST")
	homeTemplate = template.Must(template.ParseFiles("templates/base.html", "templates/contact.html",
		"templates/homePageHeader.html", "templates/homePageMainContent.html", "templates/navbarLinks.html"))
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	err := homeTemplate.ExecuteTemplate(w, "base", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
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