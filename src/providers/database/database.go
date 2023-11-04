package database

import (
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"zgrabi-mjesto.hr/backend/src/entities/product"
)

type databaseProvider struct {
	conn *gorm.DB
}

var db_ = databaseProvider{}

func migrate() databaseProvider {
	return db_
}

func DatabaseProvider() databaseProvider {
	return db_
}

func (p databaseProvider) Client() *gorm.DB {
	return db_.conn
}

func (p databaseProvider) log(format string, v ...interface{}) {
	log.Printf("|DB> "+format, v...)
}

func (p databaseProvider) Register() (err error) {
	db_.log("Connecting to database...")

	logger_ := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags|log.LUTC),
		logger.Config{
			SlowThreshold:             200 * time.Millisecond,
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)

	db_.conn, err = gorm.Open(
		mysql.Open(os.Getenv("DATABASE_URL")),
		&gorm.Config{
			Logger:                                   logger_,
			DisableForeignKeyConstraintWhenMigrating: true,
		},
	)

	db_.log("Connected to to database")
	if err != nil {
		return err
	}

	db_.log("Running the migrations...")
	err = db_.conn.AutoMigrate(
		&product.Model{},
	)
	if err != nil {
		return err
	}

	db_.log("Done with migrations")

	return nil
}
