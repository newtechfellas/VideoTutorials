package tutorials
import (
	"net/http"
	"encoding/json"
	"github.com/astaxie/beego/validation"
	"bytes"
	"log"
	"fmt"
	"errors"
)

func DecodeAndValidate(w http.ResponseWriter, r *http.Request, obj interface{}) (err error) {
	decoder := json.NewDecoder(r.Body)
	if err = decoder.Decode(obj); err != nil {
		log.Println("Error in decoding request body. Error is ", err)
		ErrorResponse(w, errors.New("Invalid json details in request body"), http.StatusBadRequest)
		return
	}
	log.Println("decoded object type and value:", fmt.Sprintf("%T", obj), obj)
	valid := validation.Validation{}
	var b bool
	if b, err = valid.RecursiveValid(obj); err != nil || !b {
		var buffer bytes.Buffer
		if valid.HasErrors() {
			for _, validationErr := range valid.Errors {
				buffer.WriteString(validationErr.Field + " " + validationErr.Message + ".")
			}
		}
		err = errors.New(buffer.String())
		ErrorResponse(w, err, http.StatusBadRequest)
	}
	return
}

func JsonResponse(w http.ResponseWriter, v interface{}, headers map[string]string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	//Any custom headers passed in
	for k, v := range headers {
		w.Header().Set(k, v)
	}
	w.WriteHeader(statusCode)
	if v != nil {
		b, _ := json.Marshal(v)
		fmt.Fprintf(w, "%s", string(b[:]))
	}
}

func ErrorResponse(w http.ResponseWriter, err error, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	b, _ := json.Marshal(map[string]interface{}{
		"message": err.Error(),
	})
	log.Println("Sending error response for \"" + err.Error() + "\" error")
	fmt.Fprintf(w, "%s", string(b[:]))
}

func Jsonify(obj interface{}) string {
	b, err := json.Marshal(obj)
	if err != nil {
		return ""
	}
	return string(b)
}
