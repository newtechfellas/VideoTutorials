package tutorials
import (
	"fmt"
	"google.golang.org/appengine"
	"net/http"
	"log"
)

func init() {
	http.HandleFunc("/", handler)
}

func handler(w http.ResponseWriter, r *http.Request) {
	log.Println("Inside init")
	appengine.NewContext(r)
	fmt.Fprint(w, "Hello, world!")
}
