package internal

import (
	"context"
	"github.com/jackc/pgx/v5"
)

type Pingable interface {
	Ping(context.Context) error
}

var Conn Pingable

func InitDB(ctx context.Context, dbURL string) error {
	conn, err := pgx.Connect(ctx, dbURL)
	if err != nil {
		return err
	}
	Conn = conn
	return nil
}

func CloseDB(ctx context.Context) error {
	if real, ok := Conn.(*pgx.Conn); ok && real != nil {
		return real.Close(ctx)
	}
	return nil
}
