package middleware

import (
	"github.com/diamondburned/arikawa/v3/gateway"
)

type Middleware[S any] func(e *gateway.InteractionCreateEvent, state *S, next ...Middleware[S]) error
