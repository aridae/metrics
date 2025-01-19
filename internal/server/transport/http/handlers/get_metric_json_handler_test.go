package handlers

import (
	"bytes"
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

func Test_getMetricJSONHandler(t *testing.T) {
	t.Parallel()

	type prereq struct {
		httpMethod  string
		urlEndpoint string
		requestBody []byte

		mockExpectedMetricKey *models.MetricKey
		mockReturnedMetric    *models.Metric
		mockControllerErr     error
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
				urlEndpoint: "/value",
			},
			want: want{
				httpCode: http.StatusMethodNotAllowed,
			},
		},
		{
			desc: "negative: invalid method: PUT",
			prereq: prereq{
				httpMethod:  http.MethodPut,
				urlEndpoint: "/value",
			},
			want: want{
				httpCode: http.StatusMethodNotAllowed,
			},
		},
		{
			desc: "negative: invalid method: PATCH",
			prereq: prereq{
				httpMethod:  http.MethodPatch,
				urlEndpoint: "/value",
			},
			want: want{
				httpCode: http.StatusMethodNotAllowed,
			},
		},
		{
			desc: "negative: invalid request body: name absent",
			prereq: prereq{
				httpMethod:  http.MethodPost,
				urlEndpoint: "/value",
				requestBody: []byte(`{"type":"counter"}`),
			},
			want: want{
				httpCode: http.StatusBadRequest,
			},
		},
		{
			desc: "negative: invalid request body: type absent",
			prereq: prereq{
				httpMethod:  http.MethodPost,
				urlEndpoint: "/value",
				requestBody: []byte(`{"id":"testName"}`),
			},
			want: want{
				httpCode: http.StatusBadRequest,
			},
		},
		{
			desc: "negative: invalid request body: type unknown",
			prereq: prereq{
				httpMethod:  http.MethodPost,
				urlEndpoint: "/value",
				requestBody: []byte(`{"id":"testName", "type":"unknown"}`),
			},
			want: want{
				httpCode: http.StatusBadRequest,
			},
		},
		{
			desc: "negative: counter not found",
			prereq: prereq{
				httpMethod:            http.MethodPost,
				urlEndpoint:           "/value",
				requestBody:           []byte(`{"id":"testName", "type":"counter"}`),
				mockExpectedMetricKey: pointer.To[models.MetricKey]("counter:testName"),
			},
			want: want{
				httpCode: http.StatusNotFound,
			},
		},
		{
			desc: "negative: gauge not found",
			prereq: prereq{
				httpMethod:            http.MethodPost,
				urlEndpoint:           "/value",
				requestBody:           []byte(`{"id":"testName", "type":"gauge"}`),
				mockExpectedMetricKey: pointer.To[models.MetricKey]("gauge:testName"),
			},
			want: want{
				httpCode: http.StatusNotFound,
			},
		},
		{
			desc: "positive: counter found",
			prereq: prereq{
				httpMethod:            http.MethodPost,
				urlEndpoint:           "/value",
				requestBody:           []byte(`{"id":"testName", "type":"counter"}`),
				mockExpectedMetricKey: pointer.To[models.MetricKey]("counter:testName"),
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
				httpMethod:            http.MethodPost,
				urlEndpoint:           "/value",
				requestBody:           []byte(`{"id":"testName", "type":"gauge"}`),
				mockExpectedMetricKey: pointer.To[models.MetricKey]("gauge:testName"),
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

			if test.prereq.mockExpectedMetricKey != nil {
				controllerMock.EXPECT().
					GetMetricByKey(gomock.Any(), *test.prereq.mockExpectedMetricKey).
					Return(test.prereq.mockReturnedMetric, test.prereq.mockControllerErr)
			}

			w := httptest.NewRecorder()

			router := NewRouter(controllerMock)
			router.getMetricJSONHandler(w, req)

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
