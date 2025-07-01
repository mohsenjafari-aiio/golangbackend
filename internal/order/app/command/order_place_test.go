package command

import (
	"context"
	"errors"
	"testing"

	orderDomain "github.com/mohsenjafari-aiio/aiiobackend/internal/order/domain"
	productDomain "github.com/mohsenjafari-aiio/aiiobackend/internal/product/domain"
	userDomain "github.com/mohsenjafari-aiio/aiiobackend/internal/user/domain"
)

// Mock implementations for testing
type MockUserRepository struct {
	users map[int64]*userDomain.User
	err   error
}

func (m *MockUserRepository) GetByID(ctx context.Context, id int64) (*userDomain.User, error) {
	if m.err != nil {
		return nil, m.err
	}
	if user, exists := m.users[id]; exists {
		return user, nil
	}
	return nil, errors.New("user not found")
}

func (m *MockUserRepository) Save(ctx context.Context, u *userDomain.User) error {
	if m.err != nil {
		return m.err
	}
	if m.users == nil {
		m.users = make(map[int64]*userDomain.User)
	}
	m.users[u.ID] = u
	return nil
}

type MockProductRepository struct {
	products map[int64]*productDomain.Product
	err      error
}

func (m *MockProductRepository) GetByID(ctx context.Context, id int64) (*productDomain.Product, error) {
	if m.err != nil {
		return nil, m.err
	}
	if product, exists := m.products[id]; exists {
		return product, nil
	}
	return nil, errors.New("product not found")
}

func (m *MockProductRepository) Save(ctx context.Context, p *productDomain.Product) error {
	if m.err != nil {
		return m.err
	}
	if m.products == nil {
		m.products = make(map[int64]*productDomain.Product)
	}
	m.products[p.ID] = p
	return nil
}

func (m *MockProductRepository) UpdateStock(ctx context.Context, p *productDomain.Product) error {
	if m.err != nil {
		return m.err
	}
	if m.products == nil {
		m.products = make(map[int64]*productDomain.Product)
	}
	m.products[p.ID] = p
	return nil
}

type MockOrderRepository struct {
	orders map[int64]*orderDomain.Order
	err    error
}

func (m *MockOrderRepository) Save(ctx context.Context, o *orderDomain.Order) error {
	if m.err != nil {
		return m.err
	}
	if m.orders == nil {
		m.orders = make(map[int64]*orderDomain.Order)
	}
	// Simulate auto-incrementing ID
	o.ID = int64(len(m.orders) + 1)
	m.orders[o.ID] = o
	return nil
}

func (m *MockOrderRepository) GetByID(ctx context.Context, id int64) (*orderDomain.Order, error) {
	if m.err != nil {
		return nil, m.err
	}
	if order, exists := m.orders[id]; exists {
		return order, nil
	}
	return nil, errors.New("order not found")
}

func TestPlaceOrderHandler_Handle_Success(t *testing.T) {
	// Arrange
	userRepo := &MockUserRepository{
		users: map[int64]*userDomain.User{
			1: {ID: 1, Email: "test@example.com", Active: true},
		},
	}
	
	productRepo := &MockProductRepository{
		products: map[int64]*productDomain.Product{
			1: {ID: 1, Name: "Test Product", Stock: 10},
		},
	}
	
	orderRepo := &MockOrderRepository{}
	
	handler := &PlaceOrderHandler{
		UserRepo:    userRepo,
		ProductRepo: productRepo,
		OrderRepo:   orderRepo,
	}
	
	cmd := PlaceOrderCommand{
		UserID:    1,
		ProductID: 1,
		Quantity:  2,
	}
	
	// Act
	err := handler.Handle(context.Background(), cmd)
	
	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	
	// Verify product stock was updated
	product, _ := productRepo.GetByID(context.Background(), 1)
	if product.Stock != 8 {
		t.Errorf("Expected product stock to be 8, got %d", product.Stock)
	}
	
	// Verify order was saved
	if len(orderRepo.orders) != 1 {
		t.Errorf("Expected 1 order to be saved, got %d", len(orderRepo.orders))
	}
	
	// Verify order details
	for _, order := range orderRepo.orders {
		if order.UserID != 1 {
			t.Errorf("Expected order UserID to be 1, got %d", order.UserID)
		}
		if order.ProductID != 1 {
			t.Errorf("Expected order ProductID to be 1, got %d", order.ProductID)
		}
		if order.Quantity != 2 {
			t.Errorf("Expected order Quantity to be 2, got %d", order.Quantity)
		}
		if order.Status != "CONFIRMED" {
			t.Errorf("Expected order Status to be 'CONFIRMED', got %s", order.Status)
		}
	}
}

func TestPlaceOrderHandler_Handle_UserNotFound(t *testing.T) {
	// Arrange
	userRepo := &MockUserRepository{users: map[int64]*userDomain.User{}}
	productRepo := &MockProductRepository{}
	orderRepo := &MockOrderRepository{}
	
	handler := &PlaceOrderHandler{
		UserRepo:    userRepo,
		ProductRepo: productRepo,
		OrderRepo:   orderRepo,
	}
	
	cmd := PlaceOrderCommand{
		UserID:    999,
		ProductID: 1,
		Quantity:  2,
	}
	
	// Act
	err := handler.Handle(context.Background(), cmd)
	
	// Assert
	if err == nil {
		t.Error("Expected error for user not found, got nil")
	}
	if err.Error() != "user not found" {
		t.Errorf("Expected error message 'user not found', got %s", err.Error())
	}
}

func TestPlaceOrderHandler_Handle_ProductNotFound(t *testing.T) {
	// Arrange
	userRepo := &MockUserRepository{
		users: map[int64]*userDomain.User{
			1: {ID: 1, Email: "test@example.com", Active: true},
		},
	}
	productRepo := &MockProductRepository{products: map[int64]*productDomain.Product{}}
	orderRepo := &MockOrderRepository{}
	
	handler := &PlaceOrderHandler{
		UserRepo:    userRepo,
		ProductRepo: productRepo,
		OrderRepo:   orderRepo,
	}
	
	cmd := PlaceOrderCommand{
		UserID:    1,
		ProductID: 999,
		Quantity:  2,
	}
	
	// Act
	err := handler.Handle(context.Background(), cmd)
	
	// Assert
	if err == nil {
		t.Error("Expected error for product not found, got nil")
	}
	if err.Error() != "product not found" {
		t.Errorf("Expected error message 'product not found', got %s", err.Error())
	}
}

func TestPlaceOrderHandler_Handle_InsufficientStock(t *testing.T) {
	// Arrange
	userRepo := &MockUserRepository{
		users: map[int64]*userDomain.User{
			1: {ID: 1, Email: "test@example.com", Active: true},
		},
	}
	
	productRepo := &MockProductRepository{
		products: map[int64]*productDomain.Product{
			1: {ID: 1, Name: "Test Product", Stock: 1}, // Only 1 in stock
		},
	}
	
	orderRepo := &MockOrderRepository{}
	
	handler := &PlaceOrderHandler{
		UserRepo:    userRepo,
		ProductRepo: productRepo,
		OrderRepo:   orderRepo,
	}
	
	cmd := PlaceOrderCommand{
		UserID:    1,
		ProductID: 1,
		Quantity:  5, // Requesting more than available
	}
	
	// Act
	err := handler.Handle(context.Background(), cmd)
	
	// Assert
	if err == nil {
		t.Error("Expected error for insufficient stock, got nil")
	}
	if err.Error() != "insufficient stock" {
		t.Errorf("Expected error message 'insufficient stock', got %s", err.Error())
	}
	
	// Verify no order was saved
	if len(orderRepo.orders) != 0 {
		t.Errorf("Expected no orders to be saved, got %d", len(orderRepo.orders))
	}
}

func TestPlaceOrderHandler_Handle_RepositoryError(t *testing.T) {
	// Arrange
	userRepo := &MockUserRepository{
		users: map[int64]*userDomain.User{
			1: {ID: 1, Email: "test@example.com", Active: true},
		},
	}
	
	productRepo := &MockProductRepository{
		products: map[int64]*productDomain.Product{
			1: {ID: 1, Name: "Test Product", Stock: 10},
		},
	}
	
	orderRepo := &MockOrderRepository{
		err: errors.New("database connection failed"), // Simulate DB error
	}
	
	handler := &PlaceOrderHandler{
		UserRepo:    userRepo,
		ProductRepo: productRepo,
		OrderRepo:   orderRepo,
	}
	
	cmd := PlaceOrderCommand{
		UserID:    1,
		ProductID: 1,
		Quantity:  2,
	}
	
	// Act
	err := handler.Handle(context.Background(), cmd)
	
	// Assert
	if err == nil {
		t.Error("Expected error from repository, got nil")
	}
	if err.Error() != "database connection failed" {
		t.Errorf("Expected error message 'database connection failed', got %s", err.Error())
	}
}
