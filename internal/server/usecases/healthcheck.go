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

	pgHealthcheckErr := c.postgresConn.Healthcheck(ctx)
	if pgHealthcheckErr != nil {
		return fmt.Errorf("Controller.postgresConn is unavailable: %w", pgHealthcheckErr)
	}

	return nil
}
