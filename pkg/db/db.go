package db

import (
	"database/sql"
	"errors"
	"time"

	_ "github.com/lib/pq"
)

// Connect makes a connection to the PostgreSQL database
// and returns the sql.DB handler representing a pool of zero or more
// underlying connections.
func Connect(connString string, maxConnections int) (*sql.DB, error) {
	db, err := sql.Open("postgres", connString)
	db.SetMaxOpenConns(maxConnections)

	if err != nil {
		return nil, err
	}

	// Try to connect up to 5 times
	for retry := 1; retry <= 5; retry++ {
		err = db.Ping()
		if err == nil {
			return db, nil
		}
		time.Sleep(1 * time.Second)
	}

	return nil, err
}

// Setup uses the existing database connection to create
// the necessary tables for the API (if they don't exist).
func Setup(db *sql.DB) error {
	if db == nil {
		return errors.New("cannot setup database, must call Connect() first")
	}

	if _, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS app (
	    	id uuid NOT NULL UNIQUE,
	    	app_id character varying NOT NULL PRIMARY KEY,
	    	app_name character varying,
			deleted_at timestamp without time zone
		);

		CREATE TABLE IF NOT EXISTS version (
			id serial NOT NULL PRIMARY KEY,
  		  	version character varying NOT NULL,
  		  	app_id character varying NOT NULL REFERENCES app(app_id),
  		  	disabled boolean DEFAULT false NOT NULL,
  		  	disabled_message character varying,
  		  	num_of_clients integer DEFAULT 0 NOT NULL,
  		  	num_of_app_launches integer DEFAULT 0 NOT NULL,
  		  	unique (app_id, version)
		);
	`); err != nil {
		return err
	}

	return nil
}
