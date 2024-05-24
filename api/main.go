package main

import (
	"log"
	"net/http"
	"time"

	"main/handlers"

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

	r.HandleFunc("/main-page", handlers.SendMainPage)
	r.HandleFunc("/profile", handlers.SendProfilePage)
	r.HandleFunc("/tag/{tag}", handlers.SendTagPage)
	r.HandleFunc("/info", handlers.SendInfoPage)
	r.HandleFunc("/img", handlers.SendCatImg)

	r.HandleFunc("/post/{id}", handlers.Post)
	r.HandleFunc("/login", handlers.Login)
	r.HandleFunc("/register", handlers.Register)
	r.HandleFunc("/tags", handlers.GetTags)
	r.HandleFunc("/user", handlers.UserRouter)
	r.HandleFunc("/comment", handlers.CommentRouter)
	r.HandleFunc("/debug/page", handlers.SendDebug)
	r.HandleFunc("/debug/contents", handlers.Debug)
	log.Fatal(srv.ListenAndServe())

}
