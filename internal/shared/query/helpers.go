package query

import (
	"gorm.io/gorm"
)

// PaginatedResult represents a paginated query result
type PaginatedResult[T any] struct {
	Data       []T   `json:"data"`
	Total      int64 `json:"total"`
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	TotalPages int   `json:"total_pages"`
}

// FindWithPagination performs a paginated query and returns structured result
func FindWithPagination[T any](db *gorm.DB, page, pageSize int, result *[]T) (*PaginatedResult[T], error) {
	var total int64

	// Count total records
	if err := db.Model(new(T)).Count(&total).Error; err != nil {
		return nil, err
	}

	// Calculate offset and total pages
	offset := (page - 1) * pageSize
	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))

	// Fetch paginated data
	if err := db.Offset(offset).Limit(pageSize).Find(result).Error; err != nil {
		return nil, err
	}

	return &PaginatedResult[T]{
		Data:       *result,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

// CommonFilters provides common filter patterns for different data types
type CommonFilters struct{}

// NewCommonFilters creates a new instance of CommonFilters
func NewCommonFilters() *CommonFilters {
	return &CommonFilters{}
}

// SearchByName creates a filter for searching by name with CONTAINS
func (cf *CommonFilters) SearchByName(name string) FilterField {
	return FilterField{
		ColumnName: "name",
		Operator:   OperatorContains,
		Value:      name,
	}
}

// FilterByStatus creates a filter for status fields
func (cf *CommonFilters) FilterByStatus(status string) FilterField {
	return FilterField{
		ColumnName: "status",
		Operator:   OperatorEquals,
		Value:      status,
	}
}

// FilterByActive creates a filter for active boolean fields
func (cf *CommonFilters) FilterByActive(active bool) FilterField {
	return FilterField{
		ColumnName: "active",
		Operator:   OperatorEquals,
		Value:      active,
	}
}

// FilterByIDs creates an IN filter for multiple IDs
func (cf *CommonFilters) FilterByIDs(ids []int64) FilterField {
	return FilterField{
		ColumnName: "id",
		Operator:   OperatorIn,
		Value:      ids,
	}
}

// FilterByDateRange creates filters for date range queries
func (cf *CommonFilters) FilterByDateRange(columnName string, from, to interface{}) []FilterField {
	var filters []FilterField

	if from != nil {
		filters = append(filters, FilterField{
			ColumnName: columnName,
			Operator:   OperatorGreaterOrEqual,
			Value:      from,
		})
	}

	if to != nil {
		filters = append(filters, FilterField{
			ColumnName: columnName,
			Operator:   OperatorLessOrEqual,
			Value:      to,
		})
	}

	return filters
}

// FilterByUserID creates a filter for user_id fields
func (cf *CommonFilters) FilterByUserID(userID int64) FilterField {
	return FilterField{
		ColumnName: "user_id",
		Operator:   OperatorEquals,
		Value:      userID,
	}
}

// FilterByProductID creates a filter for product_id fields
func (cf *CommonFilters) FilterByProductID(productID int64) FilterField {
	return FilterField{
		ColumnName: "product_id",
		Operator:   OperatorEquals,
		Value:      productID,
	}
}

// Repository helpers that can be embedded in your repositories

// BaseRepository provides common query functionality
type BaseRepository struct {
	db *gorm.DB
}

// NewBaseRepository creates a new base repository
func NewBaseRepository(db *gorm.DB) *BaseRepository {
	return &BaseRepository{db: db}
}

// QueryBuilder returns a new query builder for this repository
func (br *BaseRepository) QueryBuilder() *QueryBuilder {
	return NewQueryBuilder(br.db)
}

// FindWithFilters finds records with dynamic filters
func (br *BaseRepository) FindWithFilters(result interface{}, filters interface{}) error {
	query := NewQueryBuilder(br.db).
		ApplyFilters(filters).
		Build()

	return query.Find(result).Error
}

// FindOneWithFilters finds a single record with dynamic filters
func (br *BaseRepository) FindOneWithFilters(result interface{}, filters interface{}) error {
	query := NewQueryBuilder(br.db).
		ApplyFilters(filters).
		Build()

	return query.First(result).Error
}

// CountWithFilters counts records with dynamic filters
func (br *BaseRepository) CountWithFilters(model interface{}, filters interface{}) (int64, error) {
	var count int64
	query := NewQueryBuilder(br.db).
		ApplyFilters(filters).
		Build()

	err := query.Model(model).Count(&count).Error
	return count, err
}
