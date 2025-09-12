package database

import (
	"database/sql"
	"log"

	_ "github.com/jackc/pgx/v5" // импортируем pgx для регистрации драйвера database/sql
)

func Connect(dbParams string) (*sql.DB, error) {
	log.Println("Connecting to DB", dbParams)
	db, err := sql.Open("pgx", dbParams)
	if err != nil {
		return nil, err
	}

	return db, nil
}
