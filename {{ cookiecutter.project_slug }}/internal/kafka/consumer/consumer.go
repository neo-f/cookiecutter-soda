package consumer

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"{{ cookiecutter.project_slug }}/internal/config"

	slog "log"

	"github.com/Shopify/sarama"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"golang.org/x/exp/maps"
)

func Start(ctx context.Context) {
	// 初始化consumer协程池
	sarama.PanicHandler = func(i interface{}) {
		log.Error().Interface("error", i).Msg("kafka message recovered!!")
	}

	cfg := config.Get().Kafka
	sarama.Logger = slog.New(os.Stderr, "[Sarama] ", slog.LstdFlags)
	kafkaCfg := sarama.NewConfig()
	if cfg.Username != "" && cfg.Password != "" {
		kafkaCfg.Net.SASL.Enable = true
		kafkaCfg.Net.SASL.User = cfg.Username
		kafkaCfg.Net.SASL.Password = cfg.Password
	}
	kafkaCfg.Consumer.Return.Errors = true
	kafkaCfg.ChannelBufferSize = 1024

	consumerGroup, err := sarama.NewConsumerGroup(strings.Split(cfg.Endpoint, " "), "{{ cookiecutter.project_slug }}-service", kafkaCfg)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("failed to create consumer group")
		return
	}

	fmt.Println("[kafka] start consuming ...")
	for {
		select {
		case err := <-consumerGroup.Errors():
			log.Error().Err(err).Msg("kafka error")
		case <-ctx.Done():
			fmt.Println("[kafka] stop consuming ...")
			consumerGroup.Close()
			fmt.Println("[kafka] stopped")
			return
		default:
			if err := consumerGroup.Consume(ctx, maps.Keys(Topics), ConsumeRouter{}); err != nil {
				log.Ctx(ctx).Error().Err(err).Msg("[kafka] consuming error")
			}
		}
	}
}

type ConsumeRouter struct{}

// ConsumeClaim implements sarama.ConsumerGroupHandler
func (ConsumeRouter) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		fn, ok := Topics[message.Topic]
		if !ok {
			log.Error().Str("topic", message.Topic).Msg("specified topic route not found, ignoring")
			return nil
		}
		msg := message

		go func() {
			t := time.Now()
			defer func() {
				if recovered := recover(); recovered != nil {
					if e, ok := recovered.(error); ok {
						log.Error().Stack().Err(e).Msg("kafka consumer recovered")
					} else {
						log.Error().Stack().Interface("recoverd", recovered).Msg("kafka consumer recovered")
					}
				}
			}()
			log.Info().
				Str("topic", msg.Topic).
				Int32("partition", msg.Partition).
				Int64("offset", msg.Offset).
				Str("member_id", sess.MemberID()).
				RawJSON("value", msg.Value).
				Msg("consuming claim")

			ctx := sess.Context()
			if config.Get().Debug {
				ctx = log.With().Stack().Logger().Level(zerolog.DebugLevel).WithContext(sess.Context())
			}
			err := fn(ctx, msg)
			if err == nil {
				sess.MarkMessage(msg, "")
				log.Info().
					Str("topic", msg.Topic).
					Dur("elapse", time.Since(t)).
					Msg("consume successed")
			} else {
				log.Error().
					Stack().
					Str("topic", msg.Topic).
					RawJSON("msg", msg.Value).
					Dur("elapse", time.Since(t)).
					Err(err).
					Msg("consume failed")
			}
		}()
	}
	return nil
}

// Cleanup implements sarama.ConsumerGroupHandler
func (ConsumeRouter) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// Setup implements sarama.ConsumerGroupHandler
func (ConsumeRouter) Setup(c sarama.ConsumerGroupSession) error {
	return nil
}
