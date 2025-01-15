package usecases

import (
	"context"
	"github.com/aridae/go-metrics-store/internal/server/models"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestController_UpsertMetric_Happy_MetricTypeGauge(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	now, _ := time.Parse(time.DateTime, "2024-01-02 15:04:05")

	inputMetricUpsert := models.MetricUpsert{
		MName: "TestMetricName",
		Mtype: models.MetricTypeGauge,
		Val:   models.NewFloat64MetricValue(123.456),
	}
	expectedMetricNewState := models.Metric{
		MetricUpsert: inputMetricUpsert,
		Datetime:     now,
	}

	tk := setupTestKit(t, now)

	tk.metricsRepoMock.EXPECT().
		Save(ctx, models.Metric{MetricUpsert: inputMetricUpsert, Datetime: now}).
		Return(nil)

	actualMetricNewState, err := tk.controller.UpsertMetric(ctx, inputMetricUpsert)

	require.NoError(t, err)
	require.EqualValues(t, expectedMetricNewState, actualMetricNewState)
}

func TestController_UpsertMetric_Happy_MetricTypeCounter(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	now, _ := time.Parse(time.DateTime, "2024-01-02 15:04:05")

	metricName := "TestMetricName"
	prevMetricValue := int64(1234)
	metricIncrement := int64(3456)

	inputMetricUpsert := models.MetricUpsert{
		MName: metricName,
		Mtype: models.MetricTypeCounter,
		Val:   models.NewInt64MetricValue(metricIncrement),
	}
	prevMetricState := models.Metric{
		MetricUpsert: models.MetricUpsert{
			MName: metricName,
			Mtype: models.MetricTypeCounter,
			Val:   models.NewInt64MetricValue(prevMetricValue),
		},
	}
	expectedMetricNewState := models.Metric{
		MetricUpsert: models.MetricUpsert{
			MName: metricName,
			Mtype: models.MetricTypeCounter,
			Val:   models.NewInt64MetricValue(prevMetricValue + metricIncrement),
		},
		Datetime: now,
	}

	tk := setupTestKit(t, now)

	tk.metricsRepoMock.EXPECT().
		GetByKey(ctx, inputMetricUpsert.GetKey()).
		Return(&prevMetricState, nil)

	tk.metricsRepoMock.EXPECT().
		Save(ctx, expectedMetricNewState).
		Return(nil)

	actualMetricNewState, err := tk.controller.UpsertMetric(ctx, inputMetricUpsert)
	require.NoError(t, err)
	require.EqualValues(t, expectedMetricNewState, actualMetricNewState)
}
