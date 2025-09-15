package usecase

import "product/cmd/product/service"

type ProductUsecase struct {
	ProductService service.ProductService
}

func NewProductUsecase(productService service.ProductService) *ProductUsecase {
	return &ProductUsecase{
		ProductService: productService,
	}
}
