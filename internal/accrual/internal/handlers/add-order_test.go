package handlers

import (
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/rusalexch/loyalty-group/internal/accrual/internal/app"
	mock "github.com/rusalexch/loyalty-group/internal/accrual/internal/handlers/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testRequest(t *testing.T, ts *httptest.Server, method, path string, body io.Reader) (int, string) {
	req, err := http.NewRequest(method, ts.URL+path, body)
	require.NoError(t, err)
	if body != nil {
		req.Header.Add(contentType, appJSON)
	}
	log.Println(req.Context())

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	defer resp.Body.Close()

	return resp.StatusCode, string(respBody)
}

func Test_handlers_addOrder(t *testing.T) {
	type args struct {
		order       app.OrderGoods
		existOrder  app.Order
		getOrderErr error
		addOrderErr error
		requestBody string
	}
	type want struct {
		code int
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "add order: OK",
			args: args{
				order: app.OrderGoods{
					ID: 4561261212345467,
					Goods: []app.OrderProduct{
						{
							Description: "Чайник Bork",
							Price:       7000,
						},
					},
				},
				existOrder:  app.Order{},
				getOrderErr: app.ErrOrderNotFound,
				addOrderErr: nil,
				requestBody: `{"order":4561261212345467,"goods":[{"description":"Чайник Bork","price":7000}]}`,
			},
			want: want{
				code: 202,
			},
		},
		{
			name: "add order: bad request",
			args: args{
				order: app.OrderGoods{
					ID: 4561261212345467,
					Goods: []app.OrderProduct{
						{
							Description: "Чайник Bork",
							Price:       7000,
						},
					},
				},
				existOrder: app.Order{
					ID:      4561261212345467,
					Status:  app.Registered,
					Accrual: nil,
				},
				getOrderErr: nil,
				addOrderErr: nil,
				requestBody: `{"order":4561261212345467,"goods":[{"description":"Чайник Bork","price":7000}]}`,
			},
			want: want{
				code: 409,
			},
		},
		{
			name: "add order: already exist",
			args: args{
				order: app.OrderGoods{
					ID: 3,
					Goods: []app.OrderProduct{
						{
							Description: "Чайник Bork",
							Price:       7000,
						},
					},
				},
				existOrder:  app.Order{},
				getOrderErr: app.ErrOrderNotFound,
				addOrderErr: nil,
				requestBody: `{"order":3,"goods":[{"description":"Чайник Bork","price":7000}]}`,
			},
			want: want{
				code: 400,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctl := gomock.NewController(t)
			defer ctl.Finish()

			srv := mock.NewMockservice(ctl)
			srv.EXPECT().
				GetOrder(gomock.Any(), tt.args.order.ID).
				Return(tt.args.existOrder, tt.args.getOrderErr).
				AnyTimes()
			srv.EXPECT().
				AddOrder(gomock.Any(), tt.args.order).
				Return(tt.args.addOrderErr).
				AnyTimes()

			h := New(srv)
			ts := httptest.NewServer(h)
			defer ts.Close()

			statusCode, _ := testRequest(t, ts, http.MethodPost, "/api/orders", strings.NewReader(tt.args.requestBody))
			assert.Equal(t, tt.want.code, statusCode)
		})
	}
}
