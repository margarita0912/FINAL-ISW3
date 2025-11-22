package models

import "gorm.io/gorm"

type Venta struct {
	gorm.Model
	UsuarioID   uint    `json:"usuario_id"`
	ProductoID  uint    `json:"producto_id"`
	Cantidad    int     `json:"cantidad"`
	Descuento   float64 `json:"descuento"`
	PrecioFinal float64 `json:"precio_final"`
}
