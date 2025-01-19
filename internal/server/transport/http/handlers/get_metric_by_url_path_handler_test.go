package handlers

import (
	"context"
	"github.com/aridae/go-metrics-store/internal/server/models"
	"github.com/aridae/go-metrics-store/internal/server/transport/http/handlers/_mock"
	"github.com/aridae/go-metrics-store/pkg/pointer"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_getMetricByURLPathHandler(t *testing.T) {
	t.Parallel()

	type prereq struct {
		httpMethod  string
		urlEndpoint string
		chiParams   map[string]string

		expectedMetricKey  *models.MetricKey
		mockReturnedMetric *models.Metric
		mockControllerErr  error
	}

	type want struct {
		expectedReturnValue string
		httpCode            int
	}

	testCases := []struct {
		desc   string
		prereq prereq
		want   want
	}{
		{
			desc: "negative: invalid method: POST",
			prereq: prereq{
				httpMethod:  http.MethodPost,
				urlEndpoint: "/value/counter/testName",
			},
			want: want{
				httpCode: http.StatusMethodNotAllowed,
			},
		},
		{
			desc: "negative: invalid method: PUT",
			prereq: prereq{
				httpMethod:  http.MethodPut,
				urlEndpoint: "/value/counter/testName",
			},
			want: want{
				httpCode: http.StatusMethodNotAllowed,
			},
		},
		{
			desc: "negative: invalid method: PATCH",
			prereq: prereq{
				httpMethod:  http.MethodPatch,
				urlEndpoint: "/value/counter/testName",
			},
			want: want{
				httpCode: http.StatusMethodNotAllowed,
			},
		},
		{
			desc: "negative: invalid url: name absent",
			prereq: prereq{
				httpMethod:  http.MethodGet,
				urlEndpoint: "/value/counter",
			},
			want: want{
				httpCode: http.StatusBadRequest,
			},
		},
		{
			desc: "negative: invalid url: type unknown",
			prereq: prereq{
				httpMethod:  http.MethodGet,
				urlEndpoint: "/value/unknown/testName/",
				chiParams: map[string]string{
					urlParamMetricType: "unknown",
					urlParamMetricName: "testName",
				},
			},
			want: want{
				httpCode: http.StatusBadRequest,
			},
		},
		{
			desc: "negative: counter not found",
			prereq: prereq{
				httpMethod:  http.MethodGet,
				urlEndpoint: "/value/counter/testName/",
				chiParams: map[string]string{
					urlParamMetricType: "counter",
					urlParamMetricName: "testName",
				},
				expectedMetricKey: pointer.To[models.MetricKey]("counter:testName"),
			},
			want: want{
				httpCode: http.StatusNotFound,
			},
		},
		{
			desc: "negative: gauge not found",
			prereq: prereq{
				httpMethod:  http.MethodGet,
				urlEndpoint: "/value/gauge/testName/",
				chiParams: map[string]string{
					urlParamMetricType: "gauge",
					urlParamMetricName: "testName",
				},
				expectedMetricKey: pointer.To[models.MetricKey]("gauge:testName"),
			},
			want: want{
				httpCode: http.StatusNotFound,
			},
		},
		{
			desc: "positive: counter found",
			prereq: prereq{
				httpMethod:  http.MethodGet,
				urlEndpoint: "/value/counter/testName/",
				chiParams: map[string]string{
					urlParamMetricType: "counter",
					urlParamMetricName: "testName",
				},
				expectedMetricKey: pointer.To[models.MetricKey]("counter:testName"),
				mockReturnedMetric: &models.Metric{
					MetricUpsert: models.MetricUpsert{
						MName: "testName",
						Mtype: models.MetricTypeCounter,
						Val:   models.NewInt64MetricValue(123),
					},
				},
			},
			want: want{
				expectedReturnValue: "123",
				httpCode:            http.StatusOK,
			},
		},
		{
			desc: "positive: counter found",
			prereq: prereq{
				httpMethod:  http.MethodGet,
				urlEndpoint: "/value/gauge/testName/",
				chiParams: map[string]string{
					urlParamMetricType: "gauge",
					urlParamMetricName: "testName",
				},
				expectedMetricKey: pointer.To[models.MetricKey]("gauge:testName"),
				mockReturnedMetric: &models.Metric{
					MetricUpsert: models.MetricUpsert{
						MName: "testName",
						Mtype: models.MetricTypeGauge,
						Val:   models.NewFloat64MetricValue(123.5),
					},
				},
			},
			want: want{
				expectedReturnValue: "123.5",
				httpCode:            http.StatusOK,
			},
		},
	}

	for _, test := range testCases {
		t.Run(test.desc, func(t *testing.T) {
			rctx := chi.NewRouteContext()
			for k, v := range test.prereq.chiParams {
				rctx.URLParams.Add(k, v)
			}

			req := httptest.NewRequest(test.prereq.httpMethod, test.prereq.urlEndpoint, nil)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			ctrl := gomock.NewController(t)
			controllerMock := _mock.NewMockuseCasesController(ctrl)
			if test.prereq.expectedMetricKey != nil {
				controllerMock.EXPECT().
					GetMetricByKey(gomock.Any(), *test.prereq.expectedMetricKey).
					Return(test.prereq.mockReturnedMetric, test.prereq.mockControllerErr)
			}

			w := httptest.NewRecorder()

			router := NewRouter(controllerMock)
			router.getMetricByURLPathHandler(w, req)

			resp := w.Result()

			body, _ := io.ReadAll(resp.Body)
			_ = resp.Body.Close()

			if len(test.want.expectedReturnValue) > 0 {
				require.EqualValues(t, test.want.expectedReturnValue, body)
			}
			require.Equal(t, test.want.httpCode, resp.StatusCode)
		})
	}
}
