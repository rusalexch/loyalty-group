package db

import (
	"context"
	"errors"
	"log"
	"sync"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rusalexch/loyalty-group/internal/accrual/internal/common"
)

type database struct {
	pool *pgxpool.Pool
	mx   sync.Mutex
}

var matchStmt = "match"

func New(dbURL string) *database {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		log.Panic(err)
	}
	db := &database{
		pool: pool,
	}
	ctx, cancel = context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	err = db.init(ctx)
	if err != nil {
		log.Panic(err)
	}

	return db
}

func (db *database) init(ctx context.Context) error {
	db.mx.Lock()
	defer db.mx.Unlock()

	_, err := db.pool.Exec(ctx, sqlCreateRewards)
	return err
}

func (db *database) Add(ctx context.Context, reward common.Reward) error {
	db.mx.Lock()
	defer db.mx.Unlock()

	_, err := db.pool.Exec(ctx, sqlAddReward, reward.ID, reward.Type, reward.Reward)
	return err
}

func (db *database) Match(ctx context.Context, match []string) (map[string]common.Reward, error) {
	db.mx.Lock()
	defer db.mx.Unlock()

	trx, err := db.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer trx.Commit(ctx)

	stmt, err := trx.Prepare(ctx, matchStmt, sqlFindRewards)
	if err != nil {
		trx.Rollback(ctx)
		return nil, err
	}
	res := make(map[string]common.Reward)
	for _, s := range match {
		r := reward{}
		row := trx.QueryRow(ctx, stmt.Name, s)
		err := row.Scan(r.ID, r.Type, r.Reward)
		if err != nil && !errors.Is(err, pgx.ErrNoRows) {
			trx.Rollback(ctx)
			return nil, err
		}
		res[s] = convert(r)
	}

	return res, nil
}

func convert(r reward) common.Reward {
	return common.Reward{
		ID:     r.ID,
		Type:   common.RewardType(r.Type),
		Reward: r.Reward,
	}
}
