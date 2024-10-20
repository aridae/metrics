package metricpgrepo

import (
	"context"
	"fmt"
	"github.com/aridae/go-metrics-store/internal/server/repos/metric-repo"
	"github.com/aridae/go-metrics-store/pkg/postgres"
)

type repo struct {
	db *postgres.Client
}

func NewRepositoryImplementation(ctx context.Context, pgClient *postgres.Client) (metricrepo.Repository, error) {
	imp := &repo{db: pgClient}

	err := imp.prepareSchema(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare schema: %w", err)
	}

	return imp, nil
}

func (r *repo) prepareSchema(ctx context.Context) error {
	_, err := r.db.Exec(ctx, schemaDDL)
	if err != nil {
		return fmt.Errorf("executing schema ddl: %w", err)
	}

	return nil
}
