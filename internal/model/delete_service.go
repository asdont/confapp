package model

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"confapp/internal/app"
)

func DeleteServiceEasy(tools *app.Tools, serviceName string) error {
	ctx, cancel := context.WithTimeout(context.Background(),
		time.Second*time.Duration(tools.Conf.Postgres.QueryTimeoutSecond))
	defer cancel()

	tx, err := tools.DB.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("create transaction: %w", err)
	}

	serviceID, err := deleteServiceEasy(tx, serviceName)
	if err != nil {
		if errTx := tx.Rollback(); errTx != nil {
			return fmt.Errorf("rollback: %w: %v", err, errTx)
		}

		return fmt.Errorf("delete service: %w", err)
	}

	if err := deleteVersionsEasy(tx, serviceID); err != nil {
		if errTx := tx.Rollback(); errTx != nil {
			return fmt.Errorf("rollback: %w: %v", err, errTx)
		}

		return fmt.Errorf("delete versions: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit: %w", err)
	}

	return nil
}

func deleteServiceEasy(tx *sql.Tx, serviceName string) (int, error) {
	var serviceID int

	if err := tx.QueryRow(`
		update service
			set deleted = now()
		where
			name=$1
			and deleted is null
		returning
			service_id			
	`,
		serviceName,
	).Scan(
		&serviceID,
	); err != nil {
		return 0, fmt.Errorf("exec: %w", err)
	}

	return serviceID, nil
}

func deleteVersionsEasy(tx *sql.Tx, serviceID int) error {
	if _, err := tx.Exec(`
		update version
			set deleted = now()
		where
			service_id=$1
			and deleted is null
	`,
		serviceID,
	); err != nil {
		return fmt.Errorf("exec: %w", err)
	}

	return nil
}
