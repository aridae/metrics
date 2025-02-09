package handlers

import (
	"bytes"
	"context"
	"github.com/aridae/go-metrics-store/internal/server/models"
	"github.com/aridae/go-metrics-store/internal/server/transport/http/handlers/_mock"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_updateMetricsBatchJSONHandler(t *testing.T) {
	t.Parallel()

	type prereq struct {
		mockControllerErr          error
		httpMethod                 string
		urlEndpoint                string
		requestBody                []byte
		mockExpectedMetricsUpserts []models.MetricUpsert
		mockReturnedMetrics        []models.Metric
	}

	type want struct {
		responseBody []byte
		httpCode     int
	}

	testCases := []struct {
		desc   string
		prereq prereq
		want   want
	}{
		{
			desc: "negative: invalid method: GET",
			prereq: prereq{
				httpMethod:  http.MethodGet,
				urlEndpoint: "/updates/",
			},
			want: want{
				httpCode: http.StatusMethodNotAllowed,
			},
		},
		{
			desc: "negative: invalid method: PUT",
			prereq: prereq{
				httpMethod:  http.MethodPut,
				urlEndpoint: "/updates/",
			},
			want: want{
				httpCode: http.StatusMethodNotAllowed,
			},
		},
		{
			desc: "negative: invalid method: PATCH",
			prereq: prereq{
				httpMethod:  http.MethodPatch,
				urlEndpoint: "/updates/",
			},
			want: want{
				httpCode: http.StatusMethodNotAllowed,
			},
		},
		{
			desc: "negative: invalid request body: one of names absent",
			prereq: prereq{
				httpMethod:  http.MethodPost,
				urlEndpoint: "/updates/",
				requestBody: []byte(`[{"type":"counter", "delta":123}, {"type":"gauge", "name":"testGauge", "value":123.5}]`),
			},
			want: want{
				httpCode: http.StatusBadRequest,
			},
		},
		{
			desc: "negative: invalid request body: one of types absent",
			prereq: prereq{
				httpMethod:  http.MethodPost,
				urlEndpoint: "/updates/",
				requestBody: []byte(`[{"name":"testCounter", "delta":123}, {"type":"gauge", "name":"testGauge", "value":123.5}]`),
			},
			want: want{
				httpCode: http.StatusBadRequest,
			},
		},
		{
			desc: "negative: invalid request body: one of values absent",
			prereq: prereq{
				httpMethod:  http.MethodPost,
				urlEndpoint: "/updates/",
				requestBody: []byte(`[{"type":"counter", "name":"testCounter"}, {"type":"gauge", "name":"testGauge", "value":123.5}]`),
			},
			want: want{
				httpCode: http.StatusBadRequest,
			},
		},
		{
			desc: "negative: invalid request body: one of types unknow",
			prereq: prereq{
				httpMethod:  http.MethodPost,
				urlEndpoint: "/updates/",
				requestBody: []byte(`[{"type":"counter", "name":"testCounter", "delta":123}, {"type":"unknown", "name":"testGauge", "value":123.5}]`),
			},
			want: want{
				httpCode: http.StatusBadRequest,
			},
		},
		{
			desc: "negative: invalid request body: wrong counter updater",
			prereq: prereq{
				httpMethod:  http.MethodPost,
				urlEndpoint: "/updates/",
				requestBody: []byte(`[{"type":"counter", "name":"testCounter", "value":123}, {"type":"gauge", "name":"testGauge", "value":123.5}]`),
			},
			want: want{
				httpCode: http.StatusBadRequest,
			},
		},
		{
			desc: "negative: invalid request body: wrong gauge updater",
			prereq: prereq{
				httpMethod:  http.MethodPost,
				urlEndpoint: "/updates/",
				requestBody: []byte(`[{"type":"counter", "name":"testCounter", "delta":123}, {"type":"gauge", "name":"testGauge", "delta":123.5}]`),
			},
			want: want{
				httpCode: http.StatusBadRequest,
			},
		},
		{
			desc: "positive: single counter",
			prereq: prereq{
				httpMethod:  http.MethodPost,
				urlEndpoint: "/updates/",
				requestBody: []byte(`[{"id":"testName", "type":"counter", "delta":123}]`),
				mockExpectedMetricsUpserts: []models.MetricUpsert{
					{
						MName: "testName",
						Mtype: models.MetricTypeCounter,
						Val:   models.NewInt64MetricValue(123),
					},
				},
				mockReturnedMetrics: []models.Metric{
					{
						MetricUpsert: models.MetricUpsert{
							MName: "testName",
							Mtype: models.MetricTypeCounter,
							Val:   models.NewInt64MetricValue(123),
						},
					},
				},
			},
			want: want{
				responseBody: []byte(`[{"id":"testName", "type":"counter", "delta":123}]`),
				httpCode:     http.StatusOK,
			},
		},
		{
			desc: "positive: single gauge",
			prereq: prereq{
				httpMethod:  http.MethodPost,
				urlEndpoint: "/updates/",
				requestBody: []byte(`[{"id":"testName", "type":"gauge", "value":123.5}]`),
				mockExpectedMetricsUpserts: []models.MetricUpsert{
					{
						MName: "testName",
						Mtype: models.MetricTypeGauge,
						Val:   models.NewFloat64MetricValue(123.5),
					},
				},
				mockReturnedMetrics: []models.Metric{
					{
						MetricUpsert: models.MetricUpsert{
							MName: "testName",
							Mtype: models.MetricTypeGauge,
							Val:   models.NewFloat64MetricValue(123.5),
						},
					},
				},
			},
			want: want{
				responseBody: []byte(`[{"id":"testName", "type":"gauge", "value":123.5}]`),
				httpCode:     http.StatusOK,
			},
		},
		{
			desc: "positive: counter and gauge",
			prereq: prereq{
				httpMethod:  http.MethodPost,
				urlEndpoint: "/updates/",
				requestBody: []byte(`[{"id":"testGauge", "type":"gauge", "value":123.5},{"id":"testCounter", "type":"counter", "delta":123}]`),
				mockExpectedMetricsUpserts: []models.MetricUpsert{
					{
						MName: "testGauge",
						Mtype: models.MetricTypeGauge,
						Val:   models.NewFloat64MetricValue(123.5),
					},
					{
						MName: "testCounter",
						Mtype: models.MetricTypeCounter,
						Val:   models.NewInt64MetricValue(123),
					},
				},
				mockReturnedMetrics: []models.Metric{
					{
						MetricUpsert: models.MetricUpsert{
							MName: "testGauge",
							Mtype: models.MetricTypeGauge,
							Val:   models.NewFloat64MetricValue(123.5),
						},
					},
					{
						MetricUpsert: models.MetricUpsert{
							MName: "testCounter",
							Mtype: models.MetricTypeCounter,
							Val:   models.NewInt64MetricValue(234),
						},
					},
				},
			},
			want: want{
				responseBody: []byte(`[{"id":"testGauge", "type":"gauge", "value":123.5},{"id":"testCounter", "type":"counter", "delta":234}]`),
				httpCode:     http.StatusOK,
			},
		},
	}

	for _, test := range testCases {
		t.Run(test.desc, func(t *testing.T) {
			rctx := chi.NewRouteContext()

			req := httptest.NewRequest(test.prereq.httpMethod, test.prereq.urlEndpoint, bytes.NewBuffer(test.prereq.requestBody))
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			ctrl := gomock.NewController(t)
			controllerMock := _mock.NewMockuseCasesController(ctrl)

			if test.prereq.mockExpectedMetricsUpserts != nil {
				controllerMock.EXPECT().
					UpsertMetricsBatch(gomock.Any(), test.prereq.mockExpectedMetricsUpserts).
					Return(test.prereq.mockReturnedMetrics, test.prereq.mockControllerErr)
			}

			w := httptest.NewRecorder()

			router := NewRouter(controllerMock)
			router.updateMetricsBatchJSONHandler(w, req)

			resp := w.Result()

			body, _ := io.ReadAll(resp.Body)
			_ = resp.Body.Close()

			if len(test.want.responseBody) > 0 {
				require.JSONEq(t, string(test.want.responseBody), string(body))
			}
			require.Equal(t, test.want.httpCode, resp.StatusCode)
		})
	}
}
