package app

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	log "github.com/sirupsen/logrus"
)

const (
	defaultAttempts = 10
	defaultTimeout  = time.Second
)

func Migrations() {
	log.Info("Configuring migrations...")
	pgUrl, ok := os.LookupEnv("POSTGRES_CONN")
	if !ok || len(pgUrl) == 0 {
		log.Fatal("app - init - os.LookupEnv: POSTGRES_CONN not specified")
	}
	pgUrl += "?sslmode=disable"

	var (
		connAttempts = defaultAttempts
		err          error
		mgrt         *migrate.Migrate
	)

	for connAttempts > 0 {
		mgrt, err = migrate.New("file://migrations", pgUrl)
		if err == nil {
			break
		}

		time.Sleep(defaultTimeout)
		log.Infof("Postgres trying to connect, attempts left: %d", connAttempts)
		connAttempts--
	}

	if err != nil {
		log.Fatal(fmt.Errorf("app - init - migrate.New: %w", err))
	}
	defer mgrt.Close()

	if err = mgrt.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatal(fmt.Errorf("app - init - mgrt.Up: %w", err))
	}

	if errors.Is(err, migrate.ErrNoChange) {
		log.Info("Migration no change...")
		return
	}

	log.Info("Migration successful up...")
}
