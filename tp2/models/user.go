package models

import "gorm.io/gorm"

// go get gorm.io/gorm gorm.io/driver/postgres
// User define la estructura de nuestra tabla en la base de datos
type User struct {
	gorm.Model
	Name  string
	Email string `gorm:"unique"`
}
