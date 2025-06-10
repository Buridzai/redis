package models

type Product struct {
	ID          uint    `json:"id" gorm:"primaryKey"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Image       string  `json:"image"`
	CreatedAt   int64   `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   int64   `json:"updated_at" gorm:"autoUpdateTime"`
}
