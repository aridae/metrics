package repos

import (
	"context"
	"github.com/aridae/go-metrics-store/internal/server/logger"
)

type Repositories struct {
	MetricRepository
}

type TransactionManager interface {
	DoInTransaction(context.Context, func(*Repositories) error) error
}

// noopTxMan for all my fellow storages, which don't support transactional operations.
type noopTxMan struct {
	repos *Repositories
}

func NewNoopTransactionManager(repos *Repositories) TransactionManager {
	return &noopTxMan{repos: repos}
}

func (txman *noopTxMan) DoInTransaction(_ context.Context, fn func(*Repositories) error) error {
	logger.Obtain().Info("transactional manager is no-op, does not support transactional operations")

	err := fn(txman.repos)
	if err != nil {
		return err
	}

	return nil
}
