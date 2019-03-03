package main

import (
	"net/http"

	"zood.xyz/buster/rsrc"
)

func aboutHandler(w http.ResponseWriter, r *http.Request) {
	rsrc.ExecuteTemplate("about.html", w, map[string]interface{}{
		"title":   "About | Zood",
		"cssPath": "/css/about.css",
	})
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	rsrc.ExecuteTemplate("home.html", w, map[string]interface{}{
		"title":   "Zood",
		"cssPath": "/css/home.css",
	})
}

func locationAppHomeHandler(w http.ResponseWriter, r *http.Request) {
	rsrc.ExecuteTemplate("location-home.html", w, map[string]interface{}{
		"title":   "Zood Location | Zood",
		"cssPath": "/css/location-home.css",
	})
}

func privacyHandler(w http.ResponseWriter, r *http.Request) {
	rsrc.ExecuteTemplate("privacy.html", w, map[string]interface{}{
		"title":   "Privacy Policy | Zood",
		"cssPath": "/css/privacy.css",
	})
}
