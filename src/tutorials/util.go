package tutorials
import (
	"net/http"
	"encoding/json"
	"github.com/astaxie/beego/validation"
	"bytes"
)


func DecodeAndValidate(w http.ResponseWriter, r *http.Request, obj interface{}) (err error) {
	decoder := json.NewDecoder(r.Body)
	if err = decoder.Decode(obj); err != nil {
		http.Error(w, "Invalid json details in request body", 400)
		return
	}
	valid := validation.Validation{}
	var b bool
	if b, err = valid.RecursiveValid(obj); err != nil || !b {
		var buffer bytes.Buffer
		if valid.HasErrors() {
			for _, err := range valid.Errors {
				buffer.WriteString(err.Field + " " + err.Message + ".\t")
			}
		}
		http.Error(w, buffer.String(), 400)
	}
	return
}
