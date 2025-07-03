# Migration Guide: From Manual Queries to Dynamic Filter System

This guide shows how to migrate your existing GORM code to use the new dynamic filtering system.

## Before and After Examples

### 1. Simple Product Search

#### Before (Manual Query Building)
```go
func (r *GormProductRepository) FindProducts(ctx context.Context, name string, minStock int) ([]*domain.Product, error) {
    var products []*domain.Product
    
    query := r.db.WithContext(ctx)
    
    if name != "" {
        query = query.Where("name LIKE ?", "%"+name+"%")
    }
    
    if minStock > 0 {
        query = query.Where("stock >= ?", minStock)
    }
    
    err := query.Find(&products).Error
    return products, err
}
```

#### After (Using Dynamic Filters)
```go
func (r *GormProductRepository) FindProducts(ctx context.Context, filter query.ProductFilter) ([]*domain.Product, error) {
    var products []*domain.Product
    
    query := r.QueryBuilder().
        ApplyFilters(filter).
        Build().
        WithContext(ctx)
    
    err := query.Find(&products).Error
    return products, err
}
```

### 2. Complex Search with Pagination

#### Before (Manual Query Building)
```go
func (r *GormProductRepository) SearchWithPagination(ctx context.Context, name string, minStock, maxStock, page, pageSize int, sortBy string, sortOrder string) ([]*domain.Product, int64, error) {
    var products []*domain.Product
    var total int64
    
    baseQuery := r.db.WithContext(ctx)
    
    // Apply filters
    if name != "" {
        baseQuery = baseQuery.Where("name LIKE ?", "%"+name+"%")
    }
    if minStock > 0 {
        baseQuery = baseQuery.Where("stock >= ?", minStock)
    }
    if maxStock > 0 {
        baseQuery = baseQuery.Where("stock <= ?", maxStock)
    }
    
    // Count total
    if err := baseQuery.Model(&domain.Product{}).Count(&total).Error; err != nil {
        return nil, 0, err
    }
    
    // Apply sorting
    if sortBy != "" {
        orderClause := sortBy
        if sortOrder == "desc" {
            orderClause += " DESC"
        }
        baseQuery = baseQuery.Order(orderClause)
    }
    
    // Apply pagination
    offset := (page - 1) * pageSize
    err := baseQuery.Offset(offset).Limit(pageSize).Find(&products).Error
    
    return products, total, err
}
```

#### After (Using Dynamic Filters)
```go
func (r *GormProductRepository) SearchWithPagination(ctx context.Context, filter query.ProductFilter, page, pageSize int) (*query.PaginatedResult[*domain.Product], error) {
    return r.FindWithPagination(ctx, filter, page, pageSize)
}
```

### 3. Service Layer Changes

#### Before (Service Layer)
```go
func (s *ProductService) SearchProducts(ctx context.Context, req ProductSearchRequest) (*ProductSearchResponse, error) {
    products, total, err := s.productRepo.SearchWithPagination(
        ctx,
        req.Name,
        req.MinStock,
        req.MaxStock,
        req.Page,
        req.PageSize,
        "name",
        "asc",
    )
    if err != nil {
        return nil, err
    }
    
    totalPages := int((total + int64(req.PageSize) - 1) / int64(req.PageSize))
    
    return &ProductSearchResponse{
        Products: products,
        Total: total,
        Page: req.Page,
        PageSize: req.PageSize,
        TotalPages: totalPages,
    }, nil
}
```

#### After (Service Layer)
```go
func (s *ProductService) SearchProducts(ctx context.Context, req ProductSearchRequest) (*query.PaginatedResult[*domain.Product], error) {
    filter := query.ProductFilter{
        Name:     req.Name,
        MinStock: req.MinStock,
        MaxStock: req.MaxStock,
    }
    
    return s.productRepo.FindWithPagination(ctx, filter, req.Page, req.PageSize)
}
```

## Migration Steps

### Step 1: Install the Query System
1. Copy the query package to `internal/shared/query/`
2. Run `go mod tidy` to ensure dependencies are satisfied

### Step 2: Update Repository Interface
```go
// Add these methods to your repository interface
type ProductRepository interface {
    // Existing methods...
    GetByID(ctx context.Context, id int64) (*domain.Product, error)
    Save(ctx context.Context, p *domain.Product) error
    
    // New methods using the query system
    FindWithFilters(ctx context.Context, filters interface{}) ([]*domain.Product, error)
    FindWithPagination(ctx context.Context, filters interface{}, page, pageSize int) (*query.PaginatedResult[*domain.Product], error)
    CountWithFilters(ctx context.Context, filters interface{}) (int64, error)
}
```

### Step 3: Update Repository Implementation
```go
type GormProductRepository struct {
    db *gorm.DB
    *query.BaseRepository  // Embed the base repository
}

func NewGormProductRepository(db *gorm.DB) port.ProductRepository {
    return &GormProductRepository{
        db:             db,
        BaseRepository: query.NewBaseRepository(db),
    }
}

// Implement the new methods...
```

### Step 4: Define Filter Structs
```go
type ProductFilter struct {
    Name         string   `filter:"name,CONTAINS"`
    MinStock     int      `filter:"stock,>="`
    MaxStock     int      `filter:"stock,<="`
    IDs          []int64  `filter:"id,IN"`
    Active       *bool    `filter:"active"`
}
```

### Step 5: Update Service Layer
Replace manual parameter passing with filter structs:

```go
// Before
func (s *Service) SearchProducts(name string, minStock int, page int) (results, error)

// After  
func (s *Service) SearchProducts(filter ProductFilter, page int) (results, error)
```

### Step 6: Update API Handlers
Convert request parameters to filter structs:

```go
func (h *Handler) SearchProducts(w http.ResponseWriter, r *http.Request) {
    filter := ProductFilter{
        Name:     r.URL.Query().Get("name"),
        MinStock: parseIntParam(r.URL.Query().Get("min_stock")),
        MaxStock: parseIntParam(r.URL.Query().Get("max_stock")),
    }
    
    result, err := h.service.SearchProducts(ctx, filter, page, pageSize)
    // ...
}
```

## Benefits After Migration

### 1. Reduced Code Complexity
- **Before**: 50+ lines for complex search function
- **After**: 5-10 lines using the query builder

### 2. Type Safety
- **Before**: Easy to make mistakes with string concatenation
- **After**: Compile-time type checking with struct tags

### 3. Consistency
- **Before**: Each search function has different parameter patterns
- **After**: Consistent interface across all search functions

### 4. Maintainability
- **Before**: Adding new filters requires modifying multiple functions
- **After**: Just add fields to the filter struct

### 5. Testing
- **Before**: Complex mocking and parameter setup
- **After**: Simple struct initialization for test cases

## Common Pitfalls and Solutions

### 1. Zero Values
**Problem**: Go zero values (0, "", false) are filtered out by default

**Solution**: Use pointers for optional fields
```go
type Filter struct {
    Active *bool `filter:"active"`  // Use pointer to distinguish false from zero
}
```

### 2. Database Compatibility
**Problem**: ILIKE is PostgreSQL-specific

**Solution**: Use CONTAINS operator instead
```go
Name string `filter:"name,CONTAINS"`  // Works with all databases
```

### 3. Complex Conditions
**Problem**: Need OR conditions or complex nested logic

**Solution**: Use manual query building for complex cases, or extend the system
```go
// For complex cases, fall back to manual building
query := r.QueryBuilder().
    ApplyFilters(simpleFilters).
    Build().
    Where("(condition1 OR condition2)")
```

### 4. Performance
**Problem**: Worried about performance impact

**Solution**: The system generates the same SQL as manual queries
- Use database indexes on filtered columns
- Use pagination for large datasets
- Profile your queries to identify bottlenecks

## Testing Migration

### Test Both Approaches Side by Side
```go
func TestMigration(t *testing.T) {
    // Test old approach
    productsOld, err := repo.FindProductsOld(ctx, "test", 10)
    require.NoError(t, err)
    
    // Test new approach
    filter := query.ProductFilter{Name: "test", MinStock: 10}
    productsNew, err := repo.FindWithFilters(ctx, filter)
    require.NoError(t, err)
    
    // Compare results
    assert.Equal(t, len(productsOld), len(productsNew))
}
```

### Benchmark Performance
```go
func BenchmarkOldApproach(b *testing.B) {
    for i := 0; i < b.N; i++ {
        repo.FindProductsOld(ctx, "test", 10)
    }
}

func BenchmarkNewApproach(b *testing.B) {
    filter := query.ProductFilter{Name: "test", MinStock: 10}
    for i := 0; i < b.N; i++ {
        repo.FindWithFilters(ctx, filter)
    }
}
```

This migration guide provides a clear path to adopt the new dynamic filtering system while maintaining backwards compatibility during the transition.
