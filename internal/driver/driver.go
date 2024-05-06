// this will be the means we connect our database to

package driver

import (
	"database/sql"
	"time"

	_ "github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
)

// We put the connection pools in a struct in the event we have more tha one database, we could simply add it to here
type DB struct {
	SQL *sql.DB
}

var dbConn = &DB{}

// We now define properties of our connection pool
const maxOpenDbCon = 10
const maxIdleDbConn = 5
const maxDbLifetime = 5 * time.Minute

// this function will create a database pool for Postgres
// This will take a database connection string, and return two things, a pointer to the DB type or an err
func ConnectSQL(dsn string) (*DB, error) {
	d, err := NewDatabase(dsn)
	if err != nil {
		panic(err)
	}

	// now we set the db connection properties so they do not go out of control
	d.SetMaxOpenConns(maxOpenDbCon)
	d.SetMaxIdleConns(maxIdleDbConn)
	d.SetConnMaxLifetime(maxDbLifetime)

	// now we need assign a value to our dbConn global var
	dbConn.SQL = d
	if err != nil {
		return nil, err
	}
	return dbConn, nil
}

// attempts to ping the database to ensure we are connected
func testDB(d *sql.DB) error {
	err := d.Ping()
	if err != nil {
		return err
	}
	return nil
}

// creates a new database for the application that we will get from postgres
func NewDatabase(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	// check for an error and see if it is nil all in one step
	if err = db.Ping(); err != nil {
		return nil, err

	}
	return db, nil
}
