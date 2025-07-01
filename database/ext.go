package database

import (
	"context"
	"database/sql"
)

// Exec runs a raw SQL statement via Ent’s underlying driver.
func (c *Client) Exec(ctx context.Context, query string, args ...any) (sql.Result, error) {
	return c.config.ExecContext(ctx, query, args...)
}

// Query runs a raw SQL query via Ent’s underlying driver.
func (c *Client) Query(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	return c.config.QueryContext(ctx, query, args...)
}
