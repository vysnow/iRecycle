package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string `gorm:"type:varchar(20);not null;unique"`
	Mobile   string `gorm:"varchar(11);not null;unique"`
	Password string `gorm:"size(128);not null`
}
