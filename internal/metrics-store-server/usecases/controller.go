package usecases

import (
	"github.com/aridae/go-metrics-store/internal/metrics-store-server/usecases/counter"
	"github.com/aridae/go-metrics-store/internal/metrics-store-server/usecases/gauge"
)

type Controller struct {
	counterUseCasesHandler *counter.Handler
	gaugeUseCasesHandler   *gauge.Handler
}

func NewController(
	counterUseCasesHandler *counter.Handler,
	gaugeUseCasesHandler *gauge.Handler,
) *Controller {
	return &Controller{
		counterUseCasesHandler: counterUseCasesHandler,
		gaugeUseCasesHandler:   gaugeUseCasesHandler,
	}
}
