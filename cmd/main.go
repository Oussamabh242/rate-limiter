package main

import (
	"fmt"
	"log"
	"net/http"
	"github.com/Oussamabh242/rate-limiter/pkg/middleware"
	"time"

	"github.com/gorilla/mux"
)



func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello there")
}

func main() {
	r := mux.NewRouter()
  t := middleware.InitThing(1*time.Minute,30*time.Second ,2 ,time.Minute*2)
	r.Use(t.Middleware)
	r.HandleFunc("/", handler).Methods("GET")
  r.Use()
	srv := http.Server{
		Handler: r,
		Addr:    "0.0.0.0:3000",
	}
	log.Fatal(srv.ListenAndServe())
}
