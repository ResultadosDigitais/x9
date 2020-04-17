package config

import "strconv"

type DatabaseConfig struct {
	Host         string `env:"POSTGRES_HOST,required"`
	Port         int    `env:"POSTGRES_PORT,required"`
	User         string `env:"POSTGRES_USER,required"`
	Password     string `env:"POSTGRES_PASSWORD,required"`
	DataBaseName string `env:"POSTGRES_DB,required"`
}

func ParseDatabaseConfig() (DatabaseConfig, error) {
	dbconfig := DatabaseConfig{}
	var err error
	if dbconfig.Host, err = getEnv("POSTGRES_HOST"); err != nil {
		return dbconfig, err
	}
	stringPort, err := getEnv("POSTGRES_PORT")
	if err != nil {
		return dbconfig, err
	}
	if dbconfig.Port, err = strconv.Atoi(stringPort); err != nil {
		return dbconfig, err
	}
	if dbconfig.User, err = getEnv("POSTGRES_USER"); err != nil {
		return dbconfig, err
	}
	if dbconfig.Password, err = getEnv("POSTGRES_PASSWORD"); err != nil {
		return dbconfig, err
	}
	if dbconfig.DataBaseName, err = getEnv("POSTGRES_DB"); err != nil {
		return dbconfig, err
	}
	return dbconfig, nil
}
