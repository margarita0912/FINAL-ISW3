package mocks

import (
	"errors"
	"ventas-app/database"
)

// FakeDB simula una conexión mínima que implementa database.DBHandler
// Los métodos devuelven error o el mismo mock para permitir encadenamiento.
type FakeDB struct {
	shouldFail bool
}

func (f *FakeDB) Where(query interface{}, args ...interface{}) database.DBHandler {
	return f
}

func (f *FakeDB) First(dest interface{}, conds ...interface{}) error {
	if f.shouldFail {
		return errors.New("record not found")
	}
	return nil
}

func (f *FakeDB) Create(value interface{}) error {
	if f.shouldFail {
		return errors.New("error al crear")
	}
	return nil
}

func (f *FakeDB) Save(value interface{}) error {
	if f.shouldFail {
		return errors.New("error al guardar")
	}
	return nil
}

func (f *FakeDB) Find(dest interface{}, conds ...interface{}) error {
	if f.shouldFail {
		return errors.New("error al listar")
	}
	return nil
}

// NewFakeDB devuelve un database.DBHandler simple
func NewFakeDB(shouldFail bool) database.DBHandler {
	return &FakeDB{shouldFail: shouldFail}
}
