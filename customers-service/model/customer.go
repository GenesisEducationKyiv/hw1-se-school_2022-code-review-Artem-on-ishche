package model

import "gorm.io/gorm"

type Customer struct {
	gorm.Model
	IDTransaction string
	EmailAddress  string
	Status        string
}
