package main

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"runtime"

	"zood.xyz/buster/email"
)

func internalError(w http.ResponseWriter, r *http.Request, err error) {
	data := map[string]interface{}{
		"title":   "Server Error | Zood",
		"cssPath": "/css/server-error.css",
	}

	rsrcs := resourcesFromContext(r.Context())
	emailer := sendEmailer(r.Context())
	rsrcs.ExecuteTemplateCode("server-error.html", w, data, http.StatusInternalServerError)
	if err != nil {
		_, file, line, ok := runtime.Caller(1)
		if !ok {
			file = "???"
			line = 0
		}
		file = filepath.Base(file)
		log.Printf("%s:%d %v", file, line, err)
		mailMsg := fmt.Sprintf("Buster Error\n%s:%d\n%v", file, line, err.Error())
		mgErr := emailer.SendEmail("buster-errors@notifications.zood.xyz",
			"arash@zood.xyz",
			fmt.Sprintf("Buster Error: %s:%d", file, line),
			mailMsg, nil)
		if mgErr != nil {
			log.Printf("Failed to email internal error: %v", mgErr)
		}
	}
}

func logInternalError(err error, emailer email.SendEmailer) {
	if err != nil {
		_, file, line, ok := runtime.Caller(1)
		if !ok {
			file = "???"
			line = 0
		}
		file = filepath.Base(file)
		log.Printf("%s:%d %v", file, line, err)
		mailMsg := fmt.Sprintf("Buster Error\n%s:%d\n%v", file, line, err.Error())
		mgErr := emailer.SendEmail("buster-errors@notifications.zood.xyz",
			"arash@zood.xyz",
			fmt.Sprintf("Buster Error: %s:%d", file, line),
			mailMsg, nil)
		if mgErr != nil {
			log.Printf("Failed to email internal error: %v", mgErr)
		}
	}
}
