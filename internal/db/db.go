package db

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5"
)

var DB *pgx.Conn

func Connect(dbURL string) {
	var err error

	DB, err = pgx.Connect(context.Background(), dbURL)
	if err != nil {
		log.Fatal("DB connection failed:", err)
	}

	err = DB.Ping(context.Background())
	if err != nil {
		log.Fatal("DB ping failed:", err)
	}

	log.Println("Connected to PostgreSQL")
}
