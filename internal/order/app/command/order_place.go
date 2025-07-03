package command

import (
	"context"
	"errors"

	orderDomain "github.com/mohsenjafari-aiio/aiiobackend/internal/order/domain"
	productDomain "github.com/mohsenjafari-aiio/aiiobackend/internal/product/domain"
	userDomain "github.com/mohsenjafari-aiio/aiiobackend/internal/user/domain"
)

type PlaceOrderCommand struct {
	UserID    int64
	ProductID int64
	Quantity  int
}

type PlaceOrderHandler struct {
	OrderRepo   orderDomain.OrderRepository
	UserRepo    userDomain.UserRepository
	ProductRepo productDomain.ProductRepository
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
