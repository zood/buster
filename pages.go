package main

import (
	"net/http"

	"zood.xyz/buster/constants"
	"zood.xyz/buster/l10n"
)

func aboutHandler(w http.ResponseWriter, r *http.Request) {
	tag := l10n.MatchLanguage(r)
	rsrcs := resourcesFromContext(r.Context())
	rsrcs.ExecuteTemplate("about.html", w, map[string]interface{}{
		"title":                          "About | Zood",
		"activeHeader":                   "about",
		"ItsAboutPrivacy":                l10n.String(tag, l10n.ItsAboutPrivacy),
		"Really":                         l10n.String(tag, l10n.Really),
		"ZoodIsDifferent":                l10n.String(tag, l10n.ZoodIsDifferent),
		"DidWeMentionEncryption":         l10n.String(tag, l10n.DidWeMentionTheEncryptionInterrogative),
		"AboutPrivacyMsg":                l10n.String(tag, l10n.AboutPrivacyMsg),
		"AboutDifferentMsg":              l10n.String(tag, l10n.AboutDifferentMsg),
		"AboutEncryptionMsg":             l10n.String(tag, l10n.AboutEncryptionMsg),
		"LotsofServicesQuestion":         l10n.String(tag, l10n.LotsOfServicesSecureQuestionMsg),
		"LotsofServicesAnswer":           l10n.String(tag, l10n.LotsOfServicesSecureAnswerMsg),
		"WhenZoodLocationIosQuestionMsg": l10n.String(tag, l10n.WhenZoodLocationIosQuestionMsg),
		"WhenZoodLocationIosAnswerMsg":   l10n.String(tag, l10n.WhenZoodLocationIosAnswerMsg),
		// "WhenWillZoodLocationQuestion":   l10n.String(tag, l10n.WhenZoodLocationReleasedQuestionMsg),
		// "WhenWillZoodLocationAnswer":     l10n.String(tag, l10n.WhenZoodLocationReleasedAnswerMsg),
		"IfYouDontSurveilQuestion":    l10n.String(tag, l10n.IfYouDontSurveilQuestionMsg),
		"IfYouDontSurveilAnswer":      l10n.String(tag, l10n.IfYouDontSurveilAnswerMsg),
		"HowDoISubmitQuestion":        l10n.String(tag, l10n.HowDoISubmitQuestionMsg),
		"HowDoISubmitAnswer":          l10n.String(tag, l10n.HowDoISubmitAnswerMsg),
		"WhosBehindZoodInterrogative": l10n.String(tag, l10n.WhosBehindZoodInterrogative),
		"cssPath":                     "/css/about.css",
	})
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	tag := l10n.MatchLanguage(r)
	rsrcs := resourcesFromContext(r.Context())
	rsrcs.ExecuteTemplate("home.html", w, map[string]interface{}{
		"title":                  "Zood",
		"cssPath":                "/css/home.css",
		"activeHeader":           "home",
		"googlePlayURL":          constants.GooglePlayURL,
		"appStoreURL":            constants.AppStoreURL,
		"companyDescription":     l10n.String(tag, l10n.CompanyDescriptionMsg),
		"ThePrivacyCompany":      l10n.String(tag, l10n.ThePrivacyCompany),
		"ShareYourLocation":      l10n.String(tag, l10n.ShareYourLocation),
		"LocationDescriptionMsg": l10n.String(tag, l10n.LocationDescriptionMsg),
		"SafeAndSecure":          l10n.String(tag, l10n.SafeAndSecure),
		"SafeDescriptionMsg":     l10n.String(tag, l10n.SafeDescriptionMsg),
		"OnlyAddTrustedFriends":  l10n.String(tag, l10n.OnlyAddTrustedFriends),
		"TrustedDescriptionMsg":  l10n.String(tag, l10n.TrustedDescriptionMsg),
		"ZoodLocationBlurbMsg":   l10n.String(tag, l10n.ZoodLocationBlurbMsg),
		"ogDescription":          l10n.String(tag, l10n.CompanyDescriptionMsg),
	})
}

func mobileAppsPrivacyHandler(w http.ResponseWriter, r *http.Request) {
	tag := l10n.MatchLanguage(r)
	rsrcs := resourcesFromContext(r.Context())
	rsrcs.ExecuteTemplate("privacy-mobile-apps.html", w, map[string]interface{}{
		"title":                      "Mobile Apps Privacy Policy | Zood",
		"activeHeader":               "privacymobile",
		"cssPath":                    "/css/privacy-mobile-apps.css",
		"MobileAppsPrivacyPolicy":    l10n.String(tag, l10n.MobileAppsPrivacyPolicy),
		"MobileAppsPrivacyPolicyMsg": l10n.String(tag, l10n.MobileAppsPrivacyPolicyMsg),
		"StuffWeKnowAboutYou":        l10n.String(tag, l10n.StuffWeKnowAboutYou),
		"StuffWeKnowAboutYouMsg":     l10n.Markdown(tag, l10n.StuffWeKnowAboutYouMsg),
		"StuffWeBackupForYou":        l10n.String(tag, l10n.StuffWeBackupForYou),
		"StuffWeBackupForYouMsg":     l10n.Markdown(tag, l10n.StuffWeBackupForYouMsg),
	})
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("page not found"))
}

func privacyHandler(w http.ResponseWriter, r *http.Request) {
	rsrcs := resourcesFromContext(r.Context())
	rsrcs.ExecuteTemplate("privacy.html", w, map[string]interface{}{
		"title":        "Privacy Policy | Zood",
		"activeHeader": "privacy",
		"cssPath":      "/css/privacy.css",
	})
}
