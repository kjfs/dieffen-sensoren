package dbrepo

import (
	"database/sql"

	"github.com/kjfs/dieffe_sensor/internal/config"
	"github.com/kjfs/dieffe_sensor/internal/repository"
)

type postgresDBRepo struct {
	App *config.AppConfig
	Db  *sql.DB
}

type testDBRepo struct {
	App *config.AppConfig
	Db  *sql.DB
}

// Factory function:
func NewPostgresRepo(conn *sql.DB, a *config.AppConfig) repository.DatabaseRepo {
	return &postgresDBRepo{
		App: a,
		Db:  conn,
	}
}

func NewTestingRepo(a *config.AppConfig) repository.DatabaseRepo {
	return &testDBRepo{
		App: a,
	}
}
