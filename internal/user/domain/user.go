package domain

type User struct {
	ID     int64  `gorm:"primaryKey"`
	Active bool   `gorm:"not null"`
	Email  string `gorm:"uniqueIndex;not null"`
}

func (u *User) Activate() {
	u.Active = true
}
