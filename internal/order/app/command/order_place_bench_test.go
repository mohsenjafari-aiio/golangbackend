package command

import (
	"context"
	"testing"

	productDomain "github.com/mohsenjafari-aiio/aiiobackend/internal/product/domain"
	userDomain "github.com/mohsenjafari-aiio/aiiobackend/internal/user/domain"
)

func BenchmarkPlaceOrderHandler_Handle(b *testing.B) {
	cmd := PlaceOrderCommand{
		UserID:    1,
		ProductID: 1,
		Quantity:  1,
	}

	b.ResetTimer()

	// Benchmark
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		// Setup fresh state for each iteration (outside measurement)
		userRepo := &MockUserRepository{
			users: map[int64]*userDomain.User{
				1: {ID: 1, Email: "test@example.com", Active: true},
			},
		}

		productRepo := &MockProductRepository{
			products: map[int64]*productDomain.Product{
				1: {ID: 1, Name: "Test Product", Stock: 1000000},
			},
		}

		orderRepo := &MockOrderRepository{}

		handler := &PlaceOrderHandler{
			UserRepo:    userRepo,
			ProductRepo: productRepo,
			OrderRepo:   orderRepo,
		}
		b.StartTimer()

		// This is what we're actually measuring
		handler.Handle(context.Background(), cmd)
	}
}

// Benchmark with reusable mocks (measures handler + mock overhead)
func BenchmarkPlaceOrderHandler_Handle_WithReusableMocks(b *testing.B) {
	userRepo := &MockUserRepository{
		users: map[int64]*userDomain.User{
			1: {ID: 1, Email: "test@example.com", Active: true},
		},
	}

	productRepo := &MockProductRepository{
		products: map[int64]*productDomain.Product{
			1: {ID: 1, Name: "Test Product", Stock: 1000000},
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
		Quantity:  1,
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		// Reset only the stock without recreating objects
		productRepo.products[1].Stock = 1000000
		handler.Handle(context.Background(), cmd)
	}
}

// Benchmark error path performance
func BenchmarkPlaceOrderHandler_Handle_UserNotFound(b *testing.B) {
	cmd := PlaceOrderCommand{
		UserID:    999, // Non-existent user
		ProductID: 1,
		Quantity:  1,
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		userRepo := &MockUserRepository{users: map[int64]*userDomain.User{}}
		productRepo := &MockProductRepository{}
		orderRepo := &MockOrderRepository{}

		handler := &PlaceOrderHandler{
			UserRepo:    userRepo,
			ProductRepo: productRepo,
			OrderRepo:   orderRepo,
		}
		b.StartTimer()

		handler.Handle(context.Background(), cmd)
	}
}

// Benchmark memory allocations more precisely
func BenchmarkPlaceOrderHandler_Handle_MemoryOptimized(b *testing.B) {
	// Pre-allocate everything to minimize GC impact
	users := map[int64]*userDomain.User{
		1: {ID: 1, Email: "test@example.com", Active: true},
	}

	cmd := PlaceOrderCommand{
		UserID:    1,
		ProductID: 1,
		Quantity:  1,
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		userRepo := &MockUserRepository{users: users}
		productRepo := &MockProductRepository{
			products: map[int64]*productDomain.Product{
				1: {ID: 1, Name: "Test Product", Stock: 1000000},
			},
		}
		orderRepo := &MockOrderRepository{}

		handler := &PlaceOrderHandler{
			UserRepo:    userRepo,
			ProductRepo: productRepo,
			OrderRepo:   orderRepo,
		}

		handler.Handle(context.Background(), cmd)
	}
}
