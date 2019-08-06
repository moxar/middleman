package sqlm

import (
	"context"
	"database/sql"
)

// Middleware for sql transaction.
type Middleware interface {
	Get(next GetFn) GetFn
	Select(next SelectFn) SelectFn
	Exec(next ExecFn) ExecFn
}

// GetFn is a middleware function to tx.Get.
type GetFn = func(ctx context.Context, dest interface{}, query string, args ...interface{}) error

// ExecFn is a middleware function to tx.Exec.
type ExecFn = func(ctx context.Context, query string, args ...interface{}) (sql.Result, error)

// SelectFn is a middleware function to tx.Select.
type SelectFn = func(ctx context.Context, dest interface{}, query string, args ...interface{}) error

// Queryer describes a type that can run queries.
type Queryer interface {
	GetContext(context.Context, interface{}, string, ...interface{}) error
	SelectContext(context.Context, interface{}, string, ...interface{}) error
	ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
}

// Chain the middlewares into a new queryer.
func Chain(q Queryer, ms ...Middleware) Queryer {
	var out queryer
	out.GetFn = q.GetContext
	out.SelectFn = q.SelectContext
	out.ExecFn = q.ExecContext
	for i := len(ms) - 1; i >= 0; i-- {
		out.GetFn = ms[i].Get(out.GetFn)
		out.SelectFn = ms[i].Select(out.SelectFn)
		out.ExecFn = ms[i].Exec(out.ExecFn)
	}
	return &out
}

type queryer struct {
	GetFn
	SelectFn
	ExecFn
}

func (q *queryer) GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return q.GetFn(ctx, dest, query, args...)
}

func (q *queryer) SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return q.SelectFn(ctx, dest, query, args...)
}

func (q *queryer) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return q.ExecFn(ctx, query, args...)
}
