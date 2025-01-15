package usecases

import (
	"context"
	"github.com/aridae/go-metrics-store/internal/server/models"
	"github.com/aridae/go-metrics-store/internal/server/usecases/_mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"testing"
	"time"
)

func Test_upsertMetricOverride_Happy_Gauge_Float64MetricValue(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	ctrl := gomock.NewController(t)
	metricsRepoMock := _mock.NewMockmetricsRepo(ctrl)
	now := time.Date(2025, 9, 03, 15, 30, 24, 0, time.UTC)

	metricUpsert := models.MetricUpsert{
		MName: "TestMetricName",
		Mtype: models.MetricTypeGauge,
		Val:   models.NewFloat64MetricValue(1234),
	}
	expectedMetricState := models.Metric{
		MetricUpsert: metricUpsert,
		Datetime:     now,
	}

	metricsRepoMock.EXPECT().Save(ctx, expectedMetricState).Return(nil)

	actualMetricState, err := upsertMetricOverride(ctx, metricsRepoMock, metricUpsert, now)

	require.NoError(t, err)
	require.EqualValues(t, expectedMetricState, actualMetricState)
}

func Test_upsertMetricIncrement_Happy_Counter_Int64MetricValue_PrevValueExists(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	ctrl := gomock.NewController(t)
	metricsRepoMock := _mock.NewMockmetricsRepo(ctrl)
	now := time.Date(2025, 9, 03, 15, 30, 24, 0, time.UTC)

	prevMetricState := models.Metric{
		MetricUpsert: models.MetricUpsert{
			MName: "TestMetricName",
			Mtype: models.MetricTypeCounter,
			Val:   models.NewInt64MetricValue(2345),
		},
		Datetime: now.Add(-time.Hour * 78),
	}
	metricUpsert := models.MetricUpsert{
		MName: "TestMetricName",
		Mtype: models.MetricTypeCounter,
		Val:   models.NewInt64MetricValue(1234),
	}
	expectedMetricState := models.Metric{
		MetricUpsert: models.MetricUpsert{
			MName: "TestMetricName",
			Mtype: models.MetricTypeCounter,
			Val:   models.NewInt64MetricValue(3579),
		},
		Datetime: now,
	}

	metricsRepoMock.EXPECT().GetByKey(ctx, metricUpsert.GetKey()).Return(&prevMetricState, nil)
	metricsRepoMock.EXPECT().Save(ctx, expectedMetricState).Return(nil)

	actualMetricState, err := upsertMetricIncrement(ctx, metricsRepoMock, metricUpsert, now)

	require.NoError(t, err)
	require.EqualValues(t, expectedMetricState, actualMetricState)
}

func Test_upsertMetricIncrement_Happy_Counter_Int64MetricValue_PrevValueNotFound(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	ctrl := gomock.NewController(t)
	metricsRepoMock := _mock.NewMockmetricsRepo(ctrl)
	now := time.Date(2025, 9, 03, 15, 30, 24, 0, time.UTC)

	metricUpsert := models.MetricUpsert{
		MName: "TestMetricName",
		Mtype: models.MetricTypeCounter,
		Val:   models.NewInt64MetricValue(1234),
	}
	expectedMetricState := models.Metric{
		MetricUpsert: metricUpsert,
		Datetime:     now,
	}

	metricsRepoMock.EXPECT().GetByKey(ctx, metricUpsert.GetKey()).Return(nil, nil)
	metricsRepoMock.EXPECT().Save(ctx, expectedMetricState).Return(nil)

	actualMetricState, err := upsertMetricIncrement(ctx, metricsRepoMock, metricUpsert, now)

	require.NoError(t, err)
	require.EqualValues(t, expectedMetricState, actualMetricState)
}
