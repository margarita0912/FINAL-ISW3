package mocks

import (
	"errors"
	"ventas-app/database"
	"ventas-app/models"

	"gorm.io/gorm"
)

type MockDB struct {
	Usuarios  []models.Usuario
	Productos []models.Producto
	ShouldErr bool
	// Flags más finos para simular errores en operaciones concretas
	FailCreate bool
	FailSave   bool
	FailFind   bool
	FailFirst  bool
	// Registros creados
	Compras []models.Compra
	Ventas  []models.Venta
}

func (m *MockDB) Where(query interface{}, args ...interface{}) database.DBHandler {
	// Mantener compatibilidad: devolveremos el mismo mock para permitir encadenamiento
	// Guardar el filtro para que First pueda utilizarlo
	// query suele ser una string como "nombre = ?" y args[0] el valor buscado
	return m
}

// Implementación simplificada: los métodos devuelven error en lugar de *gorm.DB
func (m *MockDB) First(dest interface{}, conds ...interface{}) error {
	if m.ShouldErr || m.FailFirst {
		return gorm.ErrRecordNotFound
	}

	// Si se pasa un ID como conds[0], buscar en Productos/Usuarios por ID
	if len(conds) > 0 {
		switch d := dest.(type) {
		case *models.Producto:
			// conds[0] puede ser int, int64, uint, etc.
			var id uint
			switch v := conds[0].(type) {
			case int:
				id = uint(v)
			case int64:
				id = uint(v)
			case uint:
				id = v
			case uint64:
				id = uint(v)
			}
			for _, p := range m.Productos {
				if p.ID == id {
					*d = p
					return nil
				}
			}
			return gorm.ErrRecordNotFound
		}
	}

	// Si dest es *models.Usuario y no se pasó ID, intentar devolver el primer usuario
	switch d := dest.(type) {
	case *models.Usuario:
		if len(m.Usuarios) > 0 {
			*d = m.Usuarios[0]
			return nil
		}
		return gorm.ErrRecordNotFound
	case *models.Producto:
		if len(m.Productos) > 0 {
			*d = m.Productos[0]
			return nil
		}
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (m *MockDB) Create(value interface{}) error {
	if m.ShouldErr || m.FailCreate {
		return errors.New("error al crear")
	}

	// Simular creación: si es Compra o Venta, añadir a slice correspondiente
	switch v := value.(type) {
	case *models.Compra:
		m.Compras = append(m.Compras, *v)
		return nil
	case *models.Venta:
		// asignar ID secuencial simple
		nextID := uint(len(m.Ventas) + 1)
		v.ID = nextID
		m.Ventas = append(m.Ventas, *v)
		return nil
	case *models.Producto:
		// simular asignación de ID si es cero
		if v.ID == 0 {
			v.ID = uint(len(m.Productos) + 1)
		}
		m.Productos = append(m.Productos, *v)
		return nil
	}

	return nil
}

func (m *MockDB) Save(value interface{}) error {
	if m.ShouldErr || m.FailSave {
		return errors.New("error al guardar")
	}

	// Simular update: si es Producto, buscar por ID y actualizar
	switch v := value.(type) {
	case *models.Producto:
		for i, p := range m.Productos {
			if p.ID == v.ID {
				m.Productos[i] = *v
				return nil
			}
		}
		// si no existe, agregar
		m.Productos = append(m.Productos, *v)
		return nil
	}

	return nil
}

func (m *MockDB) Find(dest interface{}, conds ...interface{}) error {
	if m.ShouldErr || m.FailFind {
		return errors.New("error al listar")
	}

	// Rellenar slices con los datos almacenados en el mock
	switch d := dest.(type) {
	case *[]models.Producto:
		*d = append((*d)[:0], m.Productos...)
		return nil
	case *[]models.Usuario:
		*d = append((*d)[:0], m.Usuarios...)
		return nil
	}

	return nil
}
