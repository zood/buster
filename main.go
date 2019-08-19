package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"time"

	"github.com/gorilla/mux"
	"zood.xyz/buster/resources"
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	resPath := flag.String("resources", "", "Path to resources directory")
	devMode := flag.Bool("dev", false, "true makes the server reload templates on every request. default is false.")
	port := flag.Int("port", 1313, "Port to listen on")
	flag.Parse()

	if *resPath == "" {
		log.Fatal("Resources path is empty")
	}

	rsrcs, err := resources.New(*resPath)
	if err != nil {
		log.Fatalf("Failed to load resources: %v", err)
	}
	rsrcs.DevMode = *devMode

	r := mux.NewRouter().StrictSlash(true)
	r.PathPrefix("/css/").Handler(http.StripPrefix("/css/", http.FileServer(http.Dir(filepath.Join(*resPath, "css")))))
	r.PathPrefix("/images/").Handler(http.StripPrefix("/images/", http.FileServer(http.Dir(filepath.Join(*resPath, "images")))))
	r.PathPrefix("/fonts/").Handler(http.StripPrefix("/fonts/", http.FileServer(http.Dir(filepath.Join(*resPath, "fonts")))))

	r.HandleFunc("/", homeHandler).Methods(http.MethodGet)
	r.HandleFunc("/privacy", privacyHandler).Methods(http.MethodGet)
	r.HandleFunc("/privacy/mobile-apps", mobileAppsPrivacyHandler).Methods(http.MethodGet)
	r.HandleFunc("/about", aboutHandler).Methods(http.MethodGet)
	r.HandleFunc("/products/location", locationAppHomeHandler).Methods(http.MethodGet)

	// blog
	r.HandleFunc("/blog", blogHomeHandler).Methods(http.MethodGet)
	r.HandleFunc("/blog/archive", blogArchivesHandler).Methods(http.MethodGet)
	r.HandleFunc("/blog/{id:[0-9]+}", blogPostHandler).Methods(http.MethodGet)
	r.HandleFunc("/blog/{id:[0-9]+}/{slug}", blogPostHandler).Methods(http.MethodGet)

	r.HandleFunc("/verify-email", verifyEmailHandler).Methods(http.MethodGet)
	r.HandleFunc("/disavow-email", disavowEmailHandler).Methods(http.MethodGet)

	r.NotFoundHandler = http.HandlerFunc(notFoundHandler)
	r.MethodNotAllowedHandler = http.HandlerFunc(notFoundHandler)

	r.Use(busterMiddleware{rsrcs: rsrcs}.Middleware)

	var hostAddress string
	if *devMode {
		hostAddress = fmt.Sprintf(":%d", *port)
	} else {
		hostAddress = fmt.Sprintf("127.0.0.1:%d", *port)
	}
	server := http.Server{
		Addr:         hostAddress,
		Handler:      r,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	log.Printf("Starting app on port %d…", *port)
	log.Fatal(server.ListenAndServe())
}
