package arikawax

import (
	"log"
	"runtime/debug"
	"time"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
)

type Middleware[S any] func(e *gateway.InteractionCreateEvent, state *S, next ...Middleware[S]) error

func LoggingMiddleware[S any](e *gateway.InteractionCreateEvent, state *S, next ...Middleware[S]) error {
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

type ITokenRec struct {
	Token string
	Time  time.Time
}

var iTokenRB RingBuffer[ITokenRec]
var tdmIn chan ITokenRec
var tdmOut chan string
var BufferSize int = 128

// TimeoutDetectionMiddleware is a middleware that logs a warning when an interaction takes longer than 4 seconds (Discord requires a response within 3 seconds)
// Noted that this middleware at the moment uses a single state and goroutine, so it might not be thread safe.
func TimeoutDetectionMiddleware[S any](e *gateway.InteractionCreateEvent, state *S, next ...Middleware[S]) error {
	if tdmIn == nil {
		tdmIn = make(chan ITokenRec, BufferSize)
		tdmOut = make(chan string, BufferSize)
		go func() {
			aknowledged := map[string]struct{}{}
			for {
				select {
				case rec := <-tdmIn:
					iTokenRB.Push(rec)
				case token := <-tdmOut:
					aknowledged[token] = struct{}{}
				default:
					rec, ok := iTokenRB.Peek()
					if !ok {
						time.Sleep(time.Millisecond * 100)
						continue
					}
					if _, ok := aknowledged[rec.Token]; ok {
						delete(aknowledged, rec.Token)
						continue
					}
					if time.Since(rec.Time) < time.Second*4 {
						for time.Since(rec.Time) < time.Second*4 {
							select {
							case rec := <-tdmIn:
								iTokenRB.Push(rec)
							case token := <-tdmOut:
								aknowledged[token] = struct{}{}
							default:
								time.Sleep(time.Millisecond * 100)
							}
						}
					}
					iTokenRB.Pop()

					since := time.Since(rec.Time)
					switch data := e.Data.(type) {
					case *discord.CommandInteraction:
						log.Printf("[Timeout] %v @%d | Command | %v | target: %d | options: %v", since, e.SenderID(), data.Name, data.TargetID, data.Options)
					case *discord.ButtonInteraction:
						log.Printf("[Timeout] %v @%d | Button | %v", since, e.SenderID(), data.CustomID)
					case *discord.StringSelectInteraction:
						log.Printf("[Timeout] %v @%d | StringSelect | %v | %v", since, e.SenderID(), data.CustomID, data.Values)
					case *discord.ModalInteraction:
						log.Printf("[Timeout] %v @%d | Modal | %v", since, e.SenderID(), data)
					case *discord.AutocompleteInteraction:
						log.Printf("[Timeout] %v @%d | Autocomplete | %v", since, e.SenderID(), data)
					case *discord.PingInteraction:
						log.Printf("[Timeout] %v @%d | Ping", since, e.SenderID())
					case *discord.UnknownInteractionData:
						log.Printf("[Timeout] %v @%d | Unknown | %s", since, e.SenderID(), data.Raw)
					}
				}
			}
		}()
	}

	tdmIn <- ITokenRec{Token: e.Token, Time: time.Now()}

	if len(next) > 0 {
		return next[0](e, state, next[1:]...)
	}
	return nil
}
