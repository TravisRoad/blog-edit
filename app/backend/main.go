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

	r.Use(middleware.Recover, middleware.Logging, middleware.Auth)

	r.HandleFunc("/v1/file", handler.GetFileList).Methods("GET")
	r.HandleFunc("/v1/file/{filename}", handler.GetFile).Methods("GET")
	r.HandleFunc("/v1/file/{filename}", handler.UpdateFile).Methods("PUT")
	r.HandleFunc("/v1/file/{filename}", handler.Apply).Methods("POST")
	r.HandleFunc("/v1/sync", handler.Sync).Methods("POST")
	r.HandleFunc("/v1/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}).Methods("GET")
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}).Methods("GET")

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
