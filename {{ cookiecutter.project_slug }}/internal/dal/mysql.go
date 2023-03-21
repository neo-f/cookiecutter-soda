package dal

import (
	slog "log"
	"os"
	"sync"
	"time"

	"{{ cookiecutter.project_slug }}/internal/config"

	"github.com/rs/zerolog/log"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	sqlOnce     sync.Once
	sqlInstance *gorm.DB
)

func GetMySQLTest() *gorm.DB {
	l := logger.Default
	l = l.LogMode(logger.Info)
	db, _ := gorm.Open(sqlite.Open("file::memory:?parseTime=True&loc=Local"), &gorm.Config{Logger: l})
	return db
}

func GetMySQL() *gorm.DB {
	sqlOnce.Do(initMysql)
	return sqlInstance
}

func initMysql() {
	log.Info().Msg("loading mysql configs")
	dsn := config.Get().Databases.MySQL
	db, err := gorm.Open(
		mysql.Open(dsn),
		&gorm.Config{
			DisableForeignKeyConstraintWhenMigrating: true,
			Logger: logger.New(
				slog.New(os.Stdout, "\r\n", slog.LstdFlags), logger.Config{
					SlowThreshold:             200 * time.Millisecond,
					LogLevel:                  logger.Info,
					IgnoreRecordNotFoundError: false,
					Colorful:                  true,
				},
			),
		},
	)
	if err != nil {
		log.Fatal().Str("dsn", dsn).Err(err).Msg("connected to mysql failed")
	}
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal().Str("dsn", dsn).Err(err).Msg("connected to mysql failed")
	}
	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxIdleConns(50)

	// SetMaxOpenConns sets the maximum number of open connections to the database.
	sqlDB.SetMaxOpenConns(50)

	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	sqlDB.SetConnMaxLifetime(time.Minute * 2)
	if err != nil {
		log.Fatal().Str("dsn", dsn).Err(err).Msg("connected to mysql failed")
	}
	if config.Get().Debug {
		db = db.Debug()
	}
	sqlInstance = db
	log.Info().Msg("connected to mysql")
}
