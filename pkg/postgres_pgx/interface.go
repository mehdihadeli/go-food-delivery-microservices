package postgres

// Ref:https://github.com/henvic/pgxtutorial/blob/668784624474abea3619433c6e45510f4d156649/internal/database/interface.go
import (
	"context"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

// PGX limited interface with high-level API for pgx methods safe to be used in high-level business logic packages.
// It is satisfied by implementations *pgx.Conn and *pgxpool.Pool (and you should probably use the second one usually).
//
// Caveat: It doesn't expose a method to acquire a *pgx.Conn or handle notifications,
// so it's not compatible with LISTEN/NOTIFY.
//
// Reference: https://pkg.go.dev/github.com/jackc/pgx/v4
type PGX interface {
	// BeginTx starts a transaction with txOptions determining the transaction mode. Unlike database/sql, the context only
	// affects the begin command. i.e. there is no auto-rollback on context cancellation.
	BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error)

	// BeginTxFunc starts a transaction with txOptions determining the transaction mode and calls f. If f does not return
	// an error the transaction is committed. If f returns an error the transaction is rolled back. The context will be
	// used when executing the transaction control statements (BEGIN, ROLLBACK, and COMMIT) but does not otherwise affect
	// the execution of f.
	BeginTxFunc(ctx context.Context, txOptions pgx.TxOptions, f func(pgx.Tx) error) error

	PGXQuerier
}

// PGXQuerier interface with methods used for everything, including transactions.
type PGXQuerier interface {
	// Begin starts a transaction. Unlike database/sql, the context only affects the begin command. i.e. there is no
	// auto-rollback on context cancellation.
	Begin(ctx context.Context) (pgx.Tx, error)

	// BeginFunc starts a transaction and calls f. If f does not return an error the transaction is committed. If f returns
	// an error the transaction is rolled back. The context will be used when executing the transaction control statements
	// (BEGIN, ROLLBACK, and COMMIT) but does not otherwise affect the execution of f.
	BeginFunc(ctx context.Context, f func(pgx.Tx) error) error

	// CopyFrom uses the PostgreSQL copy protocol to perform bulk data insertion.
	// It returns the number of rows copied and an error.
	//
	// CopyFrom requires all values use the binary format. Almost all types
	// implemented by pgx use the binary format by default. Types implementing
	// Encoder can only be used if they encode to the binary format.
	CopyFrom(ctx context.Context, tableName pgx.Identifier, columnNames []string, rowSrc pgx.CopyFromSource) (int64, error)

	// Exec executes sql. sql can be either a prepared statement name or an SQL string. arguments should be referenced
	// positionally from the sql string as $1, $2, etc.
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)

	// Query executes sql with args. If there is an error the returned Rows will be returned in an error state. So it is
	// allowed to ignore the error returned from Query and handle it in Rows.
	//
	// For extra control over how the query is executed, the types QuerySimpleProtocol, QueryResultFormats, and
	// QueryResultFormatsByOID may be used as the first args to control exactly how the query is executed. This is rarely
	// needed. See the documentation for those types for details.
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)

	// QueryFunc executes sql with args. For each row returned by the query the values will scanned into the elements of
	// scans and f will be called. If any row fails to scan or f returns an error the query will be aborted and the error
	// will be returned.
	QueryFunc(ctx context.Context, sql string, args []any, scans []any, f func(pgx.QueryFuncRow) error) (pgconn.CommandTag, error)

	// QueryRow is a convenience wrapper over Query. Any error that occurs while
	// querying is deferred until calling Scan on the returned Row. That Row will
	// error with ErrNoRows if no rows are returned.
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row

	// SendBatch sends all queued queries to the server at once. All queries are run in an implicit transaction unless
	// explicit transaction control statements are executed. The returned BatchResults must be closed before the connection
	// is used again.
	SendBatch(ctx context.Context, b *pgx.Batch) pgx.BatchResults
}

// Validate if the PGX interface was derived from *pgx.Conn and *pgxpool.Pool correctly.
var (
	_ PGX = (*pgx.Conn)(nil)
	_ PGX = (*pgxpool.Pool)(nil)
)
