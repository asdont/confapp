package postgres

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"confapp/internal/config"
)

const (
	delayBeforeConnToDBSecond  = 5
	numberOfAttemptToGetConnDB = 5
)

func CreatePool(conf *config.Conf) (*sql.DB, error) {
	var errConn error

	for i := 0; i < numberOfAttemptToGetConnDB; i++ {
		log.Printf("attempt %d/%d to connect db... wait\n", i, numberOfAttemptToGetConnDB)

		conn, err := createPool(conf)
		if err != nil {
			log.Fatal(err)
		}

		if err := conn.Ping(); err != nil {
			log.Printf("fail: connect db: %v", err)

			errConn = err

			time.Sleep(time.Second * delayBeforeConnToDBSecond)

			continue
		}

		return conn, nil
	}

	return nil, fmt.Errorf("db ping: %w", errConn)
}

func createPool(conf *config.Conf) (*sql.DB, error) {
	conn, err := sql.Open("postgres", conf.Postgres.Conn)
	if err != nil {
		return nil, fmt.Errorf("sql open: %w", err)
	}

	// Максимальное количество одновременно открытых соединений
	conn.SetMaxOpenConns(conf.Postgres.MaxOpenConns)

	// Максимальное количество неактивных соединений в пуле
	conn.SetMaxIdleConns(conf.Postgres.MaxIdleConns)

	// Максимальное количество времени, в течение которого соединение может быть использовано повторно
	conn.SetConnMaxIdleTime(time.Second * time.Duration(conf.Postgres.ConnMaxIdleTimeSecond))

	// Максимальное время жизни, после того как соединение вернулось в пул
	conn.SetConnMaxLifetime(time.Second * time.Duration(conf.Postgres.ConnMaxLifeTimeSecond))

	return conn, nil
}
