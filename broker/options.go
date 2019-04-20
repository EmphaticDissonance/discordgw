package broker

import (
	"context"

	"github.com/micro/go-micro/broker"
)

type discordAuthToken struct{}

type addEventHandler struct{}

// AuthToken is a broker Option which provides for setting the Discord authentication token
func AuthToken(token string) broker.Option {
	return func(o *broker.Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, discordAuthToken{}, token)
	}
}
