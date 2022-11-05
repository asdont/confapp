package model

import (
	"context"
	"fmt"
	"time"

	"confapp/internal/app"
)

type ConfigRow struct {
	VersionID int
	Param     string
	Value     string
}

func GetServiceConfigs(tools *app.Tools, serviceName string) (map[int]map[string]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(),
		time.Second*time.Duration(tools.Conf.Postgres.QueryTimeoutSecond))
	defer cancel()

	rows, err := tools.DB.QueryContext(ctx, `
		select
		    c.version_id, c.param, c.value
		from
		    config c
		join
		    version v using (version_id)
		join
		    service s using (service_id)
		where
		    s.name = $1
			and v.deleted is null
	`,
		serviceName,
	)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}

	defer func() {
		if err := rows.Close(); err != nil {
			tools.Logger.App.Error("query: close rows: ", err)
		}
	}()

	var configsRows []ConfigRow

	for rows.Next() {
		var configRow ConfigRow

		if err := rows.Scan(
			&configRow.VersionID,
			&configRow.Param,
			&configRow.Value,
		); err != nil {
			return nil, fmt.Errorf("scan row: %w", err)
		}

		configsRows = append(configsRows, configRow)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("scan rows: %w", rows.Err())
	}

	serviceConfigs := make(map[int]map[string]string)

	for _, configRow := range configsRows {
		if serviceConfigs[configRow.VersionID] == nil {
			serviceConfigs[configRow.VersionID] = make(map[string]string)
		}

		serviceConfigs[configRow.VersionID][configRow.Param] = configRow.Value
	}

	return serviceConfigs, nil
}

func UpdateAllConfigs(tools *app.Tools, serviceConfigs map[int]map[string]string) error {
	ctx, cancel := context.WithTimeout(context.Background(),
		time.Second*time.Duration(tools.Conf.Postgres.QueryTimeoutSecond))
	defer cancel()

	tx, err := tools.DB.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("create transaction: %w", err)
	}

	for versionID, paramRows := range serviceConfigs {
		for param, value := range paramRows {
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
				if errTx := tx.Rollback(); errTx != nil {
					return fmt.Errorf("rollback: %w: %v", err, errTx)
				}

				return fmt.Errorf("exec: %w", err)
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit: %w", err)
	}

	return nil
}
