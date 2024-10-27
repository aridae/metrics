package pgtxman

import (
	"context"
	"errors"
	"fmt"
	"github.com/aridae/go-metrics-store/internal/server/repos"
	"github.com/aridae/go-metrics-store/internal/server/repos/pg-driven-repos/metric-pg-repo"
	"github.com/aridae/go-metrics-store/pkg/logger"
	"github.com/jackc/pgx/v5"
)

type db interface {
	Begin(ctx context.Context) (pgx.Tx, error)
}

type txManager struct {
	db db
}

func NewTransactionManagerImplementation(db db) repos.TransactionManager {
	return &txManager{db: db}
}

func (txman *txManager) DoInTransaction(ctx context.Context, fn func(adapters *repos.Repositories) error) error {
	tx, err := txman.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}

	defer func() {
		rbErr := tx.Rollback(ctx)
		if rbErr != nil && !errors.Is(rbErr, pgx.ErrTxClosed) {
			logger.Obtain().Errorf("transaction rollback failed: %v", rbErr)
		}
	}()

	adapters, err := initAdapters(ctx, tx)
	if err != nil {
		return fmt.Errorf("failed to transactional repos adapters: %w", err)
	}

	err = fn(adapters)
	if err != nil {
		return fmt.Errorf("failed to execute transactional function: %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func initAdapters(ctx context.Context, tx pgx.Tx) (*repos.Repositories, error) {
	metricRepo, err := metricpgrepo.NewRepositoryImplementation(ctx, tx)
	if err != nil {
		return nil, fmt.Errorf("failed to init metric repository: %w", err)
	}

	adapters := repos.Repositories{
		MetricRepository: metricRepo,
	}

	return &adapters, nil
}
