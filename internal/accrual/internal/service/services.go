package service

import (
	"context"
	"errors"
	"log"
	"math"
	"time"

	"github.com/rusalexch/loyalty-group/internal/accrual/internal/app"
)

const (
	defaultTick = 1 * time.Second
	maxTick     = 5 * time.Second
)

type ServiceConfig struct {
	Store       storager
	OrderRepo   orderRepository
	ProductRepo productRepository
	RewardRepo  rewardRepository
}

type service struct {
	store       storager
	orderRepo   orderRepository
	productRepo productRepository
	rewardRepo  rewardRepository
	tick        time.Duration
	ticker      *time.Ticker
	calcChan    chan string
}

// New конструктор сервиса
func New(conf ServiceConfig) *service {
	s := &service{
		store:       conf.Store,
		orderRepo:   conf.OrderRepo,
		productRepo: conf.ProductRepo,
		rewardRepo:  conf.RewardRepo,
		tick:        defaultTick,
		ticker:      time.NewTicker(defaultTick),
		calcChan:    make(chan string, 10),
	}

	go s.run()

	return s
}

// Ping проверка работоспособности
func (s *service) Ping(ctx context.Context) error {
	return s.store.Ping(ctx)
}

// GetOrder получение данных по заказу
func (s *service) GetOrder(ctx context.Context, orderID string) (app.Order, error) {
	return s.orderRepo.FindByID(ctx, orderID)
}

// AddReward добавление/изменение схемы начисления
func (s *service) AddReward(ctx context.Context, reward app.Reward) error {
	isExist, err := s.isRewardExist(ctx, reward.ID)
	if err != nil {
		return err
	}
	if isExist {
		return app.ErrRewardAlreadyExist
	}

	return s.rewardRepo.Add(ctx, reward)
}

func (s *service) isRewardExist(ctx context.Context, rewardID string) (bool, error) {
	reward, err := s.rewardRepo.FindByID(ctx, rewardID)
	if err == nil && reward.ID == rewardID {
		return true, nil
	}
	if errors.Is(err, app.ErrRewardNotFound) {
		return false, nil
	}

	return false, err
}

// AddOrder добавление нового заказа
func (s *service) AddOrder(ctx context.Context, order app.OrderGoods) error {
	err := s.orderRepo.Add(ctx, order.ID)
	if err != nil {
		return err
	}

	err = s.productRepo.Add(ctx, order.ID, order.Goods)
	if err != nil {
		s.orderRepo.Delete(ctx, order.ID)
	}
	return err
}

// destroy метод завершения сервиса
func (s *service) Stop() {
	s.ticker.Stop()
	close(s.calcChan)
}

// run метод запуска расчетов заказов
func (s *service) run() {
	go func() {
		for orderID := range s.calcChan {
			s.calculate(orderID)
		}
	}()

	go func() {
		for range s.ticker.C {
			s.findOrderToCalc()
		}
	}()
}

// resetTicker сброс таймера до исходного значения
func (s *service) resetTicker() {
	s.tick = defaultTick
	s.ticker.Reset(defaultTick)
}

// upTiceker увеличение таймера
func (s *service) upTiceker() {
	s.tick += 1 * time.Second
	s.ticker.Reset(s.tick)
}

// findOrderToCalc поиск заказов к расчету
func (s *service) findOrderToCalc() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	orderIDS, err := s.orderRepo.FindRegistered(ctx)
	if err != nil {
		log.Println("can't find registered orders")
		log.Println(err)
	}

	// изменяем время тикера в зависимости от результата
	if len(orderIDS) == 0 && s.tick < maxTick {
		s.upTiceker()
	} else if s.tick > defaultTick {
		s.resetTicker()
	}

	// отправляем номера заказа в канал расчета
	for _, orderID := range orderIDS {
		s.calcChan <- orderID
	}
}

// calculate метода расчета заказов
func (s *service) calculate(orderID string) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := s.changeOrderStatus(ctx, orderID, app.Processing)
	if err != nil {
		return
	}
	orderGoods, err := s.productRepo.FindByOrderID(ctx, orderID)
	if err != nil {
		log.Println("service > calculate > can't find goods by orderID")
		log.Println(err)
		s.changeOrderStatus(ctx, orderID, app.Invalid)
		return
	}

	var accrual float64 = 0

	rewards := s.findRewards(ctx, orderGoods.Goods)
	for _, product := range orderGoods.Goods {
		if reward, ok := rewards[product.Description]; ok {
			accrual += calcAccrual(product.Price, reward)
		}
	}
	accrual = math.Round(accrual * 100) / 100

	s.orderRepo.Update(ctx, app.Order{
		ID:      orderID,
		Status:  app.Processed,
		Accrual: &accrual,
	})
}

// changeOrderStatus метод изменения статуса заказа
func (s *service) changeOrderStatus(ctx context.Context, orderID string, status string) error {
	err := s.orderRepo.UpdateStatus(ctx, orderID, app.Processing)
	if err != nil {
		log.Println("service > changeOrderStatus > can't change order status")
		log.Println(err)
	}

	return err
}

// findRewards поиск схем начислений соответствующих товарам в заказе
func (s *service) findRewards(ctx context.Context, goods []app.OrderProduct) map[string]app.Reward {
	res := make(map[string]app.Reward)
	for _, product := range goods {
		reward, err := s.rewardRepo.Find(ctx, product.Description)
		if err != nil && !errors.Is(err, app.ErrRewardNotFound) {
			log.Println("service > findRewards > can't find reward")
			log.Println(err)
		} else {
			res[product.Description] = reward
		}
	}

	return res
}

// calcAccrual метод расчета начисления
func calcAccrual(price float64, acc app.Reward) float64 {
	switch acc.Type {
	case app.Percentage:
		return price * acc.Reward / 100
	case app.Fixed:
		return acc.Reward
	default:
		return 0
	}
}
