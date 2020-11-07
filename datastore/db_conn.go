package datastore

import (
	"context"
)

type Repository interface {
	Begin(ctx context.Context) (RepoTx, error)
	Query(ctx context.Context, sqlQuery string, args ...interface{}) (pgx.Rows, error)
}

type RepoTx interface {
	Begin(ctx context.Context) (RepoTx, error)
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
	Query(ctx context.Context, sqlQuery string, args ...interface{}) (pgx.Rows, error)
}

type PgxTx struct {
	tx pgx.Tx
}

func NewPgxTx(tx pgx.Tx) Repository {
	return &PgxTx{tx: tx}
}

func (p PgxTx) Commit(ctx context.Context) error {
	return p.tx.Commit(ctx)
}

func (p PgxTx) Rollback(ctx context.Context) error {
	return p.tx.Rollback(ctx)
}

func (p PgxTx) Query(ctx context.Context, sqlQuery string, args ...interface{}) (pgx.Rows, error) {
	return p.tx.Query(ctx, sqlQuery, args...)
}

func (p PgxTx) Begin(ctx context.Context) (RepoTx, error) {
	tx, err := p.tx.Begin(ctx)
	if err != nil {
		return nil, err
	}

	return &PgxTx{tx: tx}, nil
}

type PgxConn struct {
	conn *pgx.Conn
}

func NewPgxConn(c *pgx.Conn) Repository {
	return &PgxConn{conn: c}
}

func (p PgxConn) Query(ctx context.Context, sqlQuery string, args ...interface{}) (pgx.Rows, error) {
	return p.conn.Query(ctx, sqlQuery, args...)
}

func (p PgxConn) Begin(ctx context.Context) (RepoTx, error) {
	tx, err := p.conn.Begin(ctx)
	return &PgxTx{tx: tx}, err
}
