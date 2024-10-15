package usecases

import (
	"context"
	"fmt"

	"github.com/aridae/go-metrics-store/internal/server/logger"
)

func (c *Controller) Healthcheck(ctx context.Context) error {
	logger.Obtain().Debugf("usecases.Healthcheck checking if controller's dependencies are available")

	repoHealthcheckErr := c.metricsRepo.Healthcheck(ctx)
	if repoHealthcheckErr != nil {
		return fmt.Errorf("Controller.metricsRepo is unavailable: %w", repoHealthcheckErr)
	}

	return nil
}
