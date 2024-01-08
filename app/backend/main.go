package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/TravisRoad/blog-edit/global"
	"github.com/TravisRoad/blog-edit/internal/handler"
	"github.com/TravisRoad/blog-edit/internal/middleware"
	"github.com/gorilla/mux"
)

func main() {
	if err := global.LoadConfig(); err != nil {
		log.Fatal(err)
	}

	r := mux.NewRouter()

	r.Use(middleware.Recover, middleware.Logging)

	r.HandleFunc("/v1/file", handler.GetFileList).Methods("GET")
	r.HandleFunc("/v1/file/{filename}", handler.GetFile).Methods("GET")
	r.HandleFunc("/v1/file/{filename}", handler.UpdateFile).Methods("POST")
	r.HandleFunc("/v1/apply", handler.Apply).Methods("POST")
	r.HandleFunc("/v1/sync", handler.Sync).Methods("POST")

	srv := http.Server{
		Addr:         global.Config.Addr,
		Handler:      r,
		ReadTimeout:  time.Second * 15,
		WriteTimeout: time.Second * 15,
		IdleTimeout:  time.Second * 60,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	<-c

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Println(err)
		os.Exit(1)
	}
	log.Println("shutting down")
	os.Exit(0)
}
