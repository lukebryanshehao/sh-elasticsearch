package models

type Company struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt string
	UpdatedAt string
	DeletedAt string `sql:"index"`
	Name string
	Abbreviation string
	CompanyTypeId uint
	Sort int
	State int
	CompanyNum string
}