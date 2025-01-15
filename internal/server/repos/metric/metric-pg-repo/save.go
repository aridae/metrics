package metricpgrepo

import (
	"context"
	"fmt"

	"github.com/aridae/go-metrics-store/internal/server/models"
)

func (r *repo) Save(ctx context.Context, metric models.Metric) error {
	queryable := r.txGetter.DefaultTrOrDB(ctx, r.db)

	onConflict := fmt.Sprintf("ON CONFLICT(%s) DO UPDATE SET %s = EXCLUDED.%s, %s = EXCLUDED.%s;",
		keyColumn, valueColumn, valueColumn, datetimeColumn, datetimeColumn)

	qb := psql.Insert(metricTable).
		Columns(
			keyColumn,
			typeColumn,
			nameColumn,
			valueColumn,
			datetimeColumn,
		).
		Values(
			metric.GetKey().String(),
			metric.GetType().String(),
			metric.GetName(),
			metric.GetValue().String(),
			metric.Datetime,
		).
		Suffix(onConflict)

	sql, args, err := qb.ToSql()
	if err != nil {
		return fmt.Errorf("failed to build query: %w", err)
	}

	_, err = queryable.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("failed to run Exec: %w", err)
	}

	return nil
}
