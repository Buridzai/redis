package models

type Cart struct {
	ID        uint    `json:"id" gorm:"primaryKey"`
	UserID    uint    `json:"user_id"`
	ProductID uint    `json:"product_id"`
	Quantity  int     `json:"quantity"`
	Product   Product `json:"product" gorm:"foreignKey:ProductID"`
}
