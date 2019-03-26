package main

import (
	"net/http"

	"zood.xyz/buster/l10n"
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
		"title":         "Zood",
		"cssPath":       "/css/home.css",
		"googlePlayURL": l10n.GooglePlayURL,
		"appStoreURL":   l10n.AppStoreURL,
	})
}

func locationAppHomeHandler(w http.ResponseWriter, r *http.Request) {
	rsrc.ExecuteTemplate("location-home.html", w, map[string]interface{}{
		"title":         "Zood Location | Zood",
		"cssPath":       "/css/location-home.css",
		"googlePlayURL": l10n.GooglePlayURL,
		"appStoreURL":   l10n.AppStoreURL,
		"ogImage":       "/images/zood-pixel-3-xl-720px.jpg",
	})
}

func mobileAppsPrivacyHandler(w http.ResponseWriter, r *http.Request) {
	rsrc.ExecuteTemplate("privacy-mobile-apps.html", w, map[string]interface{}{
		"title":   "Mobile Apps Privacy Policy | Zood",
		"cssPath": "/css/privacy-mobile-apps.css",
	})
}

func privacyHandler(w http.ResponseWriter, r *http.Request) {
	rsrc.ExecuteTemplate("privacy.html", w, map[string]interface{}{
		"title":   "Privacy Policy | Zood",
		"cssPath": "/css/privacy.css",
	})
}
