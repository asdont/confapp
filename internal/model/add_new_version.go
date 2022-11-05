package model

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"confapp/internal/app"
)

func AddNewVersionConfig(tools *app.Tools, serviceName string, lastVersionNumber int, params map[string]string) error {
	ctx, cancel := context.WithTimeout(context.Background(),
		time.Second*time.Duration(tools.Conf.Postgres.QueryTimeoutSecond))
	defer cancel()

	tx, err := tools.DB.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("create transaction: %w", err)
	}

	versionID, err := addNewVersion(tx, serviceName, lastVersionNumber)
	if err != nil {
		if errTx := tx.Rollback(); errTx != nil {
			return fmt.Errorf("rollback: %w: %v", err, errTx)
		}

		return fmt.Errorf("add new version: %w", err)
	}

	if err := addNewConfig(tx, versionID, params); err != nil {
		if errTx := tx.Rollback(); errTx != nil {
			return fmt.Errorf("rollback: %w: %v", err, errTx)
		}

		return fmt.Errorf("add new config: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit: %w", err)
	}

	return nil
}

func addNewVersion(tx *sql.Tx, serviceName string, lastVersionNumber int) (int, error) {
	var versionID int

	if err := tx.QueryRow(`
		insert into
		    version (service_id, number)
		values
		    (
		     (select
		          service_id
		      from
		          service
		      where
		          name = $1),
		     $2)
		returning
			version_id
	`,
		serviceName,
		lastVersionNumber+1,
	).Scan(
		&versionID,
	); err != nil {
		return 0, fmt.Errorf("query: %w", err)
	}

	return versionID, nil
}

func addNewConfig(tx *sql.Tx, versionID int, params map[string]string) error {
	for param, value := range params {
		if _, err := tx.Exec(`
			insert into
				config (version_id, param, value)
			values
				($1, $2, $3)
		`,
			versionID,
			param,
			value,
		); err != nil {
			return fmt.Errorf("exec: %w", err)
		}
	}

	return nil
}
