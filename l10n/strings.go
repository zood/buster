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
	ThePrivacyCompany
	ShareYourLocation
	LocationDescriptionMsg
	SafeAndSecure
	SafeDescriptionMsg
	TrustedDescriptionMsg
	AddOnlyTrustedFriends
	ZoodLocationBlurbMsg
	ItsAboutPrivacy
	Really
	AboutPrivacyMsg
	ZoodIsDifferent
	AboutDifferentMsg
	DidWeMentionTheEncryptionInterrogative
	AboutEncryptionMsg
	LotsOfServicesSecureQuestionMsg
	LotsOfServicesSecureAnswerMsg
	WhenZoodLocationReleasedQuestionMsg
	WhenZoodLocationReleasedAnswerMsg
	IfYouDontSurveilQuestionMsg
	IfYouDontSurveilAnswerMsg
	HowDoISubmitQuestionMsg
	HowDoISubmitAnswerMsg
	GetInTouchWithUs
	GetInTouchWithUsMsg
	ContactFormErrorMissingNameMsg
	ContactFormErrorMissingEmailMsg
	ContactFormErrorMissingMessageMsg
	ContactFormErrorUnknownMsg
	MobileAppsPrivacyPolicy
	MobileAppsPrivacyPolicyMsg
	StuffWeKnowAboutYou
	StuffWeKnowAboutYouMsg
	StuffWeBackupForYou
	StuffWeBackupForYouMsg
)

// String returns a localized string or template.HTML matching the language
// represented by tag
func String(tag language.Tag, asset StringAsset) interface{} {
	s, ok := enStrings[asset]
	if !ok {
		log.Printf("WARNING: no entry found for string asset '%d'", asset)
		return "<undefined>"
	}

	return s
}
