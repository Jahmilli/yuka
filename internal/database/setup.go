package database

import (
	"context"
	"fmt"
	"yuka/internal/models"

	"github.com/cenkalti/backoff/v4"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	_ "gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func ConnectDatabase(parent context.Context,
	logger *zap.SugaredLogger,
	host string,
	user string,
	password string,
	dbname string,
	port string,
	sslmode string) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		host, user, password, dbname, port, sslmode)
	var db *gorm.DB
	connectDb := func() error {
		var err error
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
			// This is set to translate errors to a common error across databases engines (i.e Sqlite, Postgres)
			TranslateError: true,
		})
		if err != nil {
			return err
		}
		if err = db.AutoMigrate(
			&models.Device{},
			&models.Organization{},
			&models.User{},
			&models.RegisteredApplication{},
			&models.DeviceDNSQuery{},
			&models.Invitation{},
		); err != nil {
			return err
		}
		return nil
	}
	err := backoff.Retry(connectDb, backoff.WithContext(backoff.NewExponentialBackOff(), context.TODO()))
	if err != nil {
		return nil, err
	}
	logger.Info("initialised database connection")
	return db, nil
}

// ConnectDatabase creates a connection to the database, performs migration and returns the connection instance
func ConnectTestDatabase(logger *zap.Logger) (*gorm.DB, error) {
	database, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{
		// This is set to translate errors to a common error across databases engines (i.e Sqlite, Postgres)
		TranslateError: true,
	})

	if err != nil {
		return nil, err
	}

	// TODO: Will need to move to proper migrations later on
	if err = database.AutoMigrate(
		// &models.Device{},
		// &models.Organization{},
		&models.User{},
		// &models.RegisteredApplication{},
		// &models.DeviceDNSQuery{},
		// &models.Invitation{},
	); err != nil {
		println(err.Error())
		return nil, err
	}

	logger.Info("Initialised database connection")
	return database, nil
}
