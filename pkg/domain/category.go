package domain

type Category struct {
	Id         int    `json:"id"`
	Title      string `json:"title"`
	IndustryId int    `json:"industry_id"`
}
type Categories []Category

type CategoryStorage interface {
	GetCategoryById(key string) (value Category, err error)
	GetCategoriesInIndustry(key string) (value Categories, err error)
}

type CategoryUseCase interface {
	GetCategoryById(key string) (value []byte, err error)
	GetCategoriesInIndustry(key string) (value []byte, err error)
}
