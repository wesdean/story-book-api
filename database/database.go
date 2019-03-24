package database

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

type Database struct {
	db          *sql.DB
	Tx          *sql.Tx
	initOptions *DatabaseInitOptions
}

type DatabaseInitOptions struct {
	Host     string
	Port     int
	User     string
	Password string
	DBname   string
}

func NewDatabase(options *DatabaseInitOptions) (*Database, error) {
	if options == nil {
		options = &DatabaseInitOptions{}
	}

	if options.Host == "" {
		options.Host = "localhost"
	}
	if options.Port == 0 {
		options.Port = 5432
	}
	if options.User == "" {
		options.User = "storybook"
	}
	if options.Password == "" {
		options.Password = "storybook"
	}
	if options.DBname == "" {
		options.DBname = "storybook"
	}

	database := &Database{initOptions: options}

	var err error
	psqlInfo := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		database.initOptions.Host,
		database.initOptions.Port,
		database.initOptions.User,
		database.initOptions.Password,
		database.initOptions.DBname,
	)

	database.db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}

	return database, nil
}

func (db *Database) Begin() error {
	var err error

	if db.Tx == nil {
		db.Tx, err = db.db.Begin()
	}

	return err
}

func (db *Database) Commit() error {
	var err error

	if db.Tx != nil {
		err = db.Tx.Commit()
		db.Tx = nil
	}

	return err
}

func (db *Database) Rollback() error {
	var err error

	if db.Tx != nil {
		err = db.Tx.Rollback()
		db.Tx = nil
	}

	return err
}

func (db *Database) Ping() error {
	return db.db.Ping()
}

func (db *Database) GetDB() *sql.DB {
	return db.db
}
