package query

import (
	"fmt"
	"reflect"
	"strings"

	"gorm.io/gorm"
)

// Operator represents the type of filter operation
type Operator string

const (
	OperatorEquals         Operator = "="
	OperatorNotEquals      Operator = "!="
	OperatorContains       Operator = "CONTAINS"
	OperatorGreaterThan    Operator = ">"
	OperatorLessThan       Operator = "<"
	OperatorGreaterOrEqual Operator = ">="
	OperatorLessOrEqual    Operator = "<="
	OperatorIn             Operator = "IN"
	OperatorNotIn          Operator = "NOT IN"
	OperatorIsNull         Operator = "IS NULL"
	OperatorIsNotNull      Operator = "IS NOT NULL"
	OperatorStartsWith     Operator = "STARTS_WITH"
	OperatorEndsWith       Operator = "ENDS_WITH"
)

// FilterField represents a field filter with column name, operator, and value
type FilterField struct {
	ColumnName string
	Operator   Operator
	Value      interface{}
}

// PreloadConfig represents configuration for a preloaded relationship
type PreloadConfig struct {
	// Relationship is the name of the relation to preload (e.g., "User", "Product")
	Relationship string

	// Conditions allows adding conditions to the preload
	Conditions []interface{}

	// CustomPreload is an optional function that allows custom preloading behavior
	// If nil, standard preloading will be used
	CustomPreload func(db *gorm.DB) *gorm.DB
}

// SortConfig represents sorting configuration
type SortConfig struct {
	Field string
	Order SortOrder
}

type SortOrder string

const (
	SortOrderAsc  SortOrder = "ASC"
	SortOrderDesc SortOrder = "DESC"
)

// PaginationConfig represents pagination configuration
type PaginationConfig struct {
	Page     int
	PageSize int
}

// QueryBuilder provides a fluent interface for building complex queries
type QueryBuilder struct {
	db         *gorm.DB
	filters    []FilterField
	preloads   []PreloadConfig
	sorts      []SortConfig
	pagination *PaginationConfig
	distinct   bool
	groupBy    []string
	having     []interface{}
}

// NewQueryBuilder creates a new query builder instance
func NewQueryBuilder(db *gorm.DB) *QueryBuilder {
	return &QueryBuilder{
		db: db,
	}
}

// ApplyFilters applies filters to a GORM query dynamically based on filter struct
// It accepts a struct with field tags that specify how to apply the filter
// Tag format: `filter:"column_name,operator"` where operator is optional and defaults to equals
func (qb *QueryBuilder) ApplyFilters(filterStruct interface{}) *QueryBuilder {
	if filterStruct == nil {
		return qb
	}

	val := reflect.ValueOf(filterStruct)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	// If not a struct, return the query builder as is
	if val.Kind() != reflect.Struct {
		return qb
	}

	typ := val.Type()
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)

		// Skip unexported fields
		if !field.CanInterface() {
			continue
		}

		// Try to get the filter tag
		filterTag := fieldType.Tag.Get("filter")
		if filterTag == "" {
			// If no filter tag, use the field name as column name and default to equals
			filterTag = toSnakeCase(fieldType.Name)
		}

		// Parse the filter tag
		parts := strings.Split(filterTag, ",")
		columnName := parts[0]

		// Default operator is equals
		operator := OperatorEquals
		if len(parts) > 1 {
			operator = Operator(strings.TrimSpace(parts[1]))
		}

		// Handle different field types and skip empty values
		fieldValue := field.Interface()
		if !isZeroValue(field) {
			qb.filters = append(qb.filters, FilterField{
				ColumnName: columnName,
				Operator:   operator,
				Value:      fieldValue,
			})
		}
	}

	return qb
}

// AddFilter adds a single filter to the query builder
func (qb *QueryBuilder) AddFilter(columnName string, operator Operator, value interface{}) *QueryBuilder {
	qb.filters = append(qb.filters, FilterField{
		ColumnName: columnName,
		Operator:   operator,
		Value:      value,
	})
	return qb
}

// AddFilters adds multiple filters to the query builder
func (qb *QueryBuilder) AddFilters(filters []FilterField) *QueryBuilder {
	qb.filters = append(qb.filters, filters...)
	return qb
}

// AddPreload adds a preload configuration
func (qb *QueryBuilder) AddPreload(relationship string, conditions ...interface{}) *QueryBuilder {
	qb.preloads = append(qb.preloads, PreloadConfig{
		Relationship: relationship,
		Conditions:   conditions,
	})
	return qb
}

// AddCustomPreload adds a custom preload configuration
func (qb *QueryBuilder) AddCustomPreload(relationship string, customFunc func(db *gorm.DB) *gorm.DB) *QueryBuilder {
	qb.preloads = append(qb.preloads, PreloadConfig{
		Relationship:  relationship,
		CustomPreload: customFunc,
	})
	return qb
}

// AddSort adds sorting configuration
func (qb *QueryBuilder) AddSort(field string, order SortOrder) *QueryBuilder {
	qb.sorts = append(qb.sorts, SortConfig{
		Field: field,
		Order: order,
	})
	return qb
}

// SetPagination sets pagination configuration
func (qb *QueryBuilder) SetPagination(page, pageSize int) *QueryBuilder {
	qb.pagination = &PaginationConfig{
		Page:     page,
		PageSize: pageSize,
	}
	return qb
}

// SetDistinct enables distinct selection
func (qb *QueryBuilder) SetDistinct(distinct bool) *QueryBuilder {
	qb.distinct = distinct
	return qb
}

// AddGroupBy adds group by clause
func (qb *QueryBuilder) AddGroupBy(fields ...string) *QueryBuilder {
	qb.groupBy = append(qb.groupBy, fields...)
	return qb
}

// AddHaving adds having clause
func (qb *QueryBuilder) AddHaving(query interface{}, args ...interface{}) *QueryBuilder {
	qb.having = append(qb.having, query)
	qb.having = append(qb.having, args...)
	return qb
}

// Build applies all configurations and returns the final query
func (qb *QueryBuilder) Build() *gorm.DB {
	query := qb.db

	// Apply distinct
	if qb.distinct {
		query = query.Distinct()
	}

	// Apply filters
	for _, filter := range qb.filters {
		query = applyFilter(query, filter)
	}

	// Apply preloads
	for _, preload := range qb.preloads {
		if preload.CustomPreload != nil {
			query = preload.CustomPreload(query)
		} else if len(preload.Conditions) > 0 {
			query = query.Preload(preload.Relationship, preload.Conditions...)
		} else {
			query = query.Preload(preload.Relationship)
		}
	}

	// Apply group by
	if len(qb.groupBy) > 0 {
		query = query.Group(strings.Join(qb.groupBy, ", "))
	}

	// Apply having
	if len(qb.having) > 0 {
		query = query.Having(qb.having[0], qb.having[1:]...)
	}

	// Apply sorting
	for _, sort := range qb.sorts {
		query = query.Order(fmt.Sprintf("%s %s", sort.Field, sort.Order))
	}

	// Apply pagination
	if qb.pagination != nil {
		offset := (qb.pagination.Page - 1) * qb.pagination.PageSize
		query = query.Offset(offset).Limit(qb.pagination.PageSize)
	}

	return query
}

// applyFilter applies a single filter to the query
func applyFilter(query *gorm.DB, filter FilterField) *gorm.DB {
	switch filter.Operator {
	case OperatorContains:
		// Use database-agnostic LIKE with % wildcards
		return query.Where(fmt.Sprintf("%s LIKE ?", filter.ColumnName), "%"+filter.Value.(string)+"%")
	case OperatorStartsWith:
		return query.Where(fmt.Sprintf("%s LIKE ?", filter.ColumnName), filter.Value.(string)+"%")
	case OperatorEndsWith:
		return query.Where(fmt.Sprintf("%s LIKE ?", filter.ColumnName), "%"+filter.Value.(string))
	case OperatorIn:
		return query.Where(fmt.Sprintf("%s IN ?", filter.ColumnName), filter.Value)
	case OperatorNotIn:
		return query.Where(fmt.Sprintf("%s NOT IN ?", filter.ColumnName), filter.Value)
	case OperatorIsNull:
		return query.Where(fmt.Sprintf("%s IS NULL", filter.ColumnName))
	case OperatorIsNotNull:
		return query.Where(fmt.Sprintf("%s IS NOT NULL", filter.ColumnName))
	default:
		// For simple operators (=, !=, >, <, >=, <=)
		return query.Where(fmt.Sprintf("%s %s ?", filter.ColumnName, filter.Operator), filter.Value)
	}
}

// isZeroValue checks if a reflect.Value is the zero value for its type
func isZeroValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Slice, reflect.Map, reflect.Chan:
		return v.IsNil() || v.Len() == 0
	case reflect.String:
		return v.String() == ""
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	default:
		return reflect.DeepEqual(v.Interface(), reflect.Zero(v.Type()).Interface())
	}
}

// toSnakeCase converts CamelCase to snake_case
func toSnakeCase(str string) string {
	var result strings.Builder
	for i, r := range str {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result.WriteRune('_')
		}
		result.WriteRune(r)
	}
	return strings.ToLower(result.String())
}
