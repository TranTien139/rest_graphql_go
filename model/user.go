package model

import (
	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	ID     uint   `gorm:"primary_key" json "id"`
	Title    string `gorm:"not null" json:"title"`
	Description uint   `gorm:"not null" json:"description"`
}