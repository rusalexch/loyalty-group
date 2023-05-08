package core

import (
	"context"
	"testing"

	"github.com/rusalexch/loyalty-group/internal/accrual/internal/common"
	"github.com/rusalexch/loyalty-group/internal/accrual/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_calc(t *testing.T) {
	type args struct {
		price float64
		acc   common.Reward
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			name: "percentage accrual",
			args: args{
				price: 30.00,
				acc: common.Reward{
					ID:     "test",
					Reward: 3,
					Type:   common.Percentage,
				},
			},
			want: 0.9,
		},
		{
			name: "fix accrual",
			args: args{
				price: 100.00,
				acc: common.Reward{
					ID:     "test",
					Reward: 13,
					Type:   common.Fixed,
				},
			},
			want: 13.00,
		},
		{
			name: "undefined type",
			args: args{
				price: 100.00,
				acc: common.Reward{
					ID:     "test",
					Reward: 13,
					Type:   common.RewardType(3),
				},
			},
			want: 0.00,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := calc(tt.args.price, tt.args.acc)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_core_Calc(t *testing.T) {
	type fields struct {
		store storager
		goods []common.Reward
	}
	type args struct {
		goods []common.OrderProduct
	}
	type want struct {
		price float64
	}

	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "goods calc",
			fields: fields{
				store: mocks.MockAccrualStoragerNew(),
				goods: []common.Reward{
					{
						ID:     "Bosh",
						Type:   common.Fixed,
						Reward: 3,
					},
					{
						ID:     "LG",
						Type:   common.Percentage,
						Reward: 10,
					},
					{
						ID:     "Samsung",
						Type:   common.Fixed,
						Reward: 5,
					},
				},
			},
			args: args{
				goods: []common.OrderProduct{
					{
						Name:  "Bosh",
						Price: 100,
					},
					{
						Name:  "Bosh",
						Price: 1000,
					},
					{
						Name:  "LG",
						Price: 50,
					},
					{
						Name:  "Apple",
						Price: 5000,
					},
				},
			},
			want: want{
				price: 11,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			c := New(tt.fields.store)

			for _, v := range tt.fields.goods {
				c.Add(ctx, v)
			}

			got, err := c.Calc(ctx, tt.args.goods)
			require.NoError(t, err)
			assert.Equal(t, tt.want.price, got)

		})
	}
}
