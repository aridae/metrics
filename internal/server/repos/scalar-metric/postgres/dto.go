package postgres

import (
	"fmt"
	"time"

	metricsfabrics "github.com/aridae/go-metrics-store/internal/server/metrics-fabrics"
	"github.com/aridae/go-metrics-store/internal/server/models"
)

type dto struct {
	ID       int64     `db:"id"`
	Name     string    `db:"name"`
	MType    string    `db:"type"`
	Value    string    `db:"value"`
	Datetime time.Time `db:"datetime"`
}

func parseDTOs(dtos []dto) ([]models.ScalarMetric, error) {
	metrics := make([]models.ScalarMetric, 0, len(dtos))
	for _, d := range dtos {
		metric, err := parseDTO(d)
		if err != nil {
			return nil, err
		}

		metrics = append(metrics, metric)
	}

	return metrics, nil
}

func parseDTO(d dto) (models.ScalarMetric, error) {
	mtype := models.ScalarMetricType(d.MType)

	parser, ok := metricsValueParsers[mtype]
	if !ok {
		return models.ScalarMetric{}, fmt.Errorf("unsupported metrics type %q", mtype)
	}

	mvalue, err := parser(d.Value)
	if err != nil {
		return models.ScalarMetric{}, fmt.Errorf("parsing value %q: %w", d.Value, err)
	}

	return models.ScalarMetric{
		ScalarMetricToRegister: models.ScalarMetricToRegister{
			MName: d.Name,
			Mtype: mtype,
			Val:   mvalue,
		},
		Datetime: d.Datetime,
	}, nil
}

var metricsValueParsers = map[models.ScalarMetricType]func(str string) (models.ScalarMetricValue, error){
	models.ScalarMetricTypeCounter: metricsfabrics.ObtainCounterMetricFactory().ParseScalarMetricValue,
	models.ScalarMetricTypeGauge:   metricsfabrics.ObtainGaugeMetricFactory().ParseScalarMetricValue,
}
