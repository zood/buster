package main

import (
	"fmt"
	"io"
	"log"
	"mime"
	"mime/multipart"
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
	mediaType, params, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil {
		return nil, nil, nil, fmt.Errorf("while parsing media type: %w", err)
	}

	if !strings.HasPrefix(mediaType, "multipart/") {
		return nil, nil, nil, fmt.Errorf("not a multipart form - is '%s'", mediaType)
	}

	var nameData []byte
	var emailData []byte
	var msgData []byte

	limitReader := io.LimitReader(r.Body, 64*1024)
	mr := multipart.NewReader(limitReader, params["boundary"])
	for {
		p, err := mr.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, nil, nil, fmt.Errorf("while parsing form: %w", err)
		}
		switch p.FormName() {
		case "name":
			p.Read(nameData)
		case "email":
			p.Read(emailData)
		case "message":
			p.Read(msgData)
		default:
			continue
		}
	}

	if len(nameData) != 0 {
		tmp := string(nameData)
		tmp = strings.TrimSpace(tmp)
		if tmp != "" {
			name = &tmp
		}
	}
	if len(emailData) != 0 {
		tmp := string(emailData)
		tmp = strings.TrimSpace(tmp)
		if tmp != "" {
			email = &tmp
		}
	}
	if len(msgData) != 0 {
		tmp := string(msgData)
		tmp = strings.TrimSpace(tmp)
		if tmp != "" {
			msg = &tmp
		}
	}

	err = nil
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
