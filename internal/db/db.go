package db

import (
	"fmt"
	"gitlab-digital.tele2.kz/digital/eshop/backend/kaspi-pay/pkg/db/config"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	"github.com/sirupsen/logrus"
)

const migrationPath = "./migrations"

type Config struct {
	Username              string        `env:"DB_USERNAME,required"`
	Password              string        `env:"DB_PASSWORD,required"`
	Host                  string        `env:"DB_HOST,required"`
	Port                  string        `env:"DB_PORT,required"`
	Name                  string        `env:"DB_NAME,required"`
	Params                string        `env:"DB_PARAMS"`
	MaxOpenConnections    int           `env:"DB_MAX_OPEN_CONNECTIONS,required"`
	MaxIdleConnections    int           `env:"DB_MAX_IDLE_CONNECTIONS,required"`
	ConnectionMaxIdleTime time.Duration `env:"DB_CONNECTION_MAX_IDLE_TIME"`
}

func NewPostgresDB(config config.Config) (*sqlx.DB, error) {
	connectionString := fmt.Sprintf(
		"%s://%s:%s@%s:%s/%s?%s",
		config.Type,
		config.Username,
		config.Password,
		config.Host,
		config.Port,
		config.Name,
		config.Params,
	)

	sqlxDb, err := sqlx.Connect(config.Type, connectionString)
	if err != nil {
		return nil, err
	}
	sqlxDb.SetMaxOpenConns(config.MaxOpenConnections)
	sqlxDb.SetMaxIdleConns(config.MaxIdleConnections)
	sqlxDb.SetConnMaxIdleTime(config.ConnectionMaxIdleTime)

	logrus.Info("starting the migrations")
	if err = goose.Up(sqlxDb.DB, migrationPath); err != nil {
		logrus.Error(err)
	}
	logrus.Info("the migrations launch is over")

	return sqlxDb, nil
}
