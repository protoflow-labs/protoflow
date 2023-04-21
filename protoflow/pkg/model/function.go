package model

import "gorm.io/gorm"

type Function struct {
	gorm.Model
	Name       string
	ImportPath string
	Entrypoint string
}
