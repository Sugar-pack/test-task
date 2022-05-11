package repository

import (
	"context"
	"github.com/Sugar-pack/test-task/internal/model"
	"github.com/jmoiron/sqlx"
)

type SQLCompanyRepo struct {
	dbConn sqlx.DB
}

func (r *SQLCompanyRepo) CreateCompany(ctx context.Context, company *model.Company) error {
	query := `INSERT INTO companies (name, code, country, website, phone) VALUES (:name, :code, :country, :website, :phone)`
	_, err := r.dbConn.NamedExecContext(ctx, query, company)
	return err
}

func (r *SQLCompanyRepo) GetCompany(ctx context.Context, company *CompanyForFilter) (*model.Company, error) {
	//TODO implement me
	panic("implement me")
}

func (r *SQLCompanyRepo) DeleteCompany(ctx context.Context, company *CompanyForFilter) error {
	//TODO implement me
	panic("implement me")
}

func (r *SQLCompanyRepo) UpdateCompany(ctx context.Context, company *model.Company) error {
	//TODO implement me
	panic("implement me")
}

func NewPsqlRepository(dbConn *sqlx.DB) *SQLCompanyRepo {
	return &SQLCompanyRepo{*dbConn}
}
