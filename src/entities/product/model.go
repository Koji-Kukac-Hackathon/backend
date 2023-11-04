package product

import "gorm.io/gorm"

type Model struct {
	gorm.Model

	Code  string
	Price uint
}
