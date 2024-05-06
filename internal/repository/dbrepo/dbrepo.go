package dbrepo

import (
	"database/sql"
	"myapp/internal/config"
	"myapp/internal/repository"
)

type postgresDBRepo struct {
	App *config.AppConfig
	DB  *sql.DB
}

// allows us to populate postgresDBRepo and return the repo
// since we are using postgres it will intialize a postgres connection pool
func NewPostgresRepo(conn *sql.DB, a *config.AppConfig) repository.DatabaseRepo {
	return &postgresDBRepo{
		App: a,
		DB:  conn,
	}
}
