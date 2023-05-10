package reward

import (
	"context"

	"github.com/rusalexch/loyalty-group/internal/accrual/internal/common"
)

type rewardRepo interface {
	Add(ctx context.Context, reward common.Reward) error
	Find(ctx context.Context, description string) (common.Reward, error)
}

type rewardServicve struct {
	repo rewardRepo
}

func New(repo rewardRepo) *rewardServicve {
	return &rewardServicve{
		repo: repo,
	}
}
