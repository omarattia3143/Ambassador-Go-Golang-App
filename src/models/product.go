package models

type Product struct {
	Id          uint    `json:"id" gorm:"constraint:OnDelete:CASCADE"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Image       string  `json:"image"`
	Price       float64 `json:"price"`
	Links       []Link  `gorm:"many2many:link_products;constraint:OnDelete:CASCADE"`
}
