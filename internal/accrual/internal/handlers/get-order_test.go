package handlers

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/rusalexch/loyalty-group/internal/accrual/internal/app"
	mock "github.com/rusalexch/loyalty-group/internal/accrual/internal/handlers/mocks"
	"github.com/rusalexch/loyalty-group/internal/utils"
	"github.com/stretchr/testify/assert"
)

func Test_handlers_getOrder(t *testing.T) {
	type args struct {
		orderID     string
		order       app.Order
		getOrderErr error
	}
	type want struct {
		code  int
		order string
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "get order: Registered",
			args: args{
				orderID: "4561261212345467",
				order: app.Order{
					ID:      "4561261212345467",
					Status:  app.Registered,
					Accrual: nil,
				},
				getOrderErr: nil,
			},
			want: want{
				code:  200,
				order: `{"order":4561261212345467,"status":"REGISTERED"}`,
			},
		},
		{
			name: "get order: processed",
			args: args{
				orderID: "4561261212345467",
				order: app.Order{
					ID:      "4561261212345467",
					Status:  app.Processed,
					Accrual: utils.Float64ToPointer(500),
				},
				getOrderErr: nil,
			},
			want: want{
				code:  200,
				order: `{"order":4561261212345467,"status":"PROCESSED","accrual":500}`,
			},
		},
		{
			name: "get order: processed",
			args: args{
				orderID:     "13",
				order:       app.Order{},
				getOrderErr: nil,
			},
			want: want{
				code:  400,
				order: ``,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctl := gomock.NewController(t)
			defer ctl.Finish()

			srv := mock.NewMockservice(ctl)
			srv.EXPECT().GetOrder(gomock.Any(), tt.args.orderID).Return(tt.args.order, tt.args.getOrderErr).AnyTimes()

			h := New(srv)
			ts := httptest.NewServer(h)
			defer ts.Close()

			statusCode, body := testRequest(t, ts, http.MethodGet, fmt.Sprintf("/api/orders/%s", tt.args.orderID), nil)
			assert.Equal(t, tt.want.code, statusCode)
			assert.Equal(t, tt.want.order, body)
		})
	}
}
