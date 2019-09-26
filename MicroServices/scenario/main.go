package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/matscus/Hamster/MicroServices/scenario/handlers"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/matscus/Hamster/MicroServices/scenario/scn"
	"github.com/matscus/Hamster/Package/Middleware/middleware"
)

var (
	pemPath    string
	keyPath    string
	proto      string
	listenport string
)

func init() {
	go func() {
		for {
			scn.InitData()
			time.Sleep(1 * time.Minute)
		}
	}()
}

func main() {
	flag.StringVar(&pemPath, "pem", os.Getenv("SERVERREM"), "path to pem file")
	flag.StringVar(&keyPath, "key", os.Getenv("SERVERKEY"), "path to key file")
	flag.StringVar(&listenport, "port", ":10004", "port to Listen")
	flag.StringVar(&proto, "proto", "https", "http or https")
	flag.Parse()
	r := mux.NewRouter()
	r.HandleFunc("/api/v1/scenario/start", middleware.Middleware(handlers.StartScenario)).Methods("POST", "OPTIONS")
	r.HandleFunc("/api/v1/scenario/stop", middleware.Middleware(handlers.StopScenario)).Methods("POST", "OPTIONS")
	r.HandleFunc("/api/v1/scenario/new", middleware.Middleware(handlers.NewScenario)).Methods("POST", "OPTIONS")
	r.HandleFunc("/api/v1/scenario", middleware.Middleware(handlers.GetData)).Methods("GET", "OPTIONS").Queries("project", "{project}")
	r.HandleFunc("/api/v1/scenario", middleware.Middleware(handlers.UpdateData)).Methods("PUT", "OPTIONS")
	r.HandleFunc("/api/v1/scenario/lastparams", middleware.Middleware(handlers.GetLastParams)).Methods("GET", "OPTIONS").Queries("name", "{name}")
	r.HandleFunc("/api/v1/scenario/ws", handlers.Ws)
	http.Handle("/api/v1/", r)
	log.Println("ListenAndServe: " + listenport)
	err := http.ListenAndServeTLS(listenport, pemPath, keyPath, context.ClearHandler(http.DefaultServeMux))
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}