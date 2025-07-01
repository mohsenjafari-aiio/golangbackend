package command

import (
	"context"
	"errors"

	orderDomain "github.com/mohsenjafari-aiio/aiiobackend/internal/order/domain"
	orderPort "github.com/mohsenjafari-aiio/aiiobackend/internal/order/port"
	productPort "github.com/mohsenjafari-aiio/aiiobackend/internal/product/port"
	userPort "github.com/mohsenjafari-aiio/aiiobackend/internal/user/port"
)

type PlaceOrderCommand struct {
	UserID    int64
	ProductID int64
	Quantity  int
}

type PlaceOrderHandler struct {
	OrderRepo   orderPort.OrderRepository
	UserRepo    userPort.UserRepository
	ProductRepo productPort.ProductRepository
}

func (h *PlaceOrderHandler) Handle(ctx context.Context, cmd PlaceOrderCommand) error {
	// Get user by ID using repository
	u, err := h.UserRepo.GetByID(ctx, cmd.UserID)
	if err != nil {
		return errors.New("user not found")
	}

	// Get product by ID using repository
	p, err := h.ProductRepo.GetByID(ctx, cmd.ProductID)
	if err != nil {
		return errors.New("product not found")
	}

	// Reserve product stock (domain business logic)
	if err := p.Reserve(cmd.Quantity); err != nil {
		return err
	}

	// Create and confirm order
	o := orderDomain.NewOrder(u.ID, p.ID, cmd.Quantity)
	o.Confirm()

	// Update product stock using repository
	if err := h.ProductRepo.UpdateStock(ctx, p); err != nil {
		return err
	}

	// Save order using repository
	return h.OrderRepo.Save(ctx, o)
}
