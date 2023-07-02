package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"zood.xyz/buster/oscar"
)

func disavowEmailHandler(w http.ResponseWriter, r *http.Request) {
	tmplName := "disavow-email.html"
	data := map[string]interface{}{
		"title":        "Disavow Email | Zood",
		"cssPath":      "/css/email-verification.css",
		"activeHeader": "",
	}

	rsrcs := resourcesFromContext(r.Context())
	token := r.URL.Query().Get("t")
	token = strings.TrimSpace(token)
	if token == "" {
		data["line1"] = "The email token is missing."
		data["line2"] = "Double check the URL then try again."
		rsrcs.ExecuteTemplate(tmplName, w, data)
		return
	}

	endpoint := fmt.Sprintf("https://api.zood.xyz/1/email-verifications/%s", token)
	req, err := http.NewRequest(http.MethodDelete, endpoint, nil)
	if err != nil {
		internalError(w, r, err)
		return
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		internalError(w, r, err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		buf, err := io.ReadAll(resp.Body)
		if err != nil {
			err = fmt.Errorf("problem reading error response while disavowing token '%s': %v", token, err)
			internalError(w, r, err)
			return
		}
		err = fmt.Errorf("problem disavowing token '%s':\noscar responded with %d: %s", token, resp.StatusCode, buf)
		internalError(w, r, err)
		return
	}

	data["line1"] = "We've removed your email address from our system."
	data["line2"] = "Sorry for the inconvenience."
	rsrcs.ExecuteTemplate(tmplName, w, data)
}

func verifyEmailHandler(w http.ResponseWriter, r *http.Request) {
	tmplName := "verify-email.html"
	data := map[string]interface{}{
		"title":        "Verify Email | Zood",
		"cssPath":      "/css/email-verification.css",
		"activeHeader": "",
	}

	rsrcs := resourcesFromContext(r.Context())
	token := r.URL.Query().Get("t")
	token = strings.TrimSpace(token)
	if token == "" {
		data["line1"] = "The email token is missing."
		data["line2"] = "Double check the URL then try again."
		rsrcs.ExecuteTemplate(tmplName, w, data)
		return
	}

	postBody := struct {
		Token string `json:"token"`
	}{Token: token}
	postData, _ := json.Marshal(postBody)
	endpoint := "https://api.zood.xyz/1/email-verifications"
	req, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewReader(postData))
	if err != nil {
		internalError(w, r, err)
		return
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		internalError(w, r, err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		errBody := struct {
			Msg  string `json:"error_message"`
			Code int    `json:"error_code"`
		}{}
		err := json.NewDecoder(resp.Body).Decode(&errBody)
		if err != nil {
			internalError(w, r, err)
			return
		}

		if errBody.Code == oscar.ErrorMissingVerificationToken {
			data["line1"] = "Hmm… that token doesn't seem to be valid."
			data["line2"] = "Did you already verify your email? If not, double check the URL then try again."
			rsrcs.ExecuteTemplate(tmplName, w, data)
			return
		}

		// unexpected api response. probably a 500 from oscar
		err = fmt.Errorf("unexpected response from oscar while attempting to verify token '%s': %d - %s",
			token,
			errBody.Code,
			errBody.Msg)
		internalError(w, r, err)
		return
	}

	data["line1"] = "Your email has been verified!"
	data["line2"] = "It's safe to close this window."

	rsrcs.ExecuteTemplate(tmplName, w, data)
}
