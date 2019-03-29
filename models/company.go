package models

type Company struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt string
	UpdatedAt string
	DeletedAt string `sql:"index"`
	Name string
	Abbreviation string
	CompanyTypeId int
	Sort int
	State int
	CompanyNum string
}