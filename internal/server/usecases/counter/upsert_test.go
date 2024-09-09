package counter

import (
	"context"
	"github.com/aridae/go-metrics-store/internal/server/models"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"testing"
	"time"
)

func TestHandler_Upsert_HappyCase_Increment(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	ctrl := gomock.NewController(t)
	now := time.Date(2024, 01, 02, 03, 40, 50, 0, time.UTC)

	metricName := "testMetricName"
	metricPrevValue := models.CounterValue(123)
	metricIncrement := models.CounterValue(456)
	expectedNewValue := models.CounterValue(579)

	metricUpdater := models.ScalarMetricUpdater{
		Type:  models.ScalarMetricTypeCounter,
		Name:  metricName,
		Value: metricIncrement,
	}

	expectedPrevState := models.ScalarMetric{
		ScalarMetricUpdater: models.ScalarMetricUpdater{
			Type:  models.ScalarMetricTypeCounter,
			Name:  metricName,
			Value: metricPrevValue,
		},
	}

	expectedNewState := models.ScalarMetric{
		ScalarMetricUpdater: models.ScalarMetricUpdater{
			Type:  models.ScalarMetricTypeCounter,
			Name:  metricName,
			Value: expectedNewValue,
		},
		Datetime: now,
	}

	fixture := setupFixture(ctrl, now)

	_ = metricsRepo.GetLatestState
	fixture.metricsRepo.EXPECT().
		GetLatestState(ctx, metricUpdater.Type, metricUpdater.Name).
		Return(&expectedPrevState, nil)

	_ = metricsRepo.Save
	fixture.metricsRepo.EXPECT().
		Save(ctx, expectedNewState).
		Return(nil)

	err := fixture.handler.Upsert(ctx, metricUpdater)
	require.NoError(t, err)
}

func TestHandler_Upsert_HappyCase_NoPrevValue(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	ctrl := gomock.NewController(t)
	now := time.Date(2024, 01, 02, 03, 40, 50, 0, time.UTC)

	metricName := "testMetricName"
	metricIncrement := models.CounterValue(456)
	expectedNewValue := models.CounterValue(456)

	metricUpdater := models.ScalarMetricUpdater{
		Type:  models.ScalarMetricTypeCounter,
		Name:  metricName,
		Value: metricIncrement,
	}

	expectedNewState := models.ScalarMetric{
		ScalarMetricUpdater: models.ScalarMetricUpdater{
			Type:  models.ScalarMetricTypeCounter,
			Name:  metricName,
			Value: expectedNewValue,
		},
		Datetime: now,
	}

	fixture := setupFixture(ctrl, now)

	_ = metricsRepo.GetLatestState
	fixture.metricsRepo.EXPECT().
		GetLatestState(ctx, metricUpdater.Type, metricUpdater.Name).
		Return(nil, nil)

	_ = metricsRepo.Save
	fixture.metricsRepo.EXPECT().
		Save(ctx, expectedNewState).
		Return(nil)

	err := fixture.handler.Upsert(ctx, metricUpdater)
	require.NoError(t, err)
}
