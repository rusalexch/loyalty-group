package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/rusalexch/loyalty-group/internal/accrual/internal/app"
	mock "github.com/rusalexch/loyalty-group/internal/accrual/internal/service/mocks"
	"github.com/stretchr/testify/assert"
)

func Test_service_Ping(t *testing.T) {

	tests := []struct {
		name     string
		err      error
		behavior func(s *mock.Mockstorager, ctx context.Context, err error)
		wantErr  error
	}{
		{
			name: "ping ok",
			err:  nil,
			behavior: func(s *mock.Mockstorager, ctx context.Context, err error) {
				s.EXPECT().Ping(ctx).Return(err)
			},
			wantErr: nil,
		},
		{
			name: "ping fail",
			err:  errors.New("test error"),
			behavior: func(s *mock.Mockstorager, ctx context.Context, err error) {
				s.EXPECT().Ping(ctx).Return(err)
			},
			wantErr: errors.New("test error"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctl := gomock.NewController(t)
			defer ctl.Finish()

			ctx, close := context.WithTimeout(context.Background(), 5*time.Second)
			defer close()
			store := mock.NewMockstorager(ctl)
			tt.behavior(store, ctx, tt.err)

			srv := New(ServiceConfig{
				Store: store,
			})

			err := srv.Ping(ctx)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func Test_service_GetOrder(t *testing.T) {
	type args struct {
		orderID string
		order   app.Order
		err     error
	}
	type want struct {
		order app.Order
		err   error
	}

	tests := []struct {
		name     string
		args     args
		behavior func(s *mock.MockorderRepository, ctx context.Context, args args)
		want     want
	}{
		{
			name: "return registered order",
			args: args{
				orderID: "1234",
				order: app.Order{
					ID:     "1234",
					Status: app.Registered,
				},
				err: nil,
			},
			behavior: func(s *mock.MockorderRepository, ctx context.Context, args args) {
				s.EXPECT().FindByID(ctx, args.orderID).Return(args.order, args.err)
			},
			want: want{
				order: app.Order{
					ID:      "1234",
					Status:  app.Registered,
					Accrual: nil,
				},
				err: nil,
			},
		},
		{
			name: "return not found",
			args: args{
				orderID: "1234",
				order:   app.Order{},
				err:     app.ErrOrderNotFound,
			},
			behavior: func(s *mock.MockorderRepository, ctx context.Context, args args) {
				s.EXPECT().FindByID(ctx, args.orderID).Return(args.order, args.err)
			},
			want: want{
				order: app.Order{},
				err:   app.ErrOrderNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctl := gomock.NewController(t)
			defer ctl.Finish()

			ctx, close := context.WithTimeout(context.Background(), 5*time.Second)
			defer close()

			stor := mock.NewMockorderRepository(ctl)
			tt.behavior(stor, ctx, tt.args)

			srv := New(ServiceConfig{
				OrderRepo: stor,
			})

			order, err := srv.GetOrder(ctx, tt.args.orderID)

			assert.Equal(t, tt.want.err, err)
			assert.Equal(t, tt.want.order, order)
		})
	}
}

func Test_service_isRewardExist(t *testing.T) {
	type args struct {
		rewardID string
		reward   app.Reward
		err      error
	}
	type want struct {
		isExist bool
		err     error
	}

	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "not exist",
			args: args{
				rewardID: "testReward",
				reward:   app.Reward{},
				err:      app.ErrRewardNotFound,
			},
			want: want{
				isExist: false,
				err:     nil,
			},
		},
		{
			name: "exist",
			args: args{
				rewardID: "testReward",
				reward: app.Reward{
					ID:     "testReward",
					Reward: 10,
					Type:   app.Fixed,
				},
				err: nil,
			},
			want: want{
				isExist: true,
				err:     nil,
			},
		},
		{
			name: "some error",
			args: args{
				rewardID: "testReward",
				reward:   app.Reward{},
				err:      errors.New("test error"),
			},
			want: want{
				isExist: false,
				err:     errors.New("test error"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctl := gomock.NewController(t)
			defer ctl.Finish()
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			stor := mock.NewMockrewardRepository(ctl)
			stor.EXPECT().FindByID(ctx, tt.args.rewardID).Return(tt.args.reward, tt.args.err)
			srv := New(ServiceConfig{
				RewardRepo: stor,
			})

			got, err := srv.isRewardExist(ctx, tt.args.rewardID)
			assert.Equal(t, tt.want.err, err)
			assert.Equal(t, tt.want.isExist, got)
		})
	}
}

func Test_service_AddReward(t *testing.T) {
	type args struct {
		reward  app.Reward
		addErr  error
		findErr error
		exist   app.Reward
	}

	tests := []struct {
		name string
		args args
		want error
	}{
		{
			name: "success added reward",
			args: args{
				reward: app.Reward{
					ID:     "testReward",
					Type:   app.Percentage,
					Reward: 10,
				},
				exist:   app.Reward{},
				addErr:  nil,
				findErr: app.ErrRewardNotFound,
			},
			want: nil,
		},
		{
			name: "reward exist",
			args: args{
				reward: app.Reward{
					ID:     "testReward",
					Type:   app.Percentage,
					Reward: 10,
				},
				exist: app.Reward{
					ID:     "testReward",
					Type:   app.Fixed,
					Reward: 5,
				},
				addErr:  app.ErrRewardAlreadyExist,
				findErr: nil,
			},
			want: app.ErrRewardAlreadyExist,
		},
		{
			name: "error",
			args: args{
				reward: app.Reward{
					ID:     "testReward",
					Type:   app.Percentage,
					Reward: 10,
				},
				exist:   app.Reward{},
				addErr:  nil,
				findErr: errors.New("test error"),
			},
			want: errors.New("test error"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctr := gomock.NewController(t)
			defer ctr.Finish()
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			stor := mock.NewMockrewardRepository(ctr)
			stor.EXPECT().FindByID(ctx, tt.args.reward.ID).Return(tt.args.exist, tt.args.findErr)
			stor.EXPECT().Add(ctx, tt.args.reward).Return(tt.args.addErr).AnyTimes()
			srv := New(ServiceConfig{
				RewardRepo: stor,
			})

			err := srv.AddReward(ctx, tt.args.reward)
			assert.Equal(t, tt.want, err)
		})
	}
}

func Test_service_AddOrder(t *testing.T) {
	type args struct {
		order         app.OrderGoods
		addOrderErr   error
		addProductErr error
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name: "add order",
			args: args{
				order: app.OrderGoods{
					ID: "1234",
					Goods: []app.OrderProduct{
						{
							Description: "some product",
							Price:       100,
						},
						{
							Description: "another product",
							Price:       300,
						},
					},
				},
				addOrderErr:   nil,
				addProductErr: nil,
			},
			wantErr: nil,
		},
		{
			name: "can't create order",
			args: args{
				order: app.OrderGoods{
					ID: "1234",
					Goods: []app.OrderProduct{
						{
							Description: "some product",
							Price:       100,
						},
						{
							Description: "another product",
							Price:       300,
						},
					},
				},
				addOrderErr:   errors.New("can't create order"),
				addProductErr: nil,
			},
			wantErr: errors.New("can't create order"),
		},
		{
			name: "can't add goods",
			args: args{
				order: app.OrderGoods{
					ID: "1234",
					Goods: []app.OrderProduct{
						{
							Description: "some product",
							Price:       100,
						},
						{
							Description: "another product",
							Price:       300,
						},
					},
				},
				addOrderErr:   nil,
				addProductErr: errors.New("can't add goods"),
			},
			wantErr: errors.New("can't add goods"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctl := gomock.NewController(t)
			defer ctl.Finish()
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			orderStor := mock.NewMockorderRepository(ctl)
			prodStor := mock.NewMockproductRepository(ctl)
			orderStor.EXPECT().Add(ctx, tt.args.order.ID).Return(tt.args.addOrderErr).AnyTimes()
			orderStor.EXPECT().Delete(ctx, tt.args.order.ID).Return(nil).AnyTimes()
			prodStor.EXPECT().Add(ctx, tt.args.order.ID, tt.args.order.Goods).Return(tt.args.addProductErr).AnyTimes()

			srv := New(ServiceConfig{
				OrderRepo:   orderStor,
				ProductRepo: prodStor,
			})

			got := srv.AddOrder(ctx, tt.args.order)
			assert.Equal(t, tt.wantErr, got)
		})
	}
}
