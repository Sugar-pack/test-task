package repository

import (
	"context"
	"github.com/Sugar-pack/test-task/internal/model"
)

type CompanyForFilter struct {
	Name    string
	Code    string
	Country string
	Website string
	Phone   string
}

type Company struct {
	Name    string `db:"name"`
	Code    string `db:"code"`
	Country string `db:"country"`
	Website string `db:"website"`
	Phone   string `db:"phone"`
}

type CompanyRepository interface {
	CreateCompany(ctx context.Context, company *model.Company) error
	GetCompany(ctx context.Context, company *CompanyForFilter) (*model.Company, error)
	DeleteCompany(ctx context.Context, company *CompanyForFilter) error
	UpdateCompany(ctx context.Context, company *model.Company) error
}
