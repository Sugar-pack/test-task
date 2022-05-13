package model

type Company struct {
	Name    string `json:"name,omitempty"`
	Code    string `json:"code"`
	Country string `json:"country"`
	Website string `json:"website"`
	Phone   string `json:"phone"`
}

type CompanyForUpdate struct {
	FilterFields    Company          `json:"filter_fields"`
	FieldsForUpdate CompanyUpdatable `json:"fields_for_update"`
}

type CompanyUpdatable struct {
	Code    string `json:"code"`
	Country string `json:"country"`
	Website string `json:"website"`
	Phone   string `json:"phone"`
}
