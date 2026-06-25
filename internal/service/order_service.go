package service

import (
	"context"
	"myservice/internal/models"
)

type OrderRepo interface {
	Create(ctx context.Context, o *models.Order) error
}

type OrderService struct {
	repo OrderRepo
}

func NewOrderService(repo OrderRepo) *OrderService {
	return &OrderService{repo: repo}
}

func (s *OrderService) CreateOrder(ctx context.Context, userID int, req models.CreateOrderRequest) (*models.Order, error) {
	order := &models.Order{
		UserID:   userID,
		ItemName: req.ItemName,
		Quantity: req.Quantity,
		Status:   "pending", // Статус по умолчанию для нового заказа
	}

	if err := s.repo.Create(ctx, order); err != nil {
		return nil, err
	}
	return order, nil
}