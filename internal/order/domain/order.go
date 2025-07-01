package domain

import (
	productDomain "github.com/mohsenjafari-aiio/aiiobackend/internal/product/domain"
	userDomain "github.com/mohsenjafari-aiio/aiiobackend/internal/user/domain"
)

type Order struct {
	ID        int64 `gorm:"primaryKey"`
	UserID    int64
	User      userDomain.User `gorm:"foreignKey:UserID"`
	ProductID int64
	Product   productDomain.Product `gorm:"foreignKey:ProductID"`
	Quantity  int
	Status    string `gorm:"type:varchar(20);not null"`
}

func NewOrder(userID, productID int64, quantity int) *Order {
	return &Order{
		UserID:    userID,
		ProductID: productID,
		Quantity:  quantity,
		Status:    "PENDING",
	}
}

func (o *Order) Confirm() {
	o.Status = "CONFIRMED"
}
