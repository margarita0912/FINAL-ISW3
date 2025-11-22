package models

import "gorm.io/gorm"

type Compra struct {
	gorm.Model
	UsuarioID  uint    `json:"usuario_id"`
	ProductoID uint    `json:"producto_id"`
	Cantidad   int     `json:"cantidad"`
	CostoUnit  float64 `json:"costo_unit"`
}
