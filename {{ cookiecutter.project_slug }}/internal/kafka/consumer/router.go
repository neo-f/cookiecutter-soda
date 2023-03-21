package consumer

import (
	"context"

	"github.com/Shopify/sarama"
)

type ConsumeFunc func(context.Context, *sarama.ConsumerMessage) error

var Topics = map[string]ConsumeFunc{}
