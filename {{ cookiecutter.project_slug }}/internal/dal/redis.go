package dal

import (
	"sync"

	"{{ cookiecutter.project_slug }}/internal/config"

	"github.com/go-redis/redis/v8"
	"github.com/rs/zerolog/log"
)

var (
	rdsOnce     sync.Once
	rdsInstance *redis.Client
)

func GetRedis() *redis.Client {
	rdsOnce.Do(initRedis)
	return rdsInstance
}

func initRedis() {
	log.Info().Msg("loading redis configs")
	dsn := config.Get().Databases.Redis
	cfg, err := redis.ParseURL(dsn)
	if err != nil {
		log.Fatal().Err(err).Msg("parse redis url failed")
	}
	rdsInstance = redis.NewClient(cfg)
	log.Info().Msg("connected to redis")
}
