package dal

import (
	"net/http"
	"strings"
	"sync"
	"time"

	"{{ cookiecutter.project_slug }}/internal/config"

	"github.com/olivere/elastic/v7"
	"github.com/rs/zerolog/log"
)

var (
	esOnce     sync.Once
	esInstance *elastic.Client
)

func GetES() *elastic.Client {
	esOnce.Do(initES)
	return esInstance
}

func initES() {
	t := http.DefaultTransport.(*http.Transport).Clone()
	t.MaxIdleConns = 100
	t.MaxConnsPerHost = 100
	t.MaxIdleConnsPerHost = 100

	httpClient := &http.Client{
		Timeout:   60 * time.Second,
		Transport: t,
	}

	cfg := config.Get().Databases.ElasticSearch
	opts := []elastic.ClientOptionFunc{
		elastic.SetHttpClient(httpClient),
		elastic.SetSniff(false),
		elastic.SetURL(strings.Split(cfg.Endpoint, ",")...),
		elastic.SetRetrier(
			elastic.NewBackoffRetrier(
				elastic.NewSimpleBackoff(100, 200, 300, 400, 500, 600, 700),
			),
		),
		elastic.SetHealthcheck(false), // https://github.com/olivere/elastic/issues/880
		elastic.SetInfoLog(ESInfoLogger{}),
	}
	if cfg.Username != "" && cfg.Password != "" {
		opts = append(opts, elastic.SetBasicAuth(cfg.Username, cfg.Password))
	}
	client, err := elastic.NewClient(opts...)
	if err != nil {
		log.Fatal().Err(err).Msg("initial elasticsearch failed")
	}
	esInstance = client
}

type ESInfoLogger struct{}

func (l ESInfoLogger) Printf(format string, v ...interface{}) {
	log.Info().Msgf(format, v...)
}
