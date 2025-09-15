package service

import "product/cmd/product/repository"

type ProductService struct {
	ProductRepo repository.ProductRepository
}

func NewProductService(productRepo repository.ProductRepository) *ProductService {
	return &ProductService{
		ProductRepo: productRepo,
	}
}
