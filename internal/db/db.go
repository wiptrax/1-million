package db

import (
	"context"
	// "database/sql"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var counts int64

func openDB() (*pgxpool.Pool, error) {
	ctx := context.Background()
	dbURL := "postgresql://postgres:password@localhost:5432/lowserver?sslmode=disable"

	config, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		return nil, err
	}
	config.MaxConns = 50 // adjust based on your CPU & DB capacity
	config.MinIdleConns = 25
	conn, err := pgxpool.NewWithConfig(context.Background(), config)

	if err != nil {
		return nil, err
	}

	err = conn.Ping(ctx)
	if err != nil {
		return nil, err
	}

	return conn, err

}

func ConnectToDB() *pgxpool.Pool {
	for {
		connection, err := openDB()
		if err != nil {
			log.Println("Could not connect to database, Postgres is not ready...")
			counts += 1
		} else {
			log.Println("Connected to database...")
			return connection
		}

		if counts > 10 {
			log.Println(err)
			return nil
		}

		log.Println("Waiting for database to become ready...")
		time.Sleep(2 * time.Second)
		continue
	}
}
