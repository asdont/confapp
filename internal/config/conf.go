package config

import (
	"errors"
	"fmt"

	"github.com/BurntSushi/toml"
	"github.com/go-playground/validator/v10"
)

var errInvalidField = errors.New("invalid field")

type Conf struct {
	Server   server
	Postgres postgres
	Logger   logger
}

type server struct {
	Port                    string `toml:"Port" validate:"numeric"`
	ReadTimeoutSec          int    `toml:"ReadTimeoutSec" validate:"gte=3,lte=180"`
	WriteTimeoutSec         int    `toml:"WriteTimeoutSec" validate:"gte=3,lte=180"`
	DeletionOperationsLogin string `toml:"DeletionOperationsLogin" validate:"gte=3"`
	DeletionOperationsPass  string `toml:"DeletionOperationsPass" validate:"gte=3"`
}

type postgres struct {
	Conn                  string `toml:"Conn" validate:"min=10"`
	MaxOpenConns          int    `toml:"MaxOpenConns" validate:"gte=1,lte=100"`
	MaxIdleConns          int    `toml:"MaxIdleConns" validate:"gte=1,lte=100"`
	ConnMaxIdleTimeSecond int    `toml:"ConnMaxIdleTimeSecond" validate:"gte=1,lte=180"`
	ConnMaxLifeTimeSecond int    `toml:"ConnMaxLifeTimeSecond" validate:"gte=1,lte=180"`
	QueryTimeoutSecond    int    `toml:"QueryTimeoutSecond" validate:"gte=2,lte=90"`
}

type logger struct {
	AppMaxLogSizeMb          int `toml:"AppMaxLogSizeMb" validate:"gte=1,lte=3000"`
	AppMaxNumOfBackups       int `toml:"AppMaxNumOfBackups" validate:"gte=1,lte=100"`
	AppMaxBackupAgeDay       int `toml:"AppMaxBackupAgeDay" validate:"gte=1,lte=3650"`
	ServerMaxLogSizeMb       int `toml:"ServerMaxLogSizeMb" validate:"gte=1,lte=3000"`
	ServerMaxNumOfBackups    int `toml:"ServerMaxNumOfBackups" validate:"gte=1,lte=100"`
	ServerAppMaxBackupAgeDay int `toml:"ServerAppMaxBackupAgeDay" validate:"gte=1,lte=3650"`
}

func GetConf(fileName string) (*Conf, error) {
	var conf *Conf
	if _, err := toml.DecodeFile(fileName, &conf); err != nil {
		return nil, fmt.Errorf("decode file: %w", err)
	}

	if err := validator.New().Struct(*conf); err != nil {
		var vErrors validator.ValidationErrors
		if errors.As(err, &vErrors) {
			if err := checkValidatorErr(vErrors); err != nil {
				return nil, fmt.Errorf("validator: check err: %w", err)
			}

			return nil, fmt.Errorf("validator: %w", err)
		}
	}

	return conf, nil
}

func checkValidatorErr(errs validator.ValidationErrors) error {
	for _, err := range errs {
		return fmt.Errorf("%w: %s(%s): see it <%v> want <%s=%s>",
			errInvalidField,
			err.StructNamespace(),
			err.Type(),
			err.Value(),
			err.ActualTag(),
			err.Param())
	}

	return nil
}
