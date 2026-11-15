package main

import (
	"context"
	"net/http"

	"zood.xyz/buster/email"
	"zood.xyz/buster/resources"
)

type busterMiddleware struct {
	rsrcs       *resources.Resources
	sendEmailer email.SendEmailer
}

// noCache sets "Cache-Control: no-cache" on responses from next so browsers
// will revalidate cached copies before each use.
func noCache(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-cache")
		next.ServeHTTP(w, r)
	})
}

func (m busterMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), contextResourcesKey, m.rsrcs)
		ctx = context.WithValue(ctx, contextSendEmailerKey, m.sendEmailer)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
