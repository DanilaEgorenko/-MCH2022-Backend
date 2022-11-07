package categoryUseCase

import (
	"log"
	chttp "snakealive/m/pkg/customhttp"
	"snakealive/m/pkg/domain"
)

func NewCategoryUseCase(categoryStorage domain.CategoryStorage) domain.CategoryUseCase {
	return categoryUseCase{categoryStorage: categoryStorage}
}

type categoryUseCase struct {
	categoryStorage domain.CategoryStorage
}

func (c categoryUseCase) GetCategoryById(key string) (value []byte, err error) {
	category, err := c.categoryStorage.GetCategoryById(key)
	if err != nil {
		return []byte{}, err
	}
	bytes, err := chttp.ApiResp(category)
	if err != nil {
		log.Printf("error while marshalling JSON: %s", err)
	}
	return bytes, err
}

func (c categoryUseCase) GetCategoriesInIndustry(key string) (value []byte, err error) {
	category, err := c.categoryStorage.GetCategoriesInIndustry(key)
	if err != nil {
		return []byte{}, err
	}
	bytes, err := chttp.ApiResp(category)
	if err != nil {
		log.Printf("error while marshalling JSON: %s", err)
	}
	return bytes, err
}
