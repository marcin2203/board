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
		Addr:         ":1000",
		WriteTimeout: 1 * time.Second,
		ReadTimeout:  1 * time.Second,
	}

	r.HandleFunc("/main-page", sendMainPage)
	r.HandleFunc("/profile", sendProfilePage)
	r.HandleFunc("/tag/{tag}", sendTagPage)
	r.HandleFunc("/info", sendInfoPage)
	r.HandleFunc("/img", sendCatImg)

	r.HandleFunc("/post/{id}", post)
	r.HandleFunc("/login", login)
	r.HandleFunc("/register", register)
	r.HandleFunc("/tags", getTags)
	r.HandleFunc("/user", UserRouter)
	r.HandleFunc("/debug/page", sendDebug)
	r.HandleFunc("/debug/contents", debug)
	log.Fatal(srv.ListenAndServe())

}
