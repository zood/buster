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

func (m busterMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), contextResourcesKey, m.rsrcs)
		ctx = context.WithValue(ctx, contextSendEmailerKey, m.sendEmailer)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
