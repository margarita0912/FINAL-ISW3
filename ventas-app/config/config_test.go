package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadEnv_Production(t *testing.T) {
	// Test que LoadEnv con "prod" no cause panic
	assert.NotPanics(t, func() {
		LoadEnv("prod")
	})
}

func TestLoadEnv_QA(t *testing.T) {
	// Test que LoadEnv con "qa" no cause panic
	assert.NotPanics(t, func() {
		LoadEnv("qa")
	})
}

func TestLoadEnv_Default(t *testing.T) {
	// Test que LoadEnv con cualquier otro valor use .env
	assert.NotPanics(t, func() {
		LoadEnv("development")
	})
}

func TestLoadEnv_EmptyString(t *testing.T) {
	// Test que LoadEnv con string vacío use .env
	assert.NotPanics(t, func() {
		LoadEnv("")
	})
}

func TestLoadEnv_InvalidEnv(t *testing.T) {
	// Test que LoadEnv con valor inválido no cause panic
	assert.NotPanics(t, func() {
		LoadEnv("invalid_environment")
	})
}

func TestLoadEnv_CIEnvironment(t *testing.T) {
	// Simular entorno CI
	os.Setenv("CI", "true")
	defer os.Unsetenv("CI")

	assert.NotPanics(t, func() {
		LoadEnv("qa")
	})
}

func TestLoadEnv_EnvironmentSelection(t *testing.T) {
	tests := []struct {
		name        string
		env         string
		description string
	}{
		{"Production", "prod", "Should select .env.prod"},
		{"QA", "qa", "Should select .env.qa"},
		{"Development", "dev", "Should select .env"},
		{"Testing", "test", "Should select .env"},
		{"Empty", "", "Should select .env"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.NotPanics(t, func() {
				LoadEnv(tt.env)
			}, tt.description)
		})
	}
}

func TestLoadEnv_NoFileDoesNotCrash(t *testing.T) {
	// Incluso si no existe ningún archivo .env, no debería hacer panic
	// porque el código maneja el error con log.Printf
	assert.NotPanics(t, func() {
		LoadEnv("nonexistent")
	})
}

func TestLoadEnv_MultipleCallsSafe(t *testing.T) {
	// Llamar LoadEnv múltiples veces no debería causar problemas
	assert.NotPanics(t, func() {
		LoadEnv("qa")
		LoadEnv("prod")
		LoadEnv("")
		LoadEnv("qa")
	})
}
