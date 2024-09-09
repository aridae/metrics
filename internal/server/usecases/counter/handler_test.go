package counter

import (
	"github.com/aridae/go-metrics-store/internal/server/usecases/counter/_mock"
	"go.uber.org/mock/gomock"
	"time"
)

type handlerFixture struct {
	metricsRepo *_mock.MockmetricsRepo
	handler     *Handler
}

func setupFixture(ctrl *gomock.Controller, now time.Time) *handlerFixture {
	metricsRepoMock := _mock.NewMockmetricsRepo(ctrl)
	handler := NewHandler(metricsRepoMock)
	handler.now = func() time.Time { return now }

	return &handlerFixture{
		metricsRepo: metricsRepoMock,
		handler:     handler,
	}
}
