package query_test

import (
	"testing"

	"github.com/mohsenjafari-aiio/aiiobackend/internal/shared/query"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// TestProduct represents a test model
type TestProduct struct {
	ID    int64  `gorm:"primaryKey"`
	Name  string `gorm:"not null"`
	Stock int
	Price float64
}

// TestFilter represents test filters
type TestFilter struct {
	Name     string  `filter:"name,CONTAINS"`
	MinStock int     `filter:"stock,>="`
	MaxStock int     `filter:"stock,<="`
	IDs      []int64 `filter:"id,IN"`
	MinPrice float64 `filter:"price,>="`
}

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	err = db.AutoMigrate(&TestProduct{})
	assert.NoError(t, err)

	// Insert test data
	products := []TestProduct{
		{ID: 1, Name: "Product One", Stock: 10, Price: 100.0},
		{ID: 2, Name: "Product Two", Stock: 5, Price: 200.0},
		{ID: 3, Name: "Another Product", Stock: 0, Price: 150.0},
		{ID: 4, Name: "Special Item", Stock: 20, Price: 300.0},
	}

	for _, product := range products {
		db.Create(&product)
	}

	return db
}

func TestQueryBuilder_ApplyFilters(t *testing.T) {
	db := setupTestDB(t)

	tests := []struct {
		name          string
		filter        TestFilter
		expectedCount int
		description   string
	}{
		{
			name: "filter by name contains",
			filter: TestFilter{
				Name: "Product",
			},
			expectedCount: 3,
			description:   "Should find products with 'Product' in name",
		},
		{
			name: "filter by minimum stock",
			filter: TestFilter{
				MinStock: 10,
			},
			expectedCount: 2,
			description:   "Should find products with stock >= 10",
		},
		{
			name: "filter by stock range",
			filter: TestFilter{
				MinStock: 5,
				MaxStock: 15,
			},
			expectedCount: 2,
			description:   "Should find products with stock between 5 and 15",
		},
		{
			name: "filter by IDs",
			filter: TestFilter{
				IDs: []int64{1, 3},
			},
			expectedCount: 2,
			description:   "Should find products with IDs 1 and 3",
		},
		{
			name: "complex filter",
			filter: TestFilter{
				Name:     "Product",
				MinStock: 5,
				MinPrice: 150.0,
			},
			expectedCount: 1,
			description:   "Should find products matching all conditions",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var products []TestProduct

			query := query.NewQueryBuilder(db).
				ApplyFilters(tt.filter).
				Build()

			err := query.Find(&products).Error
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedCount, len(products), tt.description)
		})
	}
}

func TestQueryBuilder_Pagination(t *testing.T) {
	db := setupTestDB(t)

	// Test pagination
	var products []TestProduct

	qb := query.NewQueryBuilder(db).
		SetPagination(1, 2). // First page, 2 items per page
		AddSort("id", query.SortOrderAsc)

	err := qb.Build().Find(&products).Error
	assert.NoError(t, err)
	assert.Equal(t, 2, len(products))
	assert.Equal(t, int64(1), products[0].ID)
	assert.Equal(t, int64(2), products[1].ID)

	// Test second page
	products = []TestProduct{}
	qb = query.NewQueryBuilder(db).
		SetPagination(2, 2). // Second page, 2 items per page
		AddSort("id", query.SortOrderAsc)

	err = qb.Build().Find(&products).Error
	assert.NoError(t, err)
	assert.Equal(t, 2, len(products))
	assert.Equal(t, int64(3), products[0].ID)
	assert.Equal(t, int64(4), products[1].ID)
}

func TestQueryBuilder_Sorting(t *testing.T) {
	db := setupTestDB(t)

	// Test sorting by price descending
	var products []TestProduct

	qb := query.NewQueryBuilder(db).
		AddSort("price", query.SortOrderDesc)

	err := qb.Build().Find(&products).Error
	assert.NoError(t, err)
	assert.Equal(t, 4, len(products))
	assert.Equal(t, "Special Item", products[0].Name) // Highest price
	assert.Equal(t, "Product One", products[3].Name)  // Lowest price
}

func TestQueryBuilder_FluentInterface(t *testing.T) {
	db := setupTestDB(t)

	// Test fluent interface with multiple operations
	var products []TestProduct

	filter := TestFilter{
		MinStock: 1, // Exclude out of stock items
	}

	qb := query.NewQueryBuilder(db).
		ApplyFilters(filter).
		AddSort("price", query.SortOrderAsc).
		SetPagination(1, 2)

	err := qb.Build().Find(&products).Error
	assert.NoError(t, err)
	assert.Equal(t, 2, len(products))
	assert.True(t, products[0].Stock > 0)
	assert.True(t, products[1].Stock > 0)
	// Should be sorted by price ascending
	assert.True(t, products[0].Price <= products[1].Price)
}

func TestCommonFilters(t *testing.T) {
	db := setupTestDB(t)

	cf := query.NewCommonFilters()

	// Test search by name
	nameFilter := cf.SearchByName("Product")

	qb := query.NewQueryBuilder(db).
		AddFilter(nameFilter.ColumnName, nameFilter.Operator, nameFilter.Value)

	var products []TestProduct
	err := qb.Build().Find(&products).Error
	assert.NoError(t, err)
	assert.Equal(t, 3, len(products))

	// Test filter by IDs
	idsFilter := cf.FilterByIDs([]int64{1, 2})

	qb = query.NewQueryBuilder(db).
		AddFilter(idsFilter.ColumnName, idsFilter.Operator, idsFilter.Value)

	products = []TestProduct{}
	err = qb.Build().Find(&products).Error
	assert.NoError(t, err)
	assert.Equal(t, 2, len(products))
}

func TestFindWithPagination(t *testing.T) {
	db := setupTestDB(t)

	var products []TestProduct
	result, err := query.FindWithPagination[TestProduct](db, 1, 2, &products)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 2, len(result.Data))
	assert.Equal(t, int64(4), result.Total)
	assert.Equal(t, 1, result.Page)
	assert.Equal(t, 2, result.PageSize)
	assert.Equal(t, 2, result.TotalPages)
}
