package db

import (
	"database/sql"
	"fmt"

	"github.com/caarlos0/env/v6"
	_ "github.com/lib/pq"
)

var conn *sql.DB

type Postgres struct {
	Host         string `env:"POSTGRES_HOST,required"`
	Port         int    `env:"POSTGRES_PORT,required"`
	User         string `env:"POSTGRES_USER,required"`
	Password     string `env:"POSTGRES_PASSWORD,required"`
	DataBaseName string `env:"POSTGRES_DB,required"`
}

var config Postgres

func GetDB() error {
	err := env.Parse(&config)
	if err != nil {
		return err
	}
	err = Connect()
	return err
}

func Connect() error {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		config.Host, config.Port, config.User, config.Password, config.DataBaseName)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return err
	}

	err = db.Ping()
	if err != nil {
		return err
	}

	SetConn(db)
	return nil

}

func SetConn(dbconn *sql.DB) {
	conn = dbconn
}

func Query(query string, params ...interface{}) (*sql.Rows, error) {
	stmt, err := conn.Prepare(query)
	if err != nil {
		return nil, err
	}
	rows, err := stmt.Query(params...)
	return rows, err
}

func Exec(query string, params ...interface{}) error {
	stmt, err := conn.Prepare(query)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(params...)

	return err
}

func InitTables() error {
	createTableVuln := `CREATE TABLE IF NOT EXISTS vulnerabilities (
		id varchar(200) NOT NULL,
		internal_id varchar(200) NOT NULL UNIQUE,
		name varchar(200) NOT NULL,
		repository varchar(255) NOT NULL,
		filename varchar(255) NOT NULL,
		tool varchar(255) NOT NULL,
		value text NOT NULL,
		false_positive boolean NOT NULL DEFAULT false,
		issue_url varchar(255),
		PRIMARY KEY(id)
	)`

	createTableFP := `CREATE TABLE IF NOT EXISTS false_positives (
		context varchar(45) NOT NULL,
		fp_hash varchar(255) NOT NULL UNIQUE,
		FOREIGN KEY (vuln_id) REFERENCES vulnerabilities (id)
	)`

	_, err := conn.Exec(createTableVuln)
	if err != nil {
		return err
	}
	_, err = conn.Exec(createTableFP)
	if err != nil {
		return err
	}

	return nil
}
