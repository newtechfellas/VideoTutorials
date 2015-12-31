package tutorials

import (
	"net/http"
	"google.golang.org/appengine"
	"time"
	"errors"
	"github.com/gorilla/mux"
	"log"
	"html/template"
	"strconv"
)

var homeTemplate *template.Template
func init() {
	log.Println("Inside handler's init")
	r := mux.NewRouter()
	http.Handle("/", r)
	r.HandleFunc("/", rootHandler)
	r.HandleFunc("/refreshCache", RefreshCache).Methods("POST")
	r.HandleFunc("/search", Search).Methods("GET")
	r.HandleFunc("/newCourse", CreateOrUpdateCourse).Methods("POST")
	r.HandleFunc("/addLecture/{courseId:[0-9]+}", AddLecture).Methods("PUT")
	homeTemplate = template.Must(template.ParseFiles("templates/base.html", "templates/contact.html",
		"templates/homePageHeader.html", "templates/homePageMainContent.html", "templates/navbarLinks.html"))
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	err := homeTemplate.ExecuteTemplate(w, "base", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func AddLecture(w http.ResponseWriter, r *http.Request) {
	var l Lecture
	if err := DecodeAndValidate(w, r, &l); err != nil {
		return //http response is already handled by DecodeAndValidate
	}
	vars := mux.Vars(r)
	courseId, _ := strconv.Atoi(vars["courseId"])
	log.Println("Updating courseId ", courseId, " to add lecture ", l)
	ctx := appengine.NewContext(r)
	var c Course
	if err := GetEntity(ctx,int64(courseId),"Course", &c); err != nil {
		JsonResponse(w, nil, nil, http.StatusOK)
		return //if entity not found for given id, silently ignore and return 200. Dont want to give any hint to the hackers
	}
	var existingLecture bool = false
	for index, item := range c.Lectures {
		if item.Link == l.Link {
			c.Lectures[index] = l
			log.Println("This lecture is already part of the course. Ignoring")
			existingLecture = true
		}
	}
	if !existingLecture {
		c.Lectures = append(c.Lectures, l)
	}
	CreateOrUpdate(ctx, &c, "Course", c.Id)
	AddToCache(c)
	JsonResponse(w, nil, nil, http.StatusOK)
	return
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

func CreateOrUpdateCourse(w http.ResponseWriter, r *http.Request) {
	var c Course
	if err := DecodeAndValidate(w, r, &c); err != nil {
		return //http response is already handled by DecodeAndValidate
	}
	c.Date = time.Now()
	if c.Id == 0  {
		c.Id = c.Date.Unix()
	}
	ctx := appengine.NewContext(r)
	if err := CreateOrUpdate(ctx, &c, "Course", c.Id); err != nil {
		ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}
	AddToCache(c) //update the cache
	JsonResponse(w, nil, nil, http.StatusAccepted)
}

func RefreshCache(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	RefreshCourseCache(ctx)
}