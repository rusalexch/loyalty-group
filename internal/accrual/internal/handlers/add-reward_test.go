package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/rusalexch/loyalty-group/internal/accrual/internal/app"
	mock "github.com/rusalexch/loyalty-group/internal/accrual/internal/handlers/mocks"
	"github.com/stretchr/testify/assert"
)

func Test_handlers_addReward(t *testing.T) {
	type args struct {
		reward       app.Reward
		addRewardErr error
		requestBody  string
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
			name: "add reward: OK",
			args: args{
				reward: app.Reward{
					ID:     "Bork",
					Reward: 10,
					Type:   app.Percentage,
				},
				addRewardErr: nil,
				requestBody:  `{"match":"Bork","reward":10,"reward_type":"%"}`,
			},
			want: want{
				code: 200,
			},
		},
		{
			name: "add reward: already exist",
			args: args{
				reward: app.Reward{
					ID:     "Bork",
					Reward: 10,
					Type:   app.Percentage,
				},
				addRewardErr: app.ErrRewardAlreadyExist,
				requestBody:  `{"match":"Bork","reward":10,"reward_type":"%"}`,
			},
			want: want{
				code: 409,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctl := gomock.NewController(t)
			defer ctl.Finish()

			srv := mock.NewMockservice(ctl)
			srv.EXPECT().
				AddReward(gomock.Any(), tt.args.reward).
				Return(tt.args.addRewardErr).
				AnyTimes()

			h := New(srv)
			ts := httptest.NewServer(h)
			defer ts.Close()

			statusCode, _ := testRequest(t, ts, http.MethodPost, "/api/goods", strings.NewReader(tt.args.requestBody))
			assert.Equal(t, tt.want.code, statusCode)
		})
	}
}
