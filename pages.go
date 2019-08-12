package main

import (
	"net/http"

	"zood.xyz/buster/constants"
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
	tag := l10n.MatchLanguage(r)
	rsrc.ExecuteTemplate("home.html", w, map[string]interface{}{
		"title":              "Zood",
		"cssPath":            "/css/home.css",
		"googlePlayURL":      constants.GooglePlayURL,
		"appStoreURL":        constants.AppStoreURL,
		"companyDescription": l10n.String(tag, l10n.CompanyDescriptionMsg),
		"ogDescription":      l10n.String(tag, l10n.CompanyDescriptionMsg),
	})
}

func locationAppHomeHandler(w http.ResponseWriter, r *http.Request) {
	tag := l10n.MatchLanguage(r)
	rsrc.ExecuteTemplate("location-home.html", w, map[string]interface{}{
		"title":             "Zood Location | Zood",
		"cssPath":           "/css/location-home.css",
		"googlePlayURL":     constants.GooglePlayURL,
		"appStoreURL":       constants.AppStoreURL,
		"ogImage":           "/images/zood-location-icon-512.png",
		"zoodLocationBlurb": l10n.String(tag, l10n.ZoodLocationBlurbMsg),
		"ogDescription":     l10n.String(tag, l10n.ZoodLocationBlurbMsg),
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
