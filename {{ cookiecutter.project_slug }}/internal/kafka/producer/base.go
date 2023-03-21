package producer

import (
	"{{ cookiecutter.project_slug }}/internal/config"
	"encoding/json"
	"strings"
	"sync"

	"github.com/Shopify/sarama"
	"github.com/rs/zerolog/log"
)

var once sync.Once
var producerInstance sarama.AsyncProducer

func setupProducer() {
	kafkaCfg := sarama.NewConfig()
	kafkaCfg.Producer.Return.Errors = true
	kafkaCfg.Producer.Return.Successes = true

	cfg := config.Get().Kafka
	if cfg.Username != "" && cfg.Password != "" {
		kafkaCfg.Net.SASL.Enable = true
		kafkaCfg.Net.SASL.User = cfg.Username
		kafkaCfg.Net.SASL.Password = cfg.Password
	}
	producer, err := sarama.NewAsyncProducer(strings.Split(cfg.Endpoint, " "), kafkaCfg)
	if err != nil {
		log.Fatal().Err(err).Msg("setup producer failed")
	}
	go func() {
		for {
			select {
			case err := <-producer.Errors():
				log.Warn().Str("topic", err.Msg.Topic).Err(err).Msg("producer error")
			case pac := <-producer.Successes():
				log.Info().Str("topic", pac.Topic).Err(err).Msg("producer successed")
			}
		}
	}()
	producerInstance = producer
}

func produce(topic string, value interface{}) {
	once.Do(setupProducer)
	bv, err := json.Marshal(value)
	if err != nil {
		log.Error().Err(err).Msg("[kafka producer]json marshal failed")
		return
	}
	producerInstance.Input() <- &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(bv),
	}
}
