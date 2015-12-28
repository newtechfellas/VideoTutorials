package handler
import (
	"net/http"
	"fmt"
	"tutorials"
	"tutorials/entity"
	"log"
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
	log.Println(&c)
}