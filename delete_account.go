package main

import (
	"net/http"

	"zood.xyz/buster/l10n"
)

// deleteAccountHandler renders the page that walks a user through deleting
// their Zood Location account. All of the work happens in the browser
// (see assets/js/delete-account.ts), because logging in requires the user's
// password, which must never reach our servers.
func deleteAccountHandler(w http.ResponseWriter, r *http.Request) {
	tag := l10n.MatchLanguage(r)
	rsrcs := resourcesFromContext(r.Context())
	rsrcs.ExecuteTemplate("delete-account.html", w, map[string]interface{}{
		"title":                         "Delete Account | Zood",
		"activeHeader":                  "",
		"cssPath":                       "/css/delete-account.css",
		"DeleteYourAccount":             l10n.String(tag, l10n.DeleteYourAccount),
		"DeleteAccountConfirmationMsg":  l10n.String(tag, l10n.DeleteAccountConfirmationMsg),
		"BeginAccountDeletion":          l10n.String(tag, l10n.BeginAccountDeletion),
		"LogInToContinueMsg":            l10n.String(tag, l10n.LogInToContinueMsg),
		"Username":                      l10n.String(tag, l10n.Username),
		"Password":                      l10n.String(tag, l10n.Password),
		"LogIn":                         l10n.String(tag, l10n.LogIn),
		"AccountDeletionIsPermanentMsg": l10n.String(tag, l10n.AccountDeletionIsPermanentMsg),
		"IUnderstandDeleteMyAccount":    l10n.String(tag, l10n.IUnderstandDeleteMyAccount),
		"DeletionComplete":              l10n.String(tag, l10n.DeletionComplete),
		"IncorrectUsernameOrPassword":   l10n.String(tag, l10n.IncorrectUsernameOrPasswordMsg),
		"InvalidUsernameMsg":            l10n.String(tag, l10n.InvalidUsernameMsg),
		"NetworkErrorMsg":               l10n.String(tag, l10n.NetworkErrorMsg),
		"UnexpectedErrorMsg":            l10n.String(tag, l10n.UnexpectedErrorMsg),
		"WorkingEllipsis":               l10n.String(tag, l10n.WorkingEllipsis),
	})
}
