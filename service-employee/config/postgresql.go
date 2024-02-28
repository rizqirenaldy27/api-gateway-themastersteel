package config

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

var pgPool *pgxpool.Pool

func InitPostgresDB() {
	config, err := pgxpool.ParseConfig("postgresql://admin-postgres:dbPostgres41@postgres:5432/api-gateway")
	if err != nil {
		panic(err)
	}

	pool, err := pgxpool.ConnectConfig(context.Background(), config)
	if err != nil {
		panic(err)
	}

	if err := pool.Ping(context.Background()); err != nil {
		panic(err)
	}

	fmt.Println("Pinged your deployment. You successfully connected to PostgreSQL!")

	pgPool = pool
}

func NewPostgresContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 10*time.Second)
}

func GetPgxPool() *pgxpool.Pool {
	return pgPool
}
