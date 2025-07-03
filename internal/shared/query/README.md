# Dynamic Query System for GORM

This package provides a powerful and efficient dynamic filtering system for GORM that allows you to build complex queries using struct-based filters, fluent interfaces, and common patterns.

## Features

- **Struct-based Filtering**: Define filters using Go structs with tags
- **Fluent Interface**: Chain operations for complex queries
- **Pagination Support**: Built-in pagination with metadata
- **Type Safety**: Leverage Go's type system for compile-time safety
- **Extensible**: Easy to add new operators and patterns
- **Clean Architecture**: Follows repository pattern and dependency inversion

## Quick Start

### 1. Define Filter Structs

```go
type ProductFilter struct {
    Name         string   `filter:"name,ILIKE"`           // LIKE search
    MinStock     int      `filter:"stock,>="`             // Greater than or equal
    MaxStock     int      `filter:"stock,<="`             // Less than or equal
    IDs          []int64  `filter:"id,IN"`                // IN clause
    Active       *bool    `filter:"active"`               // Exact match (use pointer for optional)
}
```

### 2. Use in Repository

```go
func (r *GormProductRepository) FindWithFilters(ctx context.Context, filters interface{}) ([]*domain.Product, error) {
    var products []*domain.Product
    
    query := query.NewQueryBuilder(r.db).
        ApplyFilters(filters).
        Build().
        WithContext(ctx)
    
    err := query.Find(&products).Error
    return products, err
}
```

### 3. Use in Service Layer

```go
func (s *ProductService) SearchProducts(ctx context.Context, searchTerm string) ([]*domain.Product, error) {
    filter := ProductFilter{
        Name:         searchTerm,
        StockGreater: 0, // Only products in stock
    }
    
    return s.productRepo.FindWithFilters(ctx, filter)
}
```

## Available Operators

| Operator | Tag | Description | Example |
|----------|-----|-------------|---------|
| `=` | `filter:"column"` or `filter:"column,="` | Exact match | `Name string \`filter:"name"\`` |
| `!=` | `filter:"column,!="` | Not equal | `Status string \`filter:"status,!="\`` |
| `>` | `filter:"column,>"` | Greater than | `Price float64 \`filter:"price,>"\`` |
| `<` | `filter:"column,<"` | Less than | `Price float64 \`filter:"price,<"\`` |
| `>=` | `filter:"column,>="` | Greater or equal | `Stock int \`filter:"stock,>="\`` |
| `<=` | `filter:"column,<="` | Less or equal | `Stock int \`filter:"stock,<="\`` |
| `ILIKE` | `filter:"column,ILIKE"` | Case-insensitive contains | `Name string \`filter:"name,ILIKE"\`` |
| `IN` | `filter:"column,IN"` | In list | `IDs []int64 \`filter:"id,IN"\`` |
| `NOT IN` | `filter:"column,NOT IN"` | Not in list | `Status []string \`filter:"status,NOT IN"\`` |
| `IS NULL` | `filter:"column,IS NULL"` | Is null | `DeletedAt *time.Time \`filter:"deleted_at,IS NULL"\`` |
| `IS NOT NULL` | `filter:"column,IS NOT NULL"` | Is not null | `UpdatedAt *time.Time \`filter:"updated_at,IS NOT NULL"\`` |
| `STARTS_WITH` | `filter:"column,STARTS_WITH"` | Starts with | `Email string \`filter:"email,STARTS_WITH"\`` |
| `ENDS_WITH` | `filter:"column,ENDS_WITH"` | Ends with | `Domain string \`filter:"domain,ENDS_WITH"\`` |

## Advanced Usage

### Fluent Interface

```go
query := query.NewQueryBuilder(db).
    ApplyFilters(filter).
    AddSort("created_at", query.SortOrderDesc).
    SetPagination(1, 20).
    AddPreload("User").
    AddPreload("Product", "stock > ?", 0).
    Build()
```

### Pagination with Metadata

```go
result, err := s.productRepo.FindWithPagination(ctx, filter, page, pageSize)
if err != nil {
    return nil, err
}

// result.Data contains the products
// result.Total contains total count
// result.Page, result.PageSize, result.TotalPages contain pagination info
```

### Custom Preloads

```go
qb := query.NewQueryBuilder(db).
    ApplyFilters(filter).
    AddCustomPreload("User", func(db *gorm.DB) *gorm.DB {
        return db.Preload("User", "active = ?", true)
    })
```

### Manual Filter Building

```go
qb := query.NewQueryBuilder(db).
    AddFilter("name", query.OperatorContains, "product").
    AddFilter("stock", query.OperatorGreaterThan, 0).
    AddFilter("category_id", query.OperatorIn, []int64{1, 2, 3})
```

## Common Patterns

### Date Range Filtering

```go
type OrderFilter struct {
    CreatedAfter  *time.Time `filter:"created_at,>="`
    CreatedBefore *time.Time `filter:"created_at,<="`
}
```

### Optional Boolean Fields

```go
type UserFilter struct {
    Active *bool `filter:"active"` // Use pointer to distinguish false from zero value
}
```

### Multiple Sort Fields

```go
qb := query.NewQueryBuilder(db).
    AddSort("priority", query.SortOrderDesc).
    AddSort("created_at", query.SortOrderAsc)
```

### Search with Aggregation

```go
qb := query.NewQueryBuilder(db).
    ApplyFilters(filter).
    AddGroupBy("category_id").
    AddHaving("COUNT(*) > ?", 5)
```

## Repository Integration

### Base Repository Pattern

```go
type GormProductRepository struct {
    db *gorm.DB
    *query.BaseRepository
}

func NewGormProductRepository(db *gorm.DB) port.ProductRepository {
    return &GormProductRepository{
        db: db,
        BaseRepository: query.NewBaseRepository(db),
    }
}
```

### Interface Extension

```go
type ProductRepository interface {
    // Basic CRUD
    GetByID(ctx context.Context, id int64) (*domain.Product, error)
    Save(ctx context.Context, p *domain.Product) error
    
    // Advanced queries
    FindWithFilters(ctx context.Context, filters interface{}) ([]*domain.Product, error)
    FindWithPagination(ctx context.Context, filters interface{}, page, pageSize int) (*query.PaginatedResult[*domain.Product], error)
    CountWithFilters(ctx context.Context, filters interface{}) (int64, error)
}
```

## Performance Considerations

1. **Index your columns**: Make sure frequently filtered columns are indexed
2. **Limit results**: Always use pagination for large datasets
3. **Selective preloading**: Only preload relationships you need
4. **Use appropriate operators**: ILIKE is slower than exact matches

## Testing

```go
func TestProductFiltering(t *testing.T) {
    filter := ProductFilter{
        Name:     "test",
        MinStock: 10,
    }
    
    products, err := repo.FindWithFilters(ctx, filter)
    assert.NoError(t, err)
    assert.True(t, len(products) > 0)
}
```

## Best Practices

1. **Use struct tags**: Always define filter tags for clarity
2. **Handle zero values**: Use pointers for optional boolean fields
3. **Validate input**: Check pagination parameters in your service layer
4. **Error handling**: Always handle GORM errors appropriately
5. **Context usage**: Pass context for cancellation and timeouts

## Migration from Existing Code

### Before
```go
func (r *repo) FindProducts(name string, minStock int) ([]*Product, error) {
    query := r.db
    if name != "" {
        query = query.Where("name ILIKE ?", "%"+name+"%")
    }
    if minStock > 0 {
        query = query.Where("stock >= ?", minStock)
    }
    
    var products []*Product
    err := query.Find(&products).Error
    return products, err
}
```

### After
```go
func (r *repo) FindProducts(filters ProductFilter) ([]*Product, error) {
    var products []*Product
    
    query := query.NewQueryBuilder(r.db).
        ApplyFilters(filters).
        Build()
    
    err := query.Find(&products).Error
    return products, err
}
```

This system provides a clean, type-safe, and efficient way to handle dynamic queries in your Go GORM applications while maintaining the principles of clean architecture.
