package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"zood.xyz/buster/l10n"
)

func contactErrorHandler(w http.ResponseWriter, r *http.Request) {
	tag := l10n.MatchLanguage(r)
	rsrcs := resourcesFromContext(r.Context())
	cause := r.URL.Query().Get("error")
	var errMsg string
	switch cause {
	case "missing-name":
		errMsg = l10n.String(tag, l10n.ContactFormErrorMissingNameMsg)
	case "missing-email":
		errMsg = l10n.String(tag, l10n.ContactFormErrorMissingEmailMsg)
	case "missing-message":
		errMsg = l10n.String(tag, l10n.ContactFormErrorMissingMessageMsg)
	default:
		errMsg = l10n.String(tag, l10n.ContactFormErrorUnknownMsg)
	}
	rsrcs.ExecuteTemplate("contact-error.html", w, map[string]interface{}{
		"errorMessage": errMsg,
	})
}

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

func contactSuccessHandler(w http.ResponseWriter, r *http.Request) {
	// tag := l10n.MatchLanguage(r)
	rsrcs := resourcesFromContext(r.Context())
	rsrcs.ExecuteTemplate("contact-success.html", w, map[string]interface{}{})
}

func parseContactForm(r *http.Request) (name, email, msg *string, err error) {
	nameVal := r.FormValue("name")
	emailVal := r.FormValue("email")
	msgVal := r.FormValue("message")

	if len(nameVal) != 0 {
		nameVal = strings.TrimSpace(nameVal)
		if nameVal != "" {
			name = &nameVal
		}
	}
	if len(emailVal) != 0 {
		emailVal = strings.TrimSpace(emailVal)
		if emailVal != "" {
			email = &emailVal
		}
	}
	if len(msgVal) != 0 {
		msgVal = strings.TrimSpace(msgVal)
		if msgVal != "" {
			msg = &msgVal
		}
	}

	return
}

func submitContactFormHandler(w http.ResponseWriter, r *http.Request) {
	name, email, msg, err := parseContactForm(r)

	if err != nil {
		log.Printf("Error parsing contact form: %v", err)
		http.Redirect(w, r, "/contact-error", http.StatusSeeOther)
		return
	}

	if name == nil {
		http.Redirect(w, r, "/contact-error?error=missing-name", http.StatusSeeOther)
		return
	}
	if email == nil {
		http.Redirect(w, r, "/contact-error?error=missing-email", http.StatusSeeOther)
		return
	}
	if msg == nil {
		http.Redirect(w, r, "/contact-error?error=missing-message", http.StatusSeeOther)
		return
	}

	// success!
	http.Redirect(w, r, "/contact-success", http.StatusSeeOther)

	emailer := sendEmailer(r.Context())
	go func() {
		emailBody := fmt.Sprintf(`Message from '%s' (%s):\n\n%s`, *name, *email, *msg)
		err := emailer.SendEmail("buster-contact-form@notifications.zood.xyz", "arash@zood.xyz", "Buster contact form submission", emailBody, nil)
		if err != nil {
			logInternalError(err, emailer)
		}
	}()
}
