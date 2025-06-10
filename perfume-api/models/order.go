package models

type Order struct {
	ID     uint        `json:"id" gorm:"primaryKey"`
	UserID uint        `json:"user_id"`
	Total  float64     `json:"total"` // sửa từ int -> float64
	Items  []OrderItem `json:"items" gorm:"foreignKey:OrderID"`
}
