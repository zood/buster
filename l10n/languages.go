package l10n

import (
	"net/http"

	"golang.org/x/text/language"
)

var matcher = language.NewMatcher([]language.Tag{
	language.English,
})

// MatchLanguage returns the best match between languages supported
// by the server, and the languages specified in the request
func MatchLanguage(r *http.Request) language.Tag {
	t, _, _ := language.ParseAcceptLanguage(r.Header.Get("Accept-Language"))
	// ignore the error, because the matcher will select English if nil is passed in
	tag, _, _ := matcher.Match(t...)
	return tag
}
