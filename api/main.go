package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func Main() {

	r := mux.NewRouter()
	srv := &http.Server{
		Handler:      r,
		Addr:         ":1000",
		WriteTimeout: 1 * time.Second,
		ReadTimeout:  1 * time.Second,
	}

	r.HandleFunc("/main-page", SendMainPage)
	r.HandleFunc("/profile", SendProfilePage)
	r.HandleFunc("/tag/{tag}", SendTagPage)
	r.HandleFunc("/info", SendInfoPage)
	r.HandleFunc("/img", SendCatImg)

	r.HandleFunc("/post/{id}", Post)
	r.HandleFunc("/login", Login)
	r.HandleFunc("/register", Register)
	r.HandleFunc("/tags", GetTags)
	r.HandleFunc("/user", UserRouter)
	r.HandleFunc("/comment", CommentRouter)
	r.HandleFunc("/debug/page", SendDebug)
	r.HandleFunc("/debug/contents", Debug)
	log.Fatal(srv.ListenAndServe())

}
