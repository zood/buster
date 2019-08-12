package l10n

import (
	"log"

	"golang.org/x/text/language"
)

// StringAsset is a key used to lookup a localized string
type StringAsset int

// Enumeration of all the localized strings in the app
const (
	CompanyDescriptionMsg StringAsset = iota
	ZoodLocationBlurbMsg
)

// String returns a localized string of asset for the specific
// acceptLangHdr
func String(tag language.Tag, asset StringAsset) string {
	s, ok := enStrings[asset]
	if !ok {
		log.Printf("WARNING: no string found for asset '%d'", asset)
		return "<undefined>"
	}

	return s
}
