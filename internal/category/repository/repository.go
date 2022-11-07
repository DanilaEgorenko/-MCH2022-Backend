package categoryRepository

import (
	"context"
	"log"
	"snakealive/m/pkg/domain"

	"github.com/jackc/pgx/v4/pgxpool"
)

type CategoryStorage struct {
	dataHolder *pgxpool.Pool
}

func NewCategoryStorage(DB *pgxpool.Pool) domain.CategoryStorage {
	return &CategoryStorage{dataHolder: DB}
}

func (u *CategoryStorage) GetCategoryById(key string) (value domain.Category, err error) {
	var category domain.Category

	conn, err := u.dataHolder.Acquire(context.Background())
	if err != nil {
		log.Printf("Error while getting user")
		return category, err
	}
	defer conn.Release()
	log.Printf("key = ", key)
	err = conn.QueryRow(context.Background(),
		`SELECT cp.id, cp.email, cp.name, cp.legal_name, cp.itn, cp.psrn, cp.adress, cp.phone, cp.link, cp.category_id
		FROM Company AS cp
		JOIN Category AS ct ON cp.category_id = ct.id
		WHERE cp.category_id = $1`,
		key,
	).Scan(&category.Id, &category.Title, &category.IndustryId)

	return category, err
}

func (u *CategoryStorage) GetCategoriesInIndustry(key string) (value domain.Categories, err error) {
	var industry domain.Industry
	var category domain.Category
	categories := make([]domain.Category, 0)
	conn, err := u.dataHolder.Acquire(context.Background())
	if err != nil {
		log.Printf("Error while getting categories")
		return categories, err
	}
	defer conn.Release()
	err = conn.QueryRow(context.Background(),
		`SELECT id, title
		FROM Industry WHERE id = $1`,
		key,
	).Scan(&industry.Id, &industry.Title)
	if err != nil {
		log.Printf("Error while getting categories")
		return categories, err
	}
	rows, err := conn.Query(context.Background(),
		`SELECT id, title, industry_id
		FROM Category WHERE industry_id = $1;`,
		industry.Id)
	if err != nil {
		log.Printf("Error while getting categories")
		return categories, err
	}
	for rows.Next() {
		rows.Scan(&category.Id, &category.Title, &category.IndustryId)
		categories = append(categories, category)
	}
	if rows.Err() != nil {
		log.Printf("Error while scanning categories", rows.Err())
		return categories, err
	}
	if len(categories) == 0 {
		return categories, err
	}
	return categories, err
}
