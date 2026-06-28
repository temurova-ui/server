package service

import (
	"context"
	"errors"
	"myservice/internal/models"
	"myservice/internal/repository"
)

type OrderService interface {
	CreateOrder(ctx context.Context, userID int, req models.CreateOrderRequest) (*models.Order, error)
	CancelOrder(ctx context.Context, orderID int, req models.CancelOrder) error
	GetOrders(ctx context.Context, userID int) ([]models.Order, error)
}

type OrderServices struct {
	repo repository.OrderRepo
}

func NewOrderService(repo repository.OrderRepo) OrderService {
	return &OrderServices{repo: repo}
}

func (s *OrderServices) CreateOrder(ctx context.Context, userID int, req models.CreateOrderRequest) (*models.Order, error) {
	order := &models.Order{
		UserID:   userID,
		ItemName: req.ItemName,
		Quantity: req.Quantity,
		Status:   "new",
	}

	if err := s.repo.Create(ctx, order); err != nil {
		return nil, err
	}
	return order, nil
}

func (s *OrderServices) CancelOrder(ctx context.Context, orderID int, req models.CancelOrder) error {
	if orderID <= 0 {
		return errors.New("orderID must be positive")
	}

	order, err := s.repo.GetByID(ctx, orderID)
	if err != nil {
		return err
	}
	if order.UserID != req.UserID {
		return errors.New("order with this user not found")
	}

	if order.Status != "new" {
		return errors.New("cannot cancel order: status is not new")
	}

	statusToSet := req.Status
	if statusToSet == "" {
		statusToSet = "cancelled"
	}

	err = s.repo.CancelOrder(ctx, orderID, statusToSet)
	if err != nil {
		return err
	}

	return nil
}

func (s *OrderServices) GetOrders(ctx context.Context, userID int) ([]models.Order, error) {
	return s.repo.GetMyOrders(ctx, userID)
}