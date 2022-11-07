package industryUseCase

import (
	"log"
	chttp "snakealive/m/pkg/customhttp"
	"snakealive/m/pkg/domain"
)

func NewIndustryUseCase(industryStorage domain.IndustryStorage) domain.IndustryUseCase {
	return industryUseCase{industryStorage: industryStorage}
}

type industryUseCase struct {
	industryStorage domain.IndustryStorage
}

func (c industryUseCase) GetAllIndustries() (value []byte, err error) {
	category, err := c.industryStorage.GetAllIndustries()
	if err != nil {
		return []byte{}, err
	}
	bytes, err := chttp.ApiResp(category)
	if err != nil {
		log.Printf("error while marshalling JSON: %s", err)
	}
	return bytes, err
}
