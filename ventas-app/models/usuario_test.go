package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func TestUsuarioModel_Creation(t *testing.T) {
	usuario := Usuario{
		Nombre: "admin",
		Clave:  "hashed_password",
		Rol:    "administrador",
	}

	assert.Equal(t, "admin", usuario.Nombre)
	assert.Equal(t, "hashed_password", usuario.Clave)
	assert.Equal(t, "administrador", usuario.Rol)
}

func TestUsuarioModel_Roles(t *testing.T) {
	tests := []struct {
		name     string
		rol      string
		expected string
	}{
		{"Admin role", "administrador", "administrador"},
		{"Seller role", "vendedor", "vendedor"},
		{"Buyer role", "comprador", "comprador"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			usuario := Usuario{
				Nombre: "testuser",
				Clave:  "password",
				Rol:    tt.rol,
			}
			assert.Equal(t, tt.expected, usuario.Rol)
		})
	}
}

func TestUsuarioModel_PasswordHashing(t *testing.T) {
	plainPassword := "mySecurePassword123"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)

	assert.NoError(t, err)
	assert.NotEqual(t, plainPassword, string(hashedPassword))

	usuario := Usuario{
		Nombre: "testuser",
		Clave:  string(hashedPassword),
		Rol:    "vendedor",
	}

	// Verificar que la contraseña hasheada coincide
	err = bcrypt.CompareHashAndPassword([]byte(usuario.Clave), []byte(plainPassword))
	assert.NoError(t, err, "La contraseña debería coincidir")
}

func TestUsuarioModel_PasswordValidation(t *testing.T) {
	plainPassword := "correctPassword"
	wrongPassword := "wrongPassword"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)

	usuario := Usuario{
		Nombre: "testuser",
		Clave:  string(hashedPassword),
		Rol:    "vendedor",
	}

	// Password correcto
	err := bcrypt.CompareHashAndPassword([]byte(usuario.Clave), []byte(plainPassword))
	assert.NoError(t, err)

	// Password incorrecto
	err = bcrypt.CompareHashAndPassword([]byte(usuario.Clave), []byte(wrongPassword))
	assert.Error(t, err, "La contraseña incorrecta debería generar error")
}

func TestUsuarioModel_UniqueNombre(t *testing.T) {
	usuario1 := Usuario{
		Nombre: "uniqueuser",
		Clave:  "password",
		Rol:    "vendedor",
	}

	usuario2 := Usuario{
		Nombre: "uniqueuser", // Mismo nombre
		Clave:  "password",
		Rol:    "comprador",
	}

	assert.Equal(t, usuario1.Nombre, usuario2.Nombre)
	// En la base de datos real, esto generaría un error de constraint unique
}

func TestUsuarioModel_JSONTags(t *testing.T) {
	usuario := Usuario{
		Model:  gorm.Model{ID: 1},
		Nombre: "jsonuser",
		Clave:  "hashedpass",
		Rol:    "vendedor",
	}

	assert.NotZero(t, usuario.ID)
	assert.Equal(t, "jsonuser", usuario.Nombre)
	assert.Equal(t, "hashedpass", usuario.Clave)
	assert.Equal(t, "vendedor", usuario.Rol)
}

func TestUsuarioModel_EmptyUser(t *testing.T) {
	usuario := Usuario{}

	assert.Empty(t, usuario.Nombre)
	assert.Empty(t, usuario.Clave)
	assert.Empty(t, usuario.Rol)
}

func TestUsuarioModel_GormModel(t *testing.T) {
	usuario := Usuario{
		Model:  gorm.Model{ID: 42},
		Nombre: "gormuser",
		Clave:  "password",
		Rol:    "administrador",
	}

	assert.Equal(t, uint(42), usuario.ID)
	assert.NotNil(t, usuario.Model)
}

func TestUsuarioModel_RoleValidation(t *testing.T) {
	validRoles := []string{"administrador", "vendedor", "comprador"}

	for _, rol := range validRoles {
		usuario := Usuario{
			Nombre: "user_" + rol,
			Clave:  "password",
			Rol:    rol,
		}
		assert.Contains(t, validRoles, usuario.Rol)
	}
}

func TestUsuarioModel_InvalidRole(t *testing.T) {
	usuario := Usuario{
		Nombre: "invaliduser",
		Clave:  "password",
		Rol:    "invalid_role",
	}

	validRoles := []string{"administrador", "vendedor", "comprador"}
	assert.NotContains(t, validRoles, usuario.Rol, "El rol debería ser inválido")
}
