package reward

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rusalexch/loyalty-group/internal/accrual/internal/app"
)

type reward struct {
	ID     string  `db:"id"`
	Type   string  `db:"type"`
	Reward float64 `db:"reward"`
}

type rewardRepository struct {
	pool *pgxpool.Pool
}

// New конструктор репозитория начислений
func New(pool *pgxpool.Pool) *rewardRepository {
	repo := &rewardRepository{
		pool: pool,
	}

	err := repo.init()
	if err != nil {
		log.Println("can't init reward repository")
		log.Panic(err)
	}

	return repo
}

// init инициализация репозитория начислений
func (repo *rewardRepository) init() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := repo.pool.Exec(ctx, sqlCreateRewards)

	return err
}

// Add добавление нового начисления
func (repo *rewardRepository) Add(ctx context.Context, reward app.Reward) error {
	_, err := repo.pool.Exec(ctx, sqlAddReward, reward.ID, reward.Type, reward.Reward)

	return err
}

// Find поиск начисления по наименованию товара
func (repo *rewardRepository) Find(ctx context.Context, description string) (app.Reward, error) {
	var reward reward
	row := repo.pool.QueryRow(ctx, sqlFindRewards, description)
	err := row.Scan(&reward.ID, &reward.Type, &reward.Reward)
	if err != nil && errors.Is(err, pgx.ErrNoRows) {
		return app.Reward{}, app.ErrRewardNotFound
	}

	return dbToJSON(reward), err
}

// FindByID поиск начисления по наименованию
func (repo *rewardRepository) FindByID(ctx context.Context, ID string) (app.Reward, error) {
	var reward reward
	row := repo.pool.QueryRow(ctx, sqlFindByID, ID)
	err := row.Scan(&reward)
	if err != nil && errors.Is(err, pgx.ErrNoRows) {
		return app.Reward{}, app.ErrRewardNotFound
	}

	return dbToJSON(reward), err
}

// dbToJSON преобразование структуры БД в структуру JSON
func dbToJSON(reward reward) app.Reward {
	return app.Reward{
		ID:     reward.ID,
		Type:   reward.Type,
		Reward: reward.Reward,
	}
}
