package middleware

import (
	"log"
	"time"

	"github.com/diamondburned/arikawa/v3/gateway"
)

func LoggingMiddleware[S any](e *gateway.InteractionCreateEvent, state S, next ...Middleware[S]) error {
	sender := int64(e.SenderID())
	channelId := int64(e.ChannelID)

	if len(next) > 0 {
		t := time.Now()
		log.Printf("-> %d in %d %v", sender, channelId, e.Data.InteractionType())
		err := next[0](e, state, next[1:]...)
		log.Printf("<- %d in %d %v %v", sender, channelId, e.Data.InteractionType(), time.Since(t))
		return err
	}
	return nil
}
