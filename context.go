package main

import (
	"context"

	"zood.xyz/buster/resources"
)

type contextKey string

const (
	contextResourcesKey = contextKey("resources")
)

func resourcesFromContext(ctx context.Context) *resources.Resources {
	return ctx.Value(contextResourcesKey).(*resources.Resources)
}
