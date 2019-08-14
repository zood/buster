package main

import (
	"net/http"

	"zood.xyz/buster/constants"
	"zood.xyz/buster/l10n"
)

func aboutHandler(w http.ResponseWriter, r *http.Request) {
	rsrcs := resourcesFromContext(r.Context())
	rsrcs.ExecuteTemplate("about.html", w, map[string]interface{}{
		"title":       "About | Zood",
		"activeAbout": "active",
		"cssPath":     "/css/about.css",
	})
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	tag := l10n.MatchLanguage(r)
	rsrcs := resourcesFromContext(r.Context())
	rsrcs.ExecuteTemplate("home.html", w, map[string]interface{}{
		"title":              "Zood",
		"cssPath":            "/css/home.css",
		"activeHome":         "active",
		"googlePlayURL":      constants.GooglePlayURL,
		"appStoreURL":        constants.AppStoreURL,
		"companyDescription": l10n.String(tag, l10n.CompanyDescriptionMsg),
		"ogDescription":      l10n.String(tag, l10n.CompanyDescriptionMsg),
	})
}

func locationAppHomeHandler(w http.ResponseWriter, r *http.Request) {
	tag := l10n.MatchLanguage(r)
	rsrcs := resourcesFromContext(r.Context())
	rsrcs.ExecuteTemplate("location-home.html", w, map[string]interface{}{
		"title":             "Zood Location | Zood",
		"cssPath":           "/css/location-home.css",
		"googlePlayURL":     constants.GooglePlayURL,
		"appStoreURL":       constants.AppStoreURL,
		"activeLocation":    "active",
		"ogImage":           "/images/zood-location-icon-512.png",
		"zoodLocationBlurb": l10n.String(tag, l10n.ZoodLocationBlurbMsg),
		"ogDescription":     l10n.String(tag, l10n.ZoodLocationBlurbMsg),
	})
}

func mobileAppsPrivacyHandler(w http.ResponseWriter, r *http.Request) {
	rsrcs := resourcesFromContext(r.Context())
	rsrcs.ExecuteTemplate("privacy-mobile-apps.html", w, map[string]interface{}{
		"title":               "Mobile Apps Privacy Policy | Zood",
		"activePrivacyMobile": "active",
		"cssPath":             "/css/privacy-mobile-apps.css",
	})
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("page not found"))
}

func privacyHandler(w http.ResponseWriter, r *http.Request) {
	rsrcs := resourcesFromContext(r.Context())
	rsrcs.ExecuteTemplate("privacy.html", w, map[string]interface{}{
		"title":         "Privacy Policy | Zood",
		"activePrivacy": "active",
		"cssPath":       "/css/privacy.css",
	})
}
