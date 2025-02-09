package metricpgrepo

import (
	"context"
	"fmt"
	"sync"

	"github.com/aridae/go-metrics-store/internal/server/repos/metric"
	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

var (
	_once sync.Once
)

type sqlQueryable interface {
	Begin(ctx context.Context) (pgx.Tx, error)
	CopyFrom(ctx context.Context, tableName pgx.Identifier, columnNames []string, rowSrc pgx.CopyFromSource) (int64, error)
	SendBatch(ctx context.Context, b *pgx.Batch) pgx.BatchResults
	Exec(ctx context.Context, sql string, arguments ...interface{}) (commandTag pgconn.CommandTag, err error)
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
}

type sqlTransactionManager interface {
	DefaultTrOrDB(ctx context.Context, db trmpgx.Tr) trmpgx.Tr
}

type repo struct {
	db       sqlQueryable
	txGetter sqlTransactionManager
}

func NewRepositoryImplementation(
	ctx context.Context,
	db sqlQueryable,
	txGetter sqlTransactionManager,
) (metric.Repository, error) {
	imp := &repo{db: db, txGetter: txGetter}

	err := imp.prepareSchema(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare schema: %w", err)
	}

	return imp, nil
}

func (r *repo) prepareSchema(ctx context.Context) error {
	var err error
	_once.Do(func() {
		_, err = r.db.Exec(ctx, schemaDDL)
	})
	if err != nil {
		return fmt.Errorf("executing schema ddl: %w", err)
	}

	return nil
}
