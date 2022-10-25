package model

import (
	"context"
	"fmt"
	"time"

	"confapp/internal/app"
)

func DeleteConfigVersion(tools *app.Tools, serviceName string, versionNumber int) error {
	ctx, cancel := context.WithTimeout(context.Background(),
		time.Second*time.Duration(tools.Conf.Postgres.QueryTimeoutSecond))
	defer cancel()

	if _, err := tools.DB.ExecContext(ctx, `
		update version
			set deleted = now()
		where
			service_id=(
				select
					service_id
				from
					service
				where
					name = $1
			)
			and number = $2
	`,
		serviceName,
		versionNumber,
	); err != nil {
		return fmt.Errorf("exec: %w", err)
	}

	return nil
}
