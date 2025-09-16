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

type SearchProductParameter struct {
	Name     string  `json:"name"`
	Category string  `json:"category"`
	MinPrice float64 `json:"minPrice"`
	MaxPrice float64 `json:"maxPrice"`
	Page     int     `json:"page"`
	PageSize int     `json:"pageSize"`
	OrderBy  string  `json:"orderBy"`
	Sort     string  `json:"sort"`
}

type SearchProductResponse struct {
	Products    []Product `json:"products"`
	Page        int       `json:"page"`
	PageSize    int       `json:"pageSize"`
	TotalCount  int       `json:"totalCount"`
	TotalPages  int       `json:"totalPages"`
	NextPageUrl *string   `json:"nextPageUrl"`
}
