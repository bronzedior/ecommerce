package models

type Product struct {
	ID          int64   `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Stock       int     `json:"stock"`
	CategoryID  int     `json:"category_id"`
}

type ProductManagementParameter struct {
	Action string `json:"action"`
	Product
}

type ProductCategory struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type ProductCategoryManagementParameter struct {
	Action string `json:"action"`
	ProductCategory
}
