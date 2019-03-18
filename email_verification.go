package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/pkg/errors"
	"zood.xyz/buster/oscar"
	"zood.xyz/buster/rsrc"
)

func disavowEmailHandler(w http.ResponseWriter, r *http.Request) {
	tmplName := "disavow-email.html"
	data := map[string]interface{}{
		"title":   "Disavow Email | Zood",
		"cssPath": "/css/email-verification.css",
	}

	token := r.URL.Query().Get("t")
	token = strings.TrimSpace(token)
	if token == "" {
		data["line1"] = "The email token is missing."
		data["line2"] = "Double check the URL then try again."
		rsrc.ExecuteTemplate(tmplName, w, data)
		return
	}

	endpoint := fmt.Sprintf("https://api.zood.xyz/1/email-verifications/%s", token)
	req, err := http.NewRequest(http.MethodDelete, endpoint, nil)
	if err != nil {
		internalError(w, err)
		return
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		internalError(w, err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		buf, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			err = errors.Wrapf(err, "problem reading error response while disavowing token '%s'", token)
			internalError(w, err)
			return
		}
		err = errors.Errorf("problem disavowing token '%s':\noscar responded with %d: %s", token, resp.StatusCode, buf)
		internalError(w, err)
		return
	}

	data["line1"] = "We've removed your email address from our system."
	data["line2"] = "Sorry for the inconvenience."
	rsrc.ExecuteTemplate(tmplName, w, data)
}

func verifyEmailHandler(w http.ResponseWriter, r *http.Request) {
	tmplName := "verify-email.html"
	data := map[string]interface{}{
		"title":   "Verify Email | Zood",
		"cssPath": "/css/email-verification.css",
	}

	token := r.URL.Query().Get("t")
	token = strings.TrimSpace(token)
	if token == "" {
		data["line1"] = "The email token is missing."
		data["line2"] = "Double check the URL then try again."
		rsrc.ExecuteTemplate(tmplName, w, data)
		return
	}

	postBody := struct {
		Token string `json:"token"`
	}{Token: token}
	postData, _ := json.Marshal(postBody)
	endpoint := "https://api.zood.xyz/1/email-verifications"
	req, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewReader(postData))
	if err != nil {
		internalError(w, err)
		return
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		internalError(w, err)
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
			internalError(w, err)
			return
		}

		if errBody.Code == oscar.ErrorMissingVerificationToken {
			data["line1"] = "Hmm… that URL doesn't work."
			data["line2"] = "Did you already verify your email? If not, double check the URL then try again."
			rsrc.ExecuteTemplate(tmplName, w, data)
			return
		}

		// unexpected api response. probably a 500 from oscar
		err = errors.Errorf("Unexpected response from oscar while attempting to verify token '%s': %d - %s",
			token,
			errBody.Code,
			errBody.Msg)
		internalError(w, err)
		return
	}

	data["line1"] = "Your email has been verified!"
	data["line2"] = "It's safe to close this window."
	rsrc.ExecuteTemplate(tmplName, w, data)
}
