package postgres

import (
	"fmt"

	"github.com/davidafdal/post-app/config"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

func InitPostgres(cfg *config.PostgresConfig) (*sqlx.DB, error) {

	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Database,
	)

	db, err := sqlx.Connect("pgx", dsn)

	if err != nil {
		return db, err
	}
	return db, nil
}
