package model

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"confapp/internal/app"
)

const defaultServiceVersion = 1

func AddConfig(tools *app.Tools, serviceName string, params map[string]string) (int, error) {
	ctx, cancel := context.WithTimeout(
		context.Background(), time.Second*time.Duration(tools.Conf.Postgres.QueryTimeoutSecond),
	)
	defer cancel()

	tx, err := tools.DB.BeginTx(ctx, nil)
	if err != nil {
		return 0, fmt.Errorf("create transaction: %w", err)
	}

	serviceID, err := addServiceName(tx, serviceName)
	if err != nil {
		if errTx := tx.Rollback(); errTx != nil {
			return 0, fmt.Errorf("rollback: %w: %v", err, errTx)
		}

		return 0, fmt.Errorf("add service: name: %w", err)
	}

	versionID, err := addVersion(tx, serviceID)
	if err != nil {
		if errTx := tx.Rollback(); errTx != nil {
			return 0, fmt.Errorf("rollback: %w: %v", err, errTx)
		}

		return 0, fmt.Errorf("add service: version: %w", err)
	}

	if err := addConfig(tx, versionID, params); err != nil {
		if errTx := tx.Rollback(); errTx != nil {
			return 0, fmt.Errorf("rollback: %w: %v", err, errTx)
		}

		return 0, fmt.Errorf("add service: config: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("commit: %w", err)
	}

	return defaultServiceVersion, nil
}

func addServiceName(tx *sql.Tx, serviceName string) (int, error) {
	var serviceID int

	if err := tx.QueryRow(`
		insert into
		    service (name)
		values
		    ($1)
		returning
			service_id
	`,
		serviceName,
	).Scan(
		&serviceID,
	); err != nil {
		return 0, fmt.Errorf("query: %w", err)
	}

	return serviceID, nil
}

func addVersion(tx *sql.Tx, serviceID int) (int, error) {
	var versionID int

	if err := tx.QueryRow(`
		insert into
		    version (service_id, number)
		values
		    ($1, $2)
		returning
			version_id
	`,
		serviceID,
		defaultServiceVersion,
	).Scan(
		&versionID,
	); err != nil {
		return 0, fmt.Errorf("query: %w", err)
	}

	return versionID, nil
}

func addConfig(tx *sql.Tx, versionID int, params map[string]string) error {
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
