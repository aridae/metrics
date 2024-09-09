package gauge

import (
	"context"
	"github.com/aridae/go-metrics-store/internal/server/models"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"testing"
	"time"
)

func TestHandler_Upsert_HappyCase(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	ctrl := gomock.NewController(t)
	now := time.Date(2024, 01, 02, 03, 40, 50, 0, time.UTC)

	metricName := "testMetricName"
	metricNewVal := models.GaugeValue(456.543)

	metricUpdater := models.ScalarMetricUpdater{
		Type:  models.ScalarMetricTypeGauge,
		Name:  metricName,
		Value: metricNewVal,
	}

	expectedNewValue := models.ScalarMetric{
		ScalarMetricUpdater: metricUpdater,
		Datetime:            now,
	}

	fixture := setupFixture(ctrl, now)

	_ = metricsRepo.Save
	fixture.metricsRepo.EXPECT().
		Save(ctx, expectedNewValue).
		Return(nil)

	err := fixture.handler.Upsert(ctx, metricUpdater)
	require.NoError(t, err)
}
