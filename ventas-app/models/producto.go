package models

import "gorm.io/gorm"

type Producto struct {
	gorm.Model
	Nombre string  `json:"nombre"`
	Costo  float64 `json:"costo"`
	Precio float64 `json:"precio"`
	Stock  int     `json:"stock"`
}
