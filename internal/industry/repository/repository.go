package industryRepository

import (
	"context"
	"log"
	"snakealive/m/pkg/domain"

	"github.com/jackc/pgx/v4/pgxpool"
)

type IndustryStorage struct {
	dataHolder *pgxpool.Pool
}

func NewIndustryStorage(DB *pgxpool.Pool) domain.IndustryStorage {
	return &IndustryStorage{dataHolder: DB}
}

func (u *IndustryStorage) GetAllIndustries() (value domain.Industries, err error) {
	var industry domain.Industry
	industries := make([]domain.Industry, 0)
	conn, err := u.dataHolder.Acquire(context.Background())
	if err != nil {
		log.Printf("Error while getting industries")
		return industries, err
	}
	defer conn.Release()
	rows, err := conn.Query(context.Background(),
		`SELECT id, title FROM Industry`)
	if err != nil {
		log.Printf("Error while getting industries")
		return industries, err
	}
	for rows.Next() {
		rows.Scan(&industry.Id, &industry.Title)
		industries = append(industries, industry)
	}
	if rows.Err() != nil {
		log.Printf("Error while scanning industries", rows.Err())
		return industries, err
	}
	if len(industries) == 0 {
		return industries, err
	}
	return industries, err

}
