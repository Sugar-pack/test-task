package repository

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type SQLCompanyRepo struct {
	dbConn sqlx.DB
}

func CreatePattern(param string) string {
	return fmt.Sprintf("%%%s%%", param)
}

func (r *SQLCompanyRepo) GetCompany(ctx context.Context, company *CompanyForFilter) ([]Company, error) {
	query := `SELECT * FROM companies WHERE name LIKE $1 AND code LIKE $2 AND country LIKE $3 
                          AND website LIKE $4 AND phone LIKE $5`
	var companies []Company
	err := r.dbConn.SelectContext(ctx, &companies, query, CreatePattern(company.Name), CreatePattern(company.Code),
		CreatePattern(company.Country), CreatePattern(company.Website), CreatePattern(company.Phone))
	if err != nil {
		return nil, err
	}

	return companies, nil
}

func (r *SQLCompanyRepo) DeleteCompany(ctx context.Context, company *CompanyForFilter) (int64, error) {
	query := `DELETE FROM companies WHERE name 
                                LIKE $1 AND code LIKE $2 AND country LIKE $3 AND website LIKE $4 AND phone LIKE $5`
	result, err := r.dbConn.ExecContext(ctx, query, CreatePattern(company.Name), CreatePattern(company.Code),
		CreatePattern(company.Country), CreatePattern(company.Website), CreatePattern(company.Phone))
	if err != nil {
		return 0, err
	}

	affectedRow, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("cant get affected rows: %w", err)
	}

	return affectedRow, nil
}

func (r *SQLCompanyRepo) UpdateCompany(ctx context.Context, company *CompanyForUpdate) (int64, error) {
	query := `UPDATE companies SET name = $1, code = $2, country = $3, website = $4, phone = $5 WHERE name
            LIKE $6 AND code LIKE $7 AND country LIKE $8 AND website LIKE $9 AND phone LIKE $10`
	result, err := r.dbConn.ExecContext(ctx, query, company.FieldsForUpdate.Name,
		company.FieldsForUpdate.Code, company.FieldsForUpdate.Country, company.FieldsForUpdate.Website,
		company.FieldsForUpdate.Phone, CreatePattern(company.FilterFields.Name), CreatePattern(company.FilterFields.Code),
		CreatePattern(company.FilterFields.Country), CreatePattern(company.FilterFields.Website),
		CreatePattern(company.FilterFields.Phone))
	if err != nil {
		return 0, err
	}

	affectedRow, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("cant get affected rows: %w", err)
	}

	return affectedRow, nil
}

func (r *SQLCompanyRepo) CreateCompany(ctx context.Context, company *Company) error {
	query := `INSERT INTO companies (name, code, country, website, phone) 
				VALUES (:name, :code, :country, :website, :phone)`
	_, err := r.dbConn.NamedExecContext(ctx, query, company)

	return err
}

func NewPsqlRepository(dbConn *sqlx.DB) *SQLCompanyRepo {
	return &SQLCompanyRepo{*dbConn}
}
