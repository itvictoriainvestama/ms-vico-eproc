package models

import "gorm.io/gorm"

type BaseModel struct {
	gorm.Model
	CreatedBy *uint `gorm:"index" json:"created_by,omitempty"`
	UpdatedBy *uint `gorm:"index" json:"updated_by,omitempty"`
}
