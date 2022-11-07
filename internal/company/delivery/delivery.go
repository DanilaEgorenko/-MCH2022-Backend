package companyDelivery

import (
	"encoding/json"
	"fmt"
	"log"
	cnst "snakealive/m/pkg/constants"
	"snakealive/m/pkg/domain"

	cr "snakealive/m/internal/company/repository"
	cu "snakealive/m/internal/company/usecase"
	ccd "snakealive/m/internal/cookie/delivery"

	"github.com/fasthttp/router"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/valyala/fasthttp"
)

type CompanyHandler interface {
	Login(ctx *fasthttp.RequestCtx)
	Registration(ctx *fasthttp.RequestCtx)
	Logout(ctx *fasthttp.RequestCtx)
	GetCompanyById(ctx *fasthttp.RequestCtx)
	GetCompaniesByCategoryId(ctx *fasthttp.RequestCtx)
	SearchCompanies(ctx *fasthttp.RequestCtx)
}

type companyHandler struct {
	CompanyUseCase domain.CompanyUseCase
	CookieHandler  ccd.CookieHandler
}

func NewCompanyHandler(CompanyUseCase domain.CompanyUseCase, CookieHandler ccd.CookieHandler) CompanyHandler {
	return &companyHandler{
		CompanyUseCase: CompanyUseCase,
		CookieHandler:  CookieHandler,
	}
}

func CreateDelivery(db *pgxpool.Pool) CompanyHandler {
	cookieLayer := ccd.CreateDelivery(db)
	userLayer := NewCompanyHandler(cu.NewCompanyUseCase(cr.NewCompanyStorage(db)), cookieLayer)
	return userLayer
}

func SetUpCompanyRouter(db *pgxpool.Pool, r *router.Router) *router.Router {
	companyHandler := CreateDelivery(db)
	r.POST(cnst.LoginURL, companyHandler.Login)
	r.POST(cnst.RegisterURL, companyHandler.Registration)
	r.GET(cnst.CompanyURL, companyHandler.GetCompanyById)
	r.GET(cnst.CategoryURL, companyHandler.GetCompaniesByCategoryId)
	r.POST(cnst.CompanySearchURL, companyHandler.SearchCompanies)
	return r
}

func (s *companyHandler) Login(ctx *fasthttp.RequestCtx) {
	company := new(domain.Company)
	if err := json.Unmarshal(ctx.PostBody(), &company); err != nil {
		log.Printf("error while unmarshalling JSON: %s", err)
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}
	valid := s.CompanyUseCase.Validate(company)
	if !valid {
		log.Printf("error while validating user")
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}

	code, err := s.CompanyUseCase.Login(company)
	ctx.SetStatusCode(code)
	if err != nil {
		log.Printf("error while logging user in")
		return
	}

	с := fmt.Sprint(uuid.NewMD5(uuid.UUID{}, []byte(company.Email)))
	foundUser, _ := s.CompanyUseCase.GetByEmail(company.Email)
	s.CookieHandler.SetCookieAndToken(ctx, с, foundUser.Id)
}

func (s *companyHandler) Registration(ctx *fasthttp.RequestCtx) {
	company := new(domain.Company)

	if err := json.Unmarshal(ctx.PostBody(), &company); err != nil {
		log.Printf("error while unmarshalling JSON: %s", err)
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}

	code, err := s.CompanyUseCase.Registration(company)
	ctx.SetStatusCode(code)
	if err != nil {
		log.Printf("error while registering user in", err)
		return
	}

	с := fmt.Sprint(uuid.NewMD5(uuid.UUID{}, []byte(company.Email)))
	newUser, _ := s.CompanyUseCase.GetByEmail(company.Email)
	s.CookieHandler.SetCookieAndToken(ctx, с, newUser.Id)
}

func (s *companyHandler) Logout(ctx *fasthttp.RequestCtx) {
	ctx.SetStatusCode(fasthttp.StatusOK)
	s.CookieHandler.DeleteCookie(ctx, string(ctx.Request.Header.Cookie(cnst.CookieName)))
}

func (s *companyHandler) GetCompanyById(ctx *fasthttp.RequestCtx) {
	param, _ := ctx.UserValue("id").(string)
	bytes, err := s.CompanyUseCase.GetCompanyById(param)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		log.Printf(": %s", err)
		ctx.Write([]byte("{}"))
		return
	}
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.Write(bytes)
}

func (s *companyHandler) SearchCompanies(ctx *fasthttp.RequestCtx) {
	//var request = &domain.CompanySearch{}
	var request domain.CompanySearch
	if err := json.Unmarshal(ctx.Request.Body(), &request); err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.SetBody([]byte{})
		return
	}

	bytes, err := s.CompanyUseCase.SearchCompanies(request)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.SetBody([]byte{})
		return
	}

	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.Write(bytes)
}

func (s *companyHandler) GetCompaniesByCategoryId(ctx *fasthttp.RequestCtx) {
	param, _ := ctx.UserValue("id").(string)
	bytes, err := s.CompanyUseCase.GetCompaniesByCategoryId(param)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		log.Printf("companyHandler error while getting list: %s", err)
		ctx.Write([]byte("{}"))
		return
	}
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.Write(bytes)
}
