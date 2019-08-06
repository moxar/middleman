# sqlm

[![Godoc](http://img.shields.io/badge/godoc-reference-blue.svg?style=flat)](https://godoc.org/github.com/moxar/middleman/sqlm)

`sqlm` is the middleman subpackage specialized in sql. It provides constructors for sql middlewares, based on jmoiron/sqlx.

Here is a simple example that logs the queries.

```go
// Logger is a middleware that logs queries.
type Logger struct{
	log.Logger
}

func (l *Logger) Get(next sqlm.GetFn) sqlm.GetFn {
	return func(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
		l.Log(query)
		return next(ctx, dest, query, args...)
	}
}

// Exec logs the Exec query.
func (l *Logger) Exec(next sqlm.ExecFn) sqlm.ExecFn {
	return func(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
		l.Log(query)
		return next(ctx, query, args...)
	}
}

// Select logs the Select query.
func (l *Logger) Select(next sqlm.SelectFn) sqlm.SelectFn {
	return func(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
		l.Log(query)
		return next(ctx, dest, query, args...)
	}
}

// Will log the query
sqlm.Chain(db, Logger{}).SelectContext(ctx, "INSERT INTO users VALUES(?, ?)", "Batman", "DCU")
```
