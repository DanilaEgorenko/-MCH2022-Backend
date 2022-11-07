package companyUseCase

import (
	"log"
	chttp "snakealive/m/pkg/customhttp"
	"snakealive/m/pkg/domain"

	"github.com/asaskevich/govalidator"
	"github.com/valyala/fasthttp"
)

func NewCompanyUseCase(companyStorage domain.CompanyStorage) domain.CompanyUseCase {
	return companyUseCase{companyStorage: companyStorage}
}

type companyUseCase struct {
	companyStorage domain.CompanyStorage
}

func (c companyUseCase) SearchCompanies(key domain.CompanySearch) (value []byte, err error) {
	companies, err := c.companyStorage.SearchCompanies(key.Name)
	if err != nil {
		return []byte{}, err
	}
	bytes, err := chttp.ApiResp(companies)
	if err != nil {
		log.Printf("error while marshalling JSON: %s", err)
	}
	return bytes, err
}

func (c companyUseCase) GetByEmail(key string) (value domain.Company, err error) {
	return c.companyStorage.GetByEmail(key)
}

func (c companyUseCase) GetCompanyById(key string) (value []byte, err error) {
	company, err := c.companyStorage.GetCompanyById(key)
	if err != nil {
		return []byte{}, err
	}
	bytes, err := chttp.ApiResp(company)
	if err != nil {
		log.Printf("error while marshalling JSON: %s", err)
	}
	return bytes, err
}

func (c companyUseCase) GetCompaniesByCategoryId(key string) (value []byte, err error) {
	companies, err := c.companyStorage.GetCompaniesByCategoryId(key)
	if err != nil {
		return []byte{}, err
	}
	bytes, err := chttp.ApiResp(companies)
	if err != nil {
		log.Printf("error while marshalling JSON: %s", err)
	}
	return bytes, err
}

func (c companyUseCase) Add(user domain.Company) error {
	return c.companyStorage.Add(user)
}

func (c companyUseCase) Validate(company *domain.Company) bool {
	// if !govalidator.IsEmail(company.Email) ||
	// 	!govalidator.StringLength(company.Password, cnst.MinPasswordLength, cnst.MaxPasswordLength) ||
	// 	!govalidator.MaxStringLength(company.Email, cnst.MaxEmailLength) {
	// 	return false
	// }
	return true
}

// func (u CompanyUseCase) Update(id int, updatedUser domain.User) error {

// 	user, err := u.GetByEmail(updatedUser.Email)
// 	if err == nil && user.Id != id {
// 		return errors.New("user with this email already exists") // change later
// 	}

// 	return u.companyStorage.Update(id, updatedUser)
// }

// func (u CompanyUseCase) Delete(id int) error {
// 	return u.companyStorage.Delete(id)
// }

func (u companyUseCase) Login(user *domain.Company) (int, error) {
	foundUser, err := u.GetByEmail(user.Email)
	if err != nil {
		log.Printf("error while login-GetByEmail")
		log.Print(err)
		return fasthttp.StatusNotFound, err
	}

	if foundUser.Password != user.Password {
		return fasthttp.StatusBadRequest, err
	}

	return fasthttp.StatusOK, err
}

func (u companyUseCase) Registration(company *domain.Company) (int, error) {
	_, err := govalidator.ValidateStruct(company)
	if err != nil {
		log.Printf("error while validating user")
		return fasthttp.StatusBadRequest, err
	}
	log.Printf(company.Email)

	_, err = u.GetByEmail(company.Email)
	if err == nil {
		log.Printf("user with this email already exists")
		return fasthttp.StatusBadRequest, err
	}

	err = u.Add(*company)
	if err != nil {
		log.Printf("error while adding user")
		return fasthttp.StatusBadRequest, err
	}

	return fasthttp.StatusOK, err
}
