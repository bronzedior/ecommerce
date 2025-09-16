package service

import (
	"context"
	"product/cmd/product/repository"
	"product/infrastructure/log"
	"product/models"

	"github.com/sirupsen/logrus"
)

type ProductService struct {
	ProductRepository repository.ProductRepository
}

func NewProductService(ProductRepository repository.ProductRepository) *ProductService {
	return &ProductService{
		ProductRepository: ProductRepository,
	}
}

func (s *ProductService) GetProductByID(ctx context.Context, productID int64) (*models.Product, error) {
	// get from Redis
	product, err := s.ProductRepository.GetProductByIDFromRedis(ctx, productID)
	if err != nil {
		log.Logger.WithFields(logrus.Fields{
			"productID": productID,
		}).Errorf("s.ProductRepository.GetProductByIDFromRedis got error %v", err)
	}

	if product.ID != 0 {
		return product, nil
	}

	// get from DB
	product, err = s.ProductRepository.FindProductByID(ctx, productID)
	if err != nil {
		return nil, err
	}

	ctxConcurrent := context.WithValue(ctx, context.Background(), ctx.Value("request_id"))
	go func(ctx context.Context, product *models.Product, productID int64) {
		errConcurrent := s.ProductRepository.SetProductByID(ctx, product, productID)
		if errConcurrent != nil {
			log.Logger.WithFields(logrus.Fields{
				"product": product,
			}).Errorf("s.ProductRepository.SetProductByID() got error %v", errConcurrent)
		}
	}(ctxConcurrent, product, productID)

	return product, nil
}

func (s *ProductService) GetProductCategoryByID(ctx context.Context, productCategoryID int) (*models.ProductCategory, error) {
	// get from Redis
	productCategory, err := s.ProductRepository.GetProductCategoryByIDFromRedis(ctx, productCategoryID)
	if err != nil {
		log.Logger.WithFields(logrus.Fields{
			"productCategoryID": productCategoryID,
		}).Errorf("s.ProductRepository.GetProductCategoryByIDFromRedis got error %v", err)
	}

	if productCategory.ID != 0 {
		return productCategory, nil
	}

	// get from DB
	productCategory, err = s.ProductRepository.FindProductCategoryByID(ctx, productCategoryID)
	if err != nil {
		return nil, err
	}

	ctxConcurrent := context.WithValue(ctx, context.Background(), ctx.Value("request_id"))
	go func(ctx context.Context, product *models.ProductCategory, productID int) {
		errConcurrent := s.ProductRepository.SetProductCategoryByID(ctx, productCategory, productCategoryID)
		if errConcurrent != nil {
			log.Logger.WithFields(logrus.Fields{
				"productCategory": productCategory,
			}).Errorf("s.ProductRepository.SetProductCategoryByID() got error %v", errConcurrent)
		}
	}(ctxConcurrent, productCategory, productCategoryID)

	return productCategory, nil
}

func (s *ProductService) CreateNewProduct(ctx context.Context, param *models.Product) (int64, error) {
	productID, err := s.ProductRepository.InsertNewProduct(ctx, param)
	if err != nil {
		return 0, err
	}

	return productID, nil
}

func (s *ProductService) CreateNewProductCategory(ctx context.Context, param *models.ProductCategory) (int, error) {
	productCategoryID, err := s.ProductRepository.InsertNewProductCategory(ctx, param)
	if err != nil {
		return 0, err
	}

	return productCategoryID, nil
}

func (s *ProductService) EditProduct(ctx context.Context, product *models.Product) (*models.Product, error) {
	product, err := s.ProductRepository.UpdateProduct(ctx, product)
	if err != nil {
		return nil, err
	}

	return product, nil
}

func (s *ProductService) EditProductCategory(ctx context.Context, productCategory *models.ProductCategory) (*models.ProductCategory, error) {
	productCategory, err := s.ProductRepository.UpdateProductCategory(ctx, productCategory)
	if err != nil {
		return nil, err
	}

	return productCategory, nil
}

func (s *ProductService) DeleteProduct(ctx context.Context, productID int64) error {
	err := s.ProductRepository.DeleteProduct(ctx, productID)
	if err != nil {
		return err
	}

	return nil
}

func (s *ProductService) DeleteProductCategory(ctx context.Context, productCategoryID int) error {
	err := s.ProductRepository.DeleteProductCategory(ctx, int64(productCategoryID))
	if err != nil {
		return err
	}

	return nil
}

func (s *ProductService) SearchProduct(ctx context.Context, param models.SearchProductParameter) ([]models.Product, int, error) {
	products, totalCount, err := s.ProductRepository.SearchProduct(ctx, param)
	if err != nil {
		return nil, 0, err
	}

	return products, totalCount, nil
}
