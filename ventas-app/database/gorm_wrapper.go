package database

import (
    "gorm.io/gorm"
)

// GormDB envuelve *gorm.DB para implementar database.DBHandler
type GormDB struct{
    DB *gorm.DB
}

func (g *GormDB) Where(query interface{}, args ...interface{}) DBHandler {
    return &GormDB{DB: g.DB.Where(query, args...)}
}

func (g *GormDB) First(dest interface{}, conds ...interface{}) error {
    return g.DB.First(dest, conds...).Error
}

func (g *GormDB) Create(value interface{}) error {
    return g.DB.Create(value).Error
}

func (g *GormDB) Save(value interface{}) error {
    return g.DB.Save(value).Error
}

func (g *GormDB) Find(dest interface{}, conds ...interface{}) error {
    return g.DB.Find(dest, conds...).Error
}
