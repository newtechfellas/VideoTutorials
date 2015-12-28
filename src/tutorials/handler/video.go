package handler
import (
	"net/http"
	"fmt"
	"tutorials"
	"tutorials/entity"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"time"
)

func init() {
	r := Router()
	r.HandleFunc("/videos", VideoHandler).Methods("GET")
	r.HandleFunc("/NewVideo", Create).Methods("POST")
}

func VideoHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello, GET Videos!")
}

func Create(w http.ResponseWriter, r *http.Request) {
	var c entity.Course
	if err := tutorials.DecodeAndValidate(w, r, &c); err != nil {
		return //http response is already handled by validateCourseFromPost
	}
	c.Date = time.Now()
	ctx := appengine.NewContext(r)
	_, err := datastore.Put(ctx, datastore.NewIncompleteKey(ctx, "Course", nil), &c)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tutorials.JsonResponse(w, nil, nil, http.StatusAccepted)
}