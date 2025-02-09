package metricpgrepo

import (
	"fmt"
	"github.com/aridae/go-metrics-store/internal/server/models/factories"
	"time"

	"github.com/aridae/go-metrics-store/internal/server/models"
)

type metricDTO struct {
	Datetime time.Time `db:"datetime"`
	Name     string    `db:"name"`
	MType    string    `db:"type"`
	Value    string    `db:"value"`
	ID       int64     `db:"id"`
}

func parseDTOs(dtos []metricDTO) ([]models.Metric, error) {
	metrics := make([]models.Metric, 0, len(dtos))
	for _, d := range dtos {
		metric, err := parseDTO(d)
		if err != nil {
			return nil, err
		}

		metrics = append(metrics, metric)
	}

	return metrics, nil
}

func parseDTO(d metricDTO) (models.Metric, error) {
	mtype := models.MetricType(d.MType)

	parser, ok := metricsValueParsers[mtype]
	if !ok {
		return models.Metric{}, fmt.Errorf("unsupported metrics-reporting type %q", mtype)
	}

	mvalue, err := parser(d.Value)
	if err != nil {
		return models.Metric{}, fmt.Errorf("parsing value %q: %w", d.Value, err)
	}

	return models.Metric{
		MetricUpsert: models.MetricUpsert{
			MName: d.Name,
			Mtype: mtype,
			Val:   mvalue,
		},
		Datetime: d.Datetime,
	}, nil
}

var metricsValueParsers = map[models.MetricType]func(str string) (models.MetricValue, error){
	models.MetricTypeCounter: factories.ObtainCounterMetricFactory().ParseMetricValue,
	models.MetricTypeGauge:   factories.ObtainGaugeMetricFactory().ParseMetricValue,
}
