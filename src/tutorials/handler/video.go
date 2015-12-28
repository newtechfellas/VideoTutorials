package handler
import (
	"net/http"
	"fmt"
	"encoding/json"
	"tutorials/entity"
	"log"
	"github.com/astaxie/beego/validation"
	"bytes"
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
	//	var c entity.Course
	var err error
	if _, err = validateCourseFromPost(w, r); err != nil {
		return //http response is already handled by validateCourseFromPost
	}


}

func validateCourseFromPost(w http.ResponseWriter, r *http.Request) (c entity.Course, err error) {
	decoder := json.NewDecoder(r.Body)
	if err = decoder.Decode(&c); err != nil {
		http.Error(w, "Invalid Course details in request body", 400)
		return
	}
	valid := validation.Validation{}
	var b bool
	if b, err = valid.RecursiveValid(c); err != nil || !b {
		var buffer bytes.Buffer
		if valid.HasErrors() {
			for _, err := range valid.Errors {
				buffer.WriteString(err.Field + " " + err.Message+".\t")
			}
		}
		http.Error(w, buffer.String(), 400)
	}
	log.Println("Received Course Data =")
	log.Println(&c)
	return
}