package model

import "gorm.io/gorm"

type Run struct {
	gorm.Model

	Input  string
	Output string

	Function
}
