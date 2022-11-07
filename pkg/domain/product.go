package domain

type Product struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type ProductStorage interface {
	// Get(name string) (value []Product, exist bool)
	GetByEmail(key string) (value Product, err error)
	Add(value Product) error
}

type ProductUseCase interface {
	Get(key string) ([]Product, bool)
	Add(value Product) error
}
