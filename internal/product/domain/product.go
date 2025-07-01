package domain

import "errors"

type Product struct {
	ID    int64  `gorm:"primaryKey"`
	Name  string `gorm:"not null"`
	Stock int
}

func (p *Product) Reserve(qty int) error {
	if p.Stock < qty {
		return errors.New("insufficient stock")
	}
	p.Stock -= qty
	return nil
}
