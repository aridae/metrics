package metricpgrepo

import (
	"context"
	"fmt"
	"github.com/aridae/go-metrics-store/internal/server/repos"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"sync"
)

var (
	_once sync.Once
)

type database interface {
	Exec(ctx context.Context, sql string, arguments ...any) (commandTag pgconn.CommandTag, err error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

type repo struct {
	db database
}

func NewRepositoryImplementation(ctx context.Context, db database) (repos.MetricRepository, error) {
	imp := &repo{db: db}

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
