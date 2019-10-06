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
		"title":                    "About | Zood",
		"activeHeader":             "about",
		"ItsAboutPrivacy":          l10n.String(tag, l10n.ItsAboutPrivacy),
		"Really":                   l10n.String(tag, l10n.Really),
		"ZoodIsDifferent":          l10n.String(tag, l10n.ZoodIsDifferent),
		"DidWeMentionEncryption":   l10n.String(tag, l10n.DidWeMentionTheEncryptionInterrogative),
		"AboutPrivacyMsg":          l10n.String(tag, l10n.AboutPrivacyMsg),
		"AboutDifferentMsg":        l10n.String(tag, l10n.AboutDifferentMsg),
		"AboutEncryptionMsg":       l10n.String(tag, l10n.AboutEncryptionMsg),
		"IfYouDontSurveilQuestion": l10n.String(tag, l10n.IfYouDontSurveilQuestionMsg),
		"IfYouDontSurveilAnswer":   l10n.String(tag, l10n.IfYouDontSurveilAnswerMsg),
		"HowDoISubmitQuestion":     l10n.String(tag, l10n.HowDoISubmitQuestionMsg),
		"HowDoISubmitAnswer":       l10n.String(tag, l10n.HowDoISubmitAnswerMsg),
		"cssPath":                  "/css/about.css",
	})
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	tag := l10n.MatchLanguage(r)
	rsrcs := resourcesFromContext(r.Context())
	rsrcs.ExecuteTemplate("home.html", w, map[string]interface{}{
		"title":                        "Zood",
		"cssPath":                      "/css/home.css",
		"activeHeader":                 "home",
		"googlePlayURL":                constants.GooglePlayURL,
		"appStoreURL":                  constants.AppStoreURL,
		"companyDescription":           l10n.String(tag, l10n.CompanyDescriptionMsg),
		"ThePrivacyCompany":            l10n.String(tag, l10n.ThePrivacyCompany),
		"ShareYourLocation":            l10n.String(tag, l10n.ShareYourLocation),
		"LocationDescriptionMsg":       l10n.String(tag, l10n.LocationDescriptionMsg),
		"SafeAndSecure":                l10n.String(tag, l10n.SafeAndSecure),
		"SafeDescriptionMsg":           l10n.String(tag, l10n.SafeDescriptionMsg),
		"AddOnlyTrustedFriends":        l10n.String(tag, l10n.AddOnlyTrustedFriends),
		"TrustedDescriptionMsg":        l10n.String(tag, l10n.TrustedDescriptionMsg),
		"ZoodLocationBlurbMsg":         l10n.String(tag, l10n.ZoodLocationBlurbMsg),
		"ogDescription":                l10n.String(tag, l10n.CompanyDescriptionMsg),
		"LotsofServicesQuestion":       l10n.String(tag, l10n.LotsOfServicesSecureQuestionMsg),
		"LotsofServicesAnswer":         l10n.String(tag, l10n.LotsOfServicesSecureAnswerMsg),
		"WhenWillZoodLocationQuestion": l10n.String(tag, l10n.WhenZoodLocationReleasedQuestionMsg),
		"WhenWillZoodLocationAnswer":   l10n.String(tag, l10n.WhenZoodLocationReleasedAnswerMsg),
	})
}

func locationAppHomeHandler(w http.ResponseWriter, r *http.Request) {
	tag := l10n.MatchLanguage(r)
	rsrcs := resourcesFromContext(r.Context())
	rsrcs.ExecuteTemplate("location-home.html", w, map[string]interface{}{
		"title":             "Zood Location | Zood",
		"cssPath":           "/css/location-home.css",
		"googlePlayURL":     constants.GooglePlayURL,
		"appStoreURL":       constants.AppStoreURL,
		"ogImage":           "/images/zood-location-icon-512.png",
		"zoodLocationBlurb": l10n.String(tag, l10n.ZoodLocationBlurbMsg),
		"ogDescription":     l10n.String(tag, l10n.ZoodLocationBlurbMsg),
	})
}

func mobileAppsPrivacyHandler(w http.ResponseWriter, r *http.Request) {
	rsrcs := resourcesFromContext(r.Context())
	rsrcs.ExecuteTemplate("privacy-mobile-apps.html", w, map[string]interface{}{
		"title":                      "Mobile Apps Privacy Policy | Zood",
		"activeHeader":               "privacymobile",
		"cssPath":                    "/css/privacy-mobile-apps.css",
		"MobileAppsPrivacyPolicy":    "Mobile Apps Privacy Policy",
		"MobileAppsPrivacyPolicyMsg": "Everything we build aims to increase, or at the very least preserve, your privacy. When you use Zood Location, the information that you share with your family+friends is shared with them only using end-to-end encryption. That means we don't know anything about your location, so we can't spy on your or sell your data, and nobody can compel us to reveal your location either.",
		"StuffWeKnowAboutYou":        "Stuff we know about you",
		"StuffWeKnowAboutYouMsg":     "Email Address: Upon registering an account, you may optionally provide your email address. When you do so, we send you an email with a verification link to make sure we have the correct address. The email address will only be used to contact you with important information about your account (like resetting your account if you forgot your password). We won't send you any spam, sell your email address, or share it with any 3rd parties.",
		"StuffWeBackupForYou":        "Stuff we backup for you",
		"StuffWeBackupForYouMsg":     "To make it easy for you to switch phones, Zood Location sends an encrypted backup of your database of friends to our server. It's encrypted with a key derived from your password (which we also don't know), so to us it's just a blob of random data that we hold just for you. That's it! No trickery here. Just a simple service that lets you (and only you) know where your loved ones are all the time.",
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
