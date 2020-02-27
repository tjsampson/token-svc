package postgres

import (
	"database/sql"
	"fmt"

	"github.com/tjsampson/token-svc/internal/config"
	internalerrors "github.com/tjsampson/token-svc/internal/errors"

	"github.com/lib/pq"

	// Registers the sql driver
	_ "github.com/lib/pq"

	"github.com/golang-migrate/migrate/v4"

	// import the source file
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// SQLConnStr parses out the config values into the Postgres Connection String
func sqlConnStr(c *config.Config) string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable&connect_timeout=%s", c.DB.User, c.DB.Pass, c.DB.Host, c.DB.Port, c.DB.Name, c.DB.Timeout)
}

// Database runs the migrations and returns the db connection
func Database(cfg *config.Config) (*sql.DB, error) {
	m, err := migrate.New("file://migrations", sqlConnStr(cfg))

	if err != nil {
		return nil, err
	}

	err = m.Up()

	if err != nil && err.Error() != "no change" {
		return nil, err
	}

	return sql.Open("postgres", sqlConnStr(cfg))
}

// ErrorCheck checks a database error and attempts to make sense of it or
func ErrorCheck(err error) error {
	if err == nil {
		return nil
	}

	// Trap Common SQL Errors
	if err == sql.ErrNoRows {
		return &internalerrors.RestError{
			Code:          404,
			Message:       "Resource not found",
			OriginalError: err,
		}
	}

	// Trap Postgres SQL Errors
	if pgerr, ok := err.(*pq.Error); ok {
		switch pgerr.Code.Name() {
		case "unique_violation":
			return &internalerrors.RestError{
				Code:          409,
				Message:       "Invalid request. Record already exists",
				OriginalError: pgerr,
			}
		default:
			return &internalerrors.RestError{
				Code:          500,
				Message:       fmt.Sprintf("unknown database error: %v", pgerr.Error()),
				OriginalError: pgerr,
			}
		}
	}
	return err
}
