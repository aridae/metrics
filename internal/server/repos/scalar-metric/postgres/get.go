package postgres

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/aridae/go-metrics-store/internal/server/models"
	"github.com/georgysavva/scany/v2/pgxscan"
)

func (r *repo) GetLatestState(ctx context.Context, key models.MetricKey) (*models.ScalarMetric, error) {
	qb := baseSelectQuery.Where(squirrel.Eq{keyColumn: key.String()})

	sql, args, err := qb.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	rows, err := r.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query: %w", err)
	}

	var dtos []dto
	err = pgxscan.ScanAll(&dtos, rows)
	if err != nil {
		return nil, fmt.Errorf("failed to scan row into dto: %w", err)
	}

	if len(dtos) == 0 {
		return nil, nil
	}

	d := dtos[0]
	metric, err := parseDTO(d)
	if err != nil {
		return nil, fmt.Errorf("failed to parse dto into metric business model: %w", err)
	}

	return &metric, nil
}

func (r *repo) GetAllLatestStates(ctx context.Context) ([]models.ScalarMetric, error) {
	qb := baseSelectQuery.OrderBy(keyColumn)

	sql, args, err := qb.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	rows, err := r.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query: %w", err)
	}

	var dtos []dto
	err = pgxscan.ScanAll(&dtos, rows)
	if err != nil {
		return nil, fmt.Errorf("failed to scan rows into dtos: %w", err)
	}

	if len(dtos) == 0 {
		return nil, nil
	}

	metrics, err := parseDTOs(dtos)
	if err != nil {
		return nil, fmt.Errorf("failed to parse dtos into metric business models: %w", err)
	}

	return metrics, nil
}
