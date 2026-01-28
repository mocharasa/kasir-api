package models

type Product struct {
	ID           int    `json:"id"`
	CategoryID   int    `json:"category_id"`
	CategoryName string `json:"category_name,omitempty"`
	Name         string `json:"name"`
	Price        int    `json:"price"`
	Stock        int    `json:"stock"`
}
