package nooptrm

import (
	"context"
	"github.com/aridae/go-metrics-store/pkg/logger"
	"github.com/avito-tech/go-transaction-manager/trm/v2"
)

type Manager struct{}

func NewNoopTransactionManager() *Manager {
	return &Manager{}
}

func (*Manager) Do(ctx context.Context, fn func(ctx context.Context) error) (err error) {
	logger.Infof("transactional manager is no-op, does not support transactional operations")

	err = fn(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (m *Manager) DoWithSettings(ctx context.Context, _ trm.Settings, fn func(ctx context.Context) error) error {
	return m.Do(ctx, fn)
}
