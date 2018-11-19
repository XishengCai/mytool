package main

import (
	"fmt"
	"github.com/drone/routes"
	"net/http"
)

func getUser(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	uid := params.Get(":uid")
	fmt.Fprintf(w, "you are get user %s", uid)
}

func modifyUser(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	uid := params.Get(":uid")
	fmt.Fprintf(w, "you are modify user %s", uid)
}

func main() {
	mux := routes.New()
	mux.Get("/user/{id}", getUser)
	http.Handle("/", mux)
	http.ListenAndServe(":8080", nil)
}