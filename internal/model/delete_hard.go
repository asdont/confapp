package model

import (
	"fmt"
	"time"

	"confapp/internal/app"
)

func DeleteOldRowsHard(tools *app.Tools) error {
	for {
		if _, err := tools.DB.Exec(`
			delete from
				service
			where
				deleted < now() - INTERVAL '90 days' 
		`,
		); err != nil {
			return fmt.Errorf("exec: %w", err)
		}

		if _, err := tools.DB.Exec(`
			delete from
				version
			where
				deleted < now() - INTERVAL '90 days' 
		`,
		); err != nil {
			return fmt.Errorf("exec: %w", err)
		}

		time.Sleep(time.Hour)
	}
}
