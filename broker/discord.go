package broker

import (
	"github.com/micro/go-micro/cmd"
)

func init() {
	cmd.DefaultBrokers["discordgw"] = NewBroker
}
