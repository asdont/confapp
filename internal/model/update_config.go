package model

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"confapp/internal/app"
)

func UpdateConfig(tools *app.Tools, serviceName string, versionNumber int, params map[string]string) error {
	ctx, cancel := context.WithTimeout(context.Background(),
		time.Second*time.Duration(tools.Conf.Postgres.QueryTimeoutSecond))
	defer cancel()

	tx, err := tools.DB.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("create transaction: %w", err)
	}

	if err := updateConfig(tx, serviceName, versionNumber, params); err != nil {
		return fmt.Errorf("update config: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit: %w", err)
	}

	return nil
}

func updateConfig(tx *sql.Tx, serviceName string, versionNumber int, params map[string]string) error {
	var versionID int
	if err := tx.QueryRow(`
		select
		    v.version_id
		from
		    version v 
		join
		    service s using (service_id)
		where
		    s.name = $1
			and v.number = $2
			and s.deleted is null 
	`,
		serviceName,
		versionNumber,
	).Scan(
		&versionID,
	); err != nil {
		return fmt.Errorf("query: scan: %w", err)
	}

	for param, value := range params {
		if _, err := tx.Exec(`
			insert into
				config (version_id, param, value)
			values
				($1, $2, $3)
			on conflict
				(version_id, param)
			do update
				set value = $3
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
