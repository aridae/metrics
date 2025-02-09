package usecases

import (
	"testing"
	"time"

	"github.com/aridae/go-metrics-store/internal/server/usecases/_mock"
	"github.com/aridae/go-metrics-store/pkg/noop-trm"
	"go.uber.org/mock/gomock"
)

type testKit struct {
	metricsRepoMock *_mock.MockmetricsRepo
	controller      *Controller
}

func setupTestKit(t *testing.T, now time.Time) *testKit {
	ctrl := gomock.NewController(t)
	metricsRepoMock := _mock.NewMockmetricsRepo(ctrl)
	transactionManagerMock := nooptrm.NewNoopTransactionManager()

	controller := &Controller{
		metricsRepo:        metricsRepoMock,
		transactionManager: transactionManagerMock,
		now:                func() time.Time { return now },
	}

	return &testKit{
		metricsRepoMock: metricsRepoMock,
		controller:      controller,
	}
}
