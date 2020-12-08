package main

import (
	"net/http"

	"zood.xyz/buster/l10n"
)

func contactHandler(w http.ResponseWriter, r *http.Request) {
	tag := l10n.MatchLanguage(r)
	rsrcs := resourcesFromContext(r.Context())
	rsrcs.ExecuteTemplate("contact.html", w, map[string]interface{}{
		"title":               "Contact | Zood",
		"activeHeader":        "contact",
		"cssPath":             "/css/contact.css",
		"GetInTouchWithUs":    l10n.String(tag, l10n.GetInTouchWithUs),
		"GetInTouchWithUsMsg": l10n.String(tag, l10n.GetInTouchWithUsMsg),
	})
}
