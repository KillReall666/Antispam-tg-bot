package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Database struct {
	db *pgxpool.Pool
}

const createUsersTableQuery = `
      CREATE TABLE IF NOT EXISTS user_under_attack (
    UserID VARCHAR(255) PRIMARY KEY,
    UserName VARCHAR(255) UNIQUE,
    Count VARCHAR(255),
    CONSTRAINT unique_person UNIQUE (UserID)
);`

func New(connString string) (*Database, error) {
	conn, err := pgxpool.New(context.Background(), connString)
	if err != nil {
		return nil, fmt.Errorf("error connecting to db: %v", err)
	}

	_, err = conn.Exec(context.Background(), createUsersTableQuery)
	if err != nil {
		return nil, fmt.Errorf("error creating user table: %v", err)
	}

	return &Database{db: conn}, nil
}
