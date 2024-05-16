package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func main() {

	r := mux.NewRouter()
	srv := &http.Server{
		Handler:      r,
		Addr:         ":2000",
		WriteTimeout: 1 * time.Second,
		ReadTimeout:  1 * time.Second,
	}

	r.HandleFunc("/main-page", sendMainPage)
	r.HandleFunc("/profile", sendProfilePage)
	r.HandleFunc("/tag", sendTagPage)
	r.HandleFunc("/info", sendInfoPage)
	r.HandleFunc("/img", sendCatImg)

	r.HandleFunc("/login", login)
	r.HandleFunc("/register", register)

	r.HandleFunc("/debug", sendDebug)

	r.HandleFunc("/", myFunc)
	log.Fatal(srv.ListenAndServe())

}
