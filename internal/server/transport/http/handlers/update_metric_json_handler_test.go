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

func Test_updateMetricJSONHandler(t *testing.T) {
	t.Parallel()

	type prereq struct {
		mockControllerErr        error
		mockExpectedMetricUpsert *models.MetricUpsert
		mockReturnedMetric       *models.Metric
		httpMethod               string
		urlEndpoint              string
		requestBody              []byte
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
				urlEndpoint: "/update",
			},
			want: want{
				httpCode: http.StatusMethodNotAllowed,
			},
		},
		{
			desc: "negative: invalid method: PUT",
			prereq: prereq{
				httpMethod:  http.MethodPut,
				urlEndpoint: "/update",
			},
			want: want{
				httpCode: http.StatusMethodNotAllowed,
			},
		},
		{
			desc: "negative: invalid method: PATCH",
			prereq: prereq{
				httpMethod:  http.MethodPatch,
				urlEndpoint: "/update",
			},
			want: want{
				httpCode: http.StatusMethodNotAllowed,
			},
		},
		{
			desc: "negative: invalid request body: name absent",
			prereq: prereq{
				httpMethod:  http.MethodPost,
				urlEndpoint: "/update",
				requestBody: []byte(`{"type":"counter", "delta":123}`),
			},
			want: want{
				httpCode: http.StatusBadRequest,
			},
		},
		{
			desc: "negative: invalid request body: type absent",
			prereq: prereq{
				httpMethod:  http.MethodPost,
				urlEndpoint: "/update",
				requestBody: []byte(`{"id":"testName", "delta":123}`),
			},
			want: want{
				httpCode: http.StatusBadRequest,
			},
		},
		{
			desc: "negative: invalid request body: value absent",
			prereq: prereq{
				httpMethod:  http.MethodPost,
				urlEndpoint: "/update",
				requestBody: []byte(`{"id":"testName", "type":"counter"}`),
			},
			want: want{
				httpCode: http.StatusBadRequest,
			},
		},
		{
			desc: "negative: invalid request body: type unknown",
			prereq: prereq{
				httpMethod:  http.MethodPost,
				urlEndpoint: "/update",
				requestBody: []byte(`{"id":"testName", "type":"unknown", "delta":123}`),
			},
			want: want{
				httpCode: http.StatusBadRequest,
			},
		},
		{
			desc: "negative: invalid request body: wrong counter updater",
			prereq: prereq{
				httpMethod:  http.MethodPost,
				urlEndpoint: "/update",
				requestBody: []byte(`{"id":"testName", "type":"counter", "value":123}`),
			},
			want: want{
				httpCode: http.StatusBadRequest,
			},
		},
		{
			desc: "negative: invalid request body: wrong gauge updater",
			prereq: prereq{
				httpMethod:  http.MethodPost,
				urlEndpoint: "/update",
				requestBody: []byte(`{"id":"testName", "type":"gauge", "delta":123}`),
			},
			want: want{
				httpCode: http.StatusBadRequest,
			},
		},
		{
			desc: "positive: counter",
			prereq: prereq{
				httpMethod:  http.MethodPost,
				urlEndpoint: "/update",
				requestBody: []byte(`{"id":"testName", "type":"counter", "delta":123}`),
				mockExpectedMetricUpsert: &models.MetricUpsert{
					MName: "testName",
					Mtype: models.MetricTypeCounter,
					Val:   models.NewInt64MetricValue(123),
				},
				mockReturnedMetric: &models.Metric{
					MetricUpsert: models.MetricUpsert{
						MName: "testName",
						Mtype: models.MetricTypeCounter,
						Val:   models.NewInt64MetricValue(123),
					},
				},
			},
			want: want{
				responseBody: []byte(`{"id":"testName", "type":"counter", "delta":123}`),
				httpCode:     http.StatusOK,
			},
		},
		{
			desc: "positive: gauge found",
			prereq: prereq{
				httpMethod:  http.MethodPost,
				urlEndpoint: "/update",
				requestBody: []byte(`{"id":"testName", "type":"gauge", "value":123.5}`),
				mockExpectedMetricUpsert: &models.MetricUpsert{
					MName: "testName",
					Mtype: models.MetricTypeGauge,
					Val:   models.NewFloat64MetricValue(123.5),
				},
				mockReturnedMetric: &models.Metric{
					MetricUpsert: models.MetricUpsert{
						MName: "testName",
						Mtype: models.MetricTypeGauge,
						Val:   models.NewFloat64MetricValue(123.5),
					},
				},
			},
			want: want{
				responseBody: []byte(`{"id":"testName", "type":"gauge", "value":123.5}`),
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

			if test.prereq.mockExpectedMetricUpsert != nil {
				controllerMock.EXPECT().
					UpsertMetric(gomock.Any(), *test.prereq.mockExpectedMetricUpsert).
					Return(*test.prereq.mockReturnedMetric, test.prereq.mockControllerErr)
			}

			w := httptest.NewRecorder()

			router := NewRouter(controllerMock)
			router.updateMetricJSONHandler(w, req)

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
