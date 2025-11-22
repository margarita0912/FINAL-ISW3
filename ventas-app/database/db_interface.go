package database

// DBHandler define una interfaz mínima que abstrae GORM para facilitar tests.
type DBHandler interface {
	// Where devuelve otra DBHandler para permitir encadenamiento en controllers.
	Where(query interface{}, args ...interface{}) DBHandler
	// First/Find/Create/Save ejecutan la operación y retornan un error si falla.
	First(dest interface{}, conds ...interface{}) error
	Create(value interface{}) error
	Save(value interface{}) error
	Find(dest interface{}, conds ...interface{}) error
}
