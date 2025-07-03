package query

import "time"

// ProductFilter represents filters for product queries
type ProductFilter struct {
	ID           int64   `filter:"id"`
	Name         string  `filter:"name,CONTAINS"`
	MinStock     int     `filter:"stock,>="`
	MaxStock     int     `filter:"stock,<="`
	StockGreater int     `filter:"stock,>"`
	IDs          []int64 `filter:"id,IN"`
}

// UserFilter represents filters for user queries
type UserFilter struct {
	ID     int64   `filter:"id"`
	Email  string  `filter:"email,CONTAINS"`
	Active *bool   `filter:"active"` // Use pointer to distinguish between false and zero value
	IDs    []int64 `filter:"id,IN"`
}

// OrderFilter represents filters for order queries
type OrderFilter struct {
	ID            int64      `filter:"id"`
	UserID        int64      `filter:"user_id"`
	ProductID     int64      `filter:"product_id"`
	Status        string     `filter:"status"`
	MinQuantity   int        `filter:"quantity,>="`
	MaxQuantity   int        `filter:"quantity,<="`
	UserIDs       []int64    `filter:"user_id,IN"`
	ProductIDs    []int64    `filter:"product_id,IN"`
	Statuses      []string   `filter:"status,IN"`
	CreatedAfter  *time.Time `filter:"created_at,>="`
	CreatedBefore *time.Time `filter:"created_at,<="`
}

// Example usage patterns for different scenarios
type ExampleFilters struct{}

// ProductSearchFilter for complex product searches
type ProductSearchFilter struct {
	// Basic filters
	SearchTerm string  `filter:"name,CONTAINS"`
	MinPrice   float64 `filter:"price,>="`
	MaxPrice   float64 `filter:"price,<="`
	InStock    *bool   `filter:"stock,>"` // Will be converted to stock > 0

	// Advanced filters
	CategoryIDs []int64  `filter:"category_id,IN"`
	BrandIDs    []int64  `filter:"brand_id,IN"`
	Tags        []string `filter:"tags,IN"`

	// Date filters
	CreatedAfter  *time.Time `filter:"created_at,>="`
	CreatedBefore *time.Time `filter:"created_at,<="`
	UpdatedAfter  *time.Time `filter:"updated_at,>="`
}

// UserSearchFilter for complex user searches
type UserSearchFilter struct {
	SearchTerm     string     `filter:"email,CONTAINS"`
	Active         *bool      `filter:"active"`
	RoleIDs        []int64    `filter:"role_id,IN"`
	DepartmentIDs  []int64    `filter:"department_id,IN"`
	CreatedAfter   *time.Time `filter:"created_at,>="`
	LastLoginAfter *time.Time `filter:"last_login_at,>="`
}

// OrderReportFilter for order reporting and analytics
type OrderReportFilter struct {
	UserIDs       []int64    `filter:"user_id,IN"`
	ProductIDs    []int64    `filter:"product_id,IN"`
	Statuses      []string   `filter:"status,IN"`
	MinAmount     float64    `filter:"total_amount,>="`
	MaxAmount     float64    `filter:"total_amount,<="`
	DateFrom      *time.Time `filter:"created_at,>="`
	DateTo        *time.Time `filter:"created_at,<="`
	PaymentMethod string     `filter:"payment_method"`
}
