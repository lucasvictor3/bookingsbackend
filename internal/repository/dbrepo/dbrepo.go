package dbrepo

import (
	"database/sql"

	"github.com/lucasvictor3/bookingsbackend/internal/config"
	"github.com/lucasvictor3/bookingsbackend/internal/repository"
)

type postgresDBRepo struct {
	App *config.AppConfig
	DB  *sql.DB
}

func NewPostgresRepo(conn *sql.DB, a *config.AppConfig) repository.DatabaseRepo {
	return &postgresDBRepo{
		App: a,
		DB:  conn,
	}
}

type testDBRepo struct {
	App *config.AppConfig
	DB  *sql.DB
}

func NewTestingsRepo(a *config.AppConfig) repository.DatabaseRepo {
	return &testDBRepo{
		App: a,
	}
}
