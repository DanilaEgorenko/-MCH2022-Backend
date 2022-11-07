package domain

type Industry struct {
	Id    int    `json:"id"`
	Title string `json:"title"`
}

type Industries []Industry

type IndustryStorage interface {
	GetAllIndustries() (value Industries, err error)
}

type IndustryUseCase interface {
	GetAllIndustries() (value []byte, err error)
}
