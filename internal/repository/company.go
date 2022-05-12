package repository

import (
	"context"
)

type CompanyForFilter struct {
	Name    string `db:"name"`
	Code    string `db:"code"`
	Country string `db:"country"`
	Website string `db:"website"`
	Phone   string `db:"phone"`
}

type Company struct {
	Name    string `db:"name"`
	Code    string `db:"code"`
	Country string `db:"country"`
	Website string `db:"website"`
	Phone   string `db:"phone"`
}

type CompanyForUpdate struct {
	FilterFields    CompanyForFilter
	FieldsForUpdate Company
}

type CompanyRepository interface {
	CreateCompany(ctx context.Context, company *Company) error
	GetCompany(ctx context.Context, company *CompanyForFilter) ([]Company, error)
	DeleteCompany(ctx context.Context, company *CompanyForFilter) (int64, error)
	UpdateCompany(ctx context.Context, company *CompanyForUpdate) (int64, error)
}
