package handler

import (
	"fmt"
	"net/http"
	"github.com/gorilla/mux"
)

var rootRouter *mux.Router = nil

func init() {
	r:= Router()
	r.HandleFunc("/", rootHandler)
}

func Router() *mux.Router {
	if rootRouter == nil {
		rootRouter = mux.NewRouter()
		http.Handle("/", rootRouter)
	}
	return rootRouter
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello, world!")
}
