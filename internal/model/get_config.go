package model

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"confapp/internal/app"
)

func GetConfig(tools *app.Tools, serviceName string, versionNumber int) (map[string]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(),
		time.Second*time.Duration(tools.Conf.Postgres.QueryTimeoutSecond))
	defer cancel()

	rows, err := tools.DB.QueryContext(ctx, `
		select
		    param, value
		from
		    config
		join
		    version v using (version_id)
		join
		    service s using (service_id)
		where
		    s.name = $1
			and v.number = $2
			and v.deleted is null
`,
		serviceName,
		versionNumber,
	)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}

	defer func() {
		if err := rows.Close(); err != nil {
			tools.Logger.App.Errorf("query: close rows: %v", err)
		}
	}()

	params, err := getConfigScanRows(rows)
	if err != nil {
		return nil, err
	}

	return params, nil
}

func getConfigScanRows(rows *sql.Rows) (map[string]string, error) {
	params := make(map[string]string)

	for rows.Next() {
		var param, value string

		if err := rows.Scan(
			&param,
			&value,
		); err != nil {
			return nil, fmt.Errorf("scan row: %w", err)
		}

		params[param] = value
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("scan rows: %w", rows.Err())
	}

	return params, nil
}

func GetLastVersionConfig(tools *app.Tools, serviceName string) (int, map[string]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(),
		time.Second*time.Duration(tools.Conf.Postgres.QueryTimeoutSecond))
	defer cancel()

	rows, err := tools.DB.QueryContext(ctx, `
		select
			number, param, value
		from
			(select
				dense_rank() over (order by v.number desc) row_num,
				v.number,
				c.param,
				c.value
			from
				config c
		   	join
				version v using (version_id)
		   	join
				service s using (service_id)
			where
			    s.name = $1
				and v.deleted is null
			) temp
			where
			    row_num = 1   
	`,
		serviceName,
	)
	if err != nil {
		return 0, nil, fmt.Errorf("query: %w", err)
	}

	defer func() {
		if err := rows.Close(); err != nil {
			tools.Logger.App.Errorf("get last versions: close rows: %v", err)
		}
	}()

	versionNumber, params, err := getLastVersionConfigScanRows(rows)
	if err != nil {
		return 0, nil, err
	}

	return versionNumber, params, nil
}

func getLastVersionConfigScanRows(rows *sql.Rows) (int, map[string]string, error) {
	params := make(map[string]string)

	var versionNumber int

	for rows.Next() {
		var param, value string

		if err := rows.Scan(
			&versionNumber,
			&param,
			&value,
		); err != nil {
			return 0, nil, fmt.Errorf("scan row: %w", err)
		}

		params[param] = value
	}

	if rows.Err() != nil {
		return 0, nil, fmt.Errorf("scan rows: %w", rows.Err())
	}

	return versionNumber, params, nil
}
