package middleware

import (
	"log"
	"runtime/debug"

	"github.com/diamondburned/arikawa/v3/gateway"
)

func PanicRecoveryMiddleware[S any](e *gateway.InteractionCreateEvent, state *S, next ...Middleware[S]) error {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("PANIC: %+v\n%s", r, debug.Stack())
		}
	}()
	if len(next) > 0 {
		return next[0](e, state, next[1:]...)
	}
	return nil
}
