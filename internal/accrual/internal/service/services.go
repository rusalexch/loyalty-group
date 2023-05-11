package service

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/rusalexch/loyalty-group/internal/accrual/internal/common"
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
	defer s.destroy()

	go s.run()

	return s
}

// Ping проверка работоспособности
func (s *service) Ping(ctx context.Context) error {
	return s.store.Ping(ctx)
}

// GetOrder получение данных по заказу
func (s *service) GetOrder(ctx context.Context, orderID string) (common.Order, error) {
	return s.orderRepo.FindByID(ctx, orderID)
}

// AddReward добавление/изменение схемы начисления
func (s *service) AddReward(ctx context.Context, reward common.Reward) error {
	return s.rewardRepo.Add(ctx, reward)
}

// AddOrder добавление нового заказа
func (s *service) AddOrder(ctx context.Context, order common.OrderGoods) error {
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
func (s *service) destroy() {
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
		for _ = range s.ticker.C {
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
	for _, orderId := range orderIDS {
		s.calcChan <- orderId
	}
}

// calculate метода расчета заказов
func (s *service) calculate(orderID string) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := s.changeOrderStatus(ctx, orderID, common.Processing)
	if err != nil {
		return
	}
	orderGoods, err := s.productRepo.FindByOrderID(ctx, orderID)
	if err != nil {
		log.Println("service > calculate > can't find goods by orderID")
		log.Println(err)
		s.changeOrderStatus(ctx, orderID, common.Invalid)
		return
	}

	var accrual float64 = 0

	rewards := s.findRewards(ctx, orderGoods.Goods)
	for _, product := range orderGoods.Goods {
		if reward, ok := rewards[product.Description]; ok {
			accrual += calcAccrual(product.Price, reward)
		}
	}

	s.orderRepo.Update(ctx, common.Order{
		ID:      orderID,
		Status:  common.Processed,
		Accrual: &accrual,
	})
}

// changeOrderStatus метод изменения статуса заказа
func (s *service) changeOrderStatus(ctx context.Context, orderID string, status common.OrderStatus) error {
	err := s.orderRepo.UpdateStatus(ctx, orderID, common.Processing)
	if err != nil {
		log.Println("service > changeOrderStatus > can't change order status")
		log.Println(err)
	}

	return err
}

// findRewards поиск схем начислений соответствующих товарам в заказе
func (s *service) findRewards(ctx context.Context, goods []common.OrderProduct) map[string]common.Reward {
	res := make(map[string]common.Reward)
	for _, product := range goods {
		reward, err := s.rewardRepo.Find(ctx, product.Description)
		if err != nil && !errors.Is(err, common.ErrOrderNotFound) {
			log.Println("service > findRewards > can't find reward")
			log.Println(err)
		} else {
			res[product.Description] = reward
		}
	}

	return res
}

// calcAccrual метод расчета начисления
func calcAccrual(price float64, acc common.Reward) float64 {
	switch acc.Type {
	case common.Percentage:
		return price * acc.Reward / 100
	case common.Fixed:
		return acc.Reward
	default:
		return 0
	}
}
