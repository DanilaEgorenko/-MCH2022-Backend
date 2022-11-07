package domain

type Company struct {
	Id         int      `json:"id"`
	Email      string   `json:"email"`
	Password   string   `json:"password"`
	Name       string   `json:"name"`
	LegalName  string   `json:"legal_name"`
	Itn        string   `json:"itn"`
	Psrn       string   `json:"psrn"`
	Adress     string   `json:"adress"`
	Phone      string   `json:"phone"`
	Link       string   `json:"link"`
	CategoryId int      `json:"category_id"`
	Docks      []string `json:"docs"`
}
type CompanySearch struct {
	Name string `json:"name"`
}

type Companies []Company

type CompanyStorage interface {
	SearchCompanies(key string) (value Companies, err error)
	GetByEmail(key string) (value Company, err error)
	GetCompanyById(key string) (value Company, err error)
	GetCompaniesByCategoryId(key string) (value Companies, err error)
	Add(value Company) error
}

type CompanyUseCase interface {
	SearchCompanies(key CompanySearch) (value []byte, err error)
	GetByEmail(key string) (value Company, err error)
	GetCompanyById(key string) (value []byte, err error)
	GetCompaniesByCategoryId(key string) (value []byte, err error)
	Add(value Company) error
	Validate(company *Company) bool
	Login(company *Company) (int, error)
	Registration(company *Company) (int, error)
}
