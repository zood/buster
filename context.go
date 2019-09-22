package main

import (
	"context"

	"zood.xyz/buster/email"
	"zood.xyz/buster/resources"
)

type contextKey string

const (
	contextResourcesKey   = contextKey("resources")
	contextSendEmailerKey = contextKey("send_emailer")
)

func resourcesFromContext(ctx context.Context) *resources.Resources {
	return ctx.Value(contextResourcesKey).(*resources.Resources)
}

func sendEmailer(ctx context.Context) email.SendEmailer {
	return ctx.Value(contextSendEmailerKey).(email.SendEmailer)
}
