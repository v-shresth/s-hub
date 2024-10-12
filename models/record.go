package models

import "gorm.io/gorm"

type TypeName map[string]interface{}
type Record struct {
	gorm.Model
	TypeName
}
