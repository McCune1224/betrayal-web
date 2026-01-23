package internal

import (
	"context"
	"github.com/jackc/pgx/v5"
)

var Conn *pgx.Conn

func InitDB(ctx context.Context, dbURL string) error {
	conn, err := pgx.Connect(ctx, dbURL)
	if err != nil {
		return err
	}
	Conn = conn
	return nil
}

func CloseDB(ctx context.Context) error {
	if Conn != nil {
		return Conn.Close(ctx)
	}
	return nil
}
