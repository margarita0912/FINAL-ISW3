package models

import (
	"gorm.io/gorm"
)

type Usuario struct {
	gorm.Model
	Nombre string `json:"nombre" gorm:"unique;not null"`
	Clave  string `json:"clave" gorm:"not null"` // Hasheada con bcrypt
	Rol    string `json:"rol" gorm:"not null"`   // Validado en el controlador
}
