package broker

import (
	"context"
	"errors"

	"github.com/bwmarrin/discordgo"
	"github.com/google/uuid"
	"github.com/micro/go-micro/broker"
)

// discordBroker implements a go-micro broker for Discord Gateway events
type discordBroker struct {
	session *discordgo.Session
	options broker.Options
}

func (b *discordBroker) Options() broker.Options {
	return b.options
}

// Address returns an empty string as this isn't applicable for the Discord gateway
func (b *discordBroker) Address() string {
	return ""
}

// Connect causes the broker to initiaie a WebSocket session with the Discord Gateway
func (b *discordBroker) Connect() error {
	return b.session.Open()
}

// Disconnect causes the broker to disconnect any active WebSocket sessions with the Discord Gateway
func (b *discordBroker) Disconnect() error {
	return b.session.Close()
}

// Init initialises a Discord Gateway broker from parameters provided as arguments
func (b *discordBroker) Init(opts ...broker.Option) (err error) {
	token := "" // fix this
	b.session, err = discordgo.New(token)
	return
}

func (b *discordBroker) Publish(topic string, msg *broker.Message, opts ...broker.PublishOption) error {
	// todo event type
	ctx := context.Background()

	m := &discordgo.Event{
		Operation: op,
		Sequence:  s,
		Type:      t,
		RawData:   d,
	}

	return errors.New("UNIMPLEMENTED")
}

func (b *discordBroker) Subscribe(topic string, h broker.Handler, opts ...broker.SubscribeOption) (broker.Subscriber, error) {
	options := broker.SubscribeOptions{
		AutoAck: true,
		Queue:   "q-" + uuid.New().String(),
		Context: b.options.Context,
	}

	for _, o := range opts {
		o(&options)
	}

	ctx := context.Background()

	subscriber := &subscriber{
		options:   options,
		exit:      make(chan bool),
		eventType: topic,
		session:   b.session,
	}

	go subscriber.add(h)

	return subscriber, nil
}

// NewBroker creates a new Discord Gateway broker
func NewBroker(opts ...broker.Option) broker.Broker {
	options := broker.Options{
		Context: context.Background(),
	}

	for _, o := range opts {
		o(&options)
	}

	// retrieve Discord token
	token, _ := options.Context.Value(discordAuthToken{}).(string)
	// retrieve .... todo

	// create new Discord Gateway session; does NOT connect yet
	s, err := discordgo.New(token)
	if err != nil {
		panic(err.Error())
	}

	return &discordBroker{
		session: s,
		options: options,
	}
}
