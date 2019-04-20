package broker

import (
	"strconv"

	"github.com/bwmarrin/discordgo"
	"github.com/micro/go-micro/broker"
)

type subscriber struct {
	options   broker.SubscribeOptions
	exit      chan bool
	eventType string

	//broker       *discordBroker
	session       *discordgo.Session
	eventHandler  *discordgo.EventHandler
	brokerHandler broker.Handler
	remover       func()
}

func (s *subscriber) Options() broker.SubscribeOptions {
	return s.options
}

func (s *subscriber) Topic() string {
	return (*s.eventHandler).Type()
}

func (s *subscriber) Type() string {
	return s.eventType
}

type Publication discordgo.Event

func (p Publication) Topic() string {
	return p.Type
}

func (p Publication) Message() *broker.Message {
	return &broker.Message{
		Header: map[string]string{
			"Operation": strconv.Itoa(p.Operation),
			"Type":      p.Type,
			"Sequence":  strconv.FormatInt(p.Sequence, 10),
		},
		Body: p.RawData,
	}
}

func (p Publication) Ack() error {
	return nil
}

func (s *subscriber) Handle(ds *discordgo.Session, m interface{}) {
	// event, ok := m.(discordgo.Event)
	publication, ok := m.(Publication)
	if !ok {
		panic("received incompatible message")
	}

	// publication := Publication(event)
	s.brokerHandler(publication)
}

func (s *subscriber) Unsubscribe() error {
	select {
	case <-s.exit:
		return nil
	default:
		close(s.exit)
		s.remover()
		return nil
	}
}

func (s *subscriber) add(hdlr broker.Handler) {
	//ctx, cancel := context.WithCancel(context.Background())
	s.brokerHandler = hdlr
	s.session.AddHandler(s)
}
