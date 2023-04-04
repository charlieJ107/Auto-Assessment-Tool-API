package models

import "gorm.io/gorm"

type Tag struct {
	gorm.Model
	Name string `gorm:"index:idx_name; unique; not null"`
}
