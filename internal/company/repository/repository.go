package companyRepository

import (
	"context"
	"log"
	"snakealive/m/pkg/domain"

	"github.com/jackc/pgx/v4/pgxpool"
)

type CompanyStorage struct {
	dataHolder *pgxpool.Pool
}

func NewCompanyStorage(DB *pgxpool.Pool) domain.CompanyStorage {
	return &CompanyStorage{dataHolder: DB}
}

func (u *CompanyStorage) Add(value domain.Company) error {
	conn, err := u.dataHolder.Acquire(context.Background())
	if err != nil {
		log.Printf("Connection error while adding user ", err)
		return err
	}
	defer conn.Release()

	_, err = conn.Exec(context.Background(),
		`INSERT INTO Company ("email", "password", "name", "legal_name", "itn", "psrn", "adress", "phone", "link", "category_id" ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`,

		value.Email,
		value.Password,
		value.Name,
		value.LegalName,
		value.Itn,
		value.Psrn,
		value.Adress,
		value.Phone,
		value.Link,
		value.CategoryId,
	)
	return err
}

func (u *CompanyStorage) GetByEmail(key string) (value domain.Company, err error) {
	var company domain.Company

	conn, err := u.dataHolder.Acquire(context.Background())
	if err != nil {
		log.Printf("Error while getting user")
		return company, err
	}
	defer conn.Release()

	err = conn.QueryRow(context.Background(),
		`SELECT id, name, legal_name, password, email
		FROM Company WHERE email = $1`,
		key,
	).Scan(&company.Id, &company.Name, &company.LegalName, &company.Password, &company.Email)

	return company, err
}
func (u *CompanyStorage) SearchCompanies(key string) (value domain.Companies, err error) {
	companies := make([]domain.Company, 0)
	param := "^" + key
	conn, err := u.dataHolder.Acquire(context.Background())
	if err != nil {
		log.Printf("Connection error while adding trip ", err)
		return companies, err
	}
	defer conn.Release()
	rows, err := conn.Query(context.Background(),
		`SELECT cp.id, cp.email, cp.name, cp.legal_name, cp.itn, cp.psrn, cp.adress, cp.phone, cp.link, cp.category_id
	FROM Company AS cp
	WHERE cp.name ~ $1`,
		param)
	if err != nil {
		log.Printf("Error while getting companies")
		return companies, err
	}
	var cp domain.Company
	for rows.Next() {
		rows.Scan(&cp.Id, &cp.Email, &cp.Name, &cp.LegalName, &cp.Itn, &cp.Psrn, &cp.Adress, &cp.Phone, &cp.Link, &cp.CategoryId)
		companies = append(companies, cp)
	}
	if rows.Err() != nil {
		log.Printf("Error while scanning places", rows.Err())
		return companies, err
	}
	if len(companies) == 0 {
		return companies, err
	}
	return companies, err
}

func (u *CompanyStorage) GetCompanyById(key string) (value domain.Company, err error) {
	var company domain.Company

	conn, err := u.dataHolder.Acquire(context.Background())
	if err != nil {
		log.Printf("Error while getting user")
		return company, err
	}
	defer conn.Release()
	log.Printf("key = ", key)
	err = conn.QueryRow(context.Background(),
		`SELECT id, email, name, legal_name, itn, psrn, adress, phone, link, category_id
		FROM Company WHERE id = $1`,
		key,
	).Scan(&company.Id, &company.Email, &company.Name, &company.LegalName, &company.Itn, &company.Psrn, &company.Adress, &company.Phone, &company.Link, &company.CategoryId)

	return company, err
}

// func (u *CompanyStorage) GetCompaniesByCategoryTitle(key string) (value domain.Companies, err error) {
// 	var companies domain.Companies

// 	conn, err := u.dataHolder.Acquire(context.Background())
// 	if err != nil {
// 		log.Printf("Error while getting companies")
// 		return companies, err
// 	}
// 	defer conn.Release()
// 	err = conn.QueryRow(context.Background(),
// 		`SELECT cp.id, cp.email, cp.name, cp.legal_name, cp.itn, cp.psrn, cp.adress, cp.phone, cp.link, cp.category_id
// 		FROM Company AS cp
// 		JOIN Category AS ct ON cp.category_id = ct.id
// 		WHERE cp.category_id = $1`,
// 		key,
// 	).Scan(&category.Id, &category.Title, &category.IndustryId)

// 	return category, err
// }

func (u *CompanyStorage) GetCompaniesByCategoryId(key string) (value domain.Companies, err error) {
	companies := make([]domain.Company, 0)
	var category domain.Category

	conn, err := u.dataHolder.Acquire(context.Background())
	if err != nil {
		log.Printf("Connection error while adding trip ", err)
		return companies, err
	}
	defer conn.Release()

	err = conn.QueryRow(context.Background(),
		`SELECT id, title, industry_id
		FROM Category WHERE id = $1`,
		key,
	).Scan(&category.Id, &category.Title, &category.IndustryId)
	if err != nil {
		log.Printf("Error while getting category")
		return companies, err
	}

	rows, err := conn.Query(context.Background(),
		`SELECT cp.id, cp.email, cp.name, cp.legal_name, cp.itn, cp.psrn, cp.adress, cp.phone, cp.link, cp.category_id
	FROM Company AS cp
	JOIN Category AS ct ON cp.category_id = ct.id
	WHERE cp.category_id = $1`,
		category.Id)
	if err != nil {
		log.Printf("Error while getting companies")
		return companies, err
	}
	var cp domain.Company
	for rows.Next() {
		rows.Scan(&cp.Id, &cp.Email, &cp.Name, &cp.LegalName, &cp.Itn, &cp.Psrn, &cp.Adress, &cp.Phone, &cp.Link, &cp.CategoryId)
		companies = append(companies, cp)
	}
	if rows.Err() != nil {
		log.Printf("Error while scanning places", rows.Err())
		return companies, err
	}
	if len(companies) == 0 {
		return companies, err
	}
	return companies, err
}
