# VentasTP7DeMarcos-strumia

README – TP7: Testing & CI/CD
# Proyecto

Sistema de Ventas – Ingeniería de Software III

Aplicación fullstack compuesta por:

Backend: Go (Gin + GORM)

Frontend: React (Vite + TypeScript)

Pipeline CI/CD: GitLab CI

Testing: Unitarios, integración, E2E (Cypress), y análisis estático con SonarCloud

# Estructura del pipeline

El pipeline se compone de 6 etapas en orden secuencial:

Etapa	Descripción	Herramienta
- build_backend	Compila el backend de Go y valida dependencias	Go 1.24
- build_frontend	Instala dependencias y construye el bundle de React	Node 20
- test_backend	Ejecuta tests de Go con go test y genera cobertura XML	go test + gocover-cobertura
- test_frontend	Ejecuta tests de Jest con reporte LCOV	Jest + React Testing Library
- sonarcloud_analysis	Análisis estático del código y métricas de calidad	SonarCloud
- e2e_tests	Pruebas de flujo completo sobre el frontend	Cypress
- Testing
- Backend (Go)

Se utilizan tests unitarios en el paquete controllers

Cobertura actual: ≈ 94%

Comando local:

go test ./controllers -v -coverprofile=coverage.out
go tool cover -html=coverage.out

Ver coverage.html o capturas en carpeta de evidencias


💻 Frontend (React)

Se usa Jest con --coverage

Cobertura: ≈ 60–70%

Archivos .spec.tsx en src/__tests__/

Ver captura en evidencia

🌐 E2E (Cypress)

Pruebas completas de flujo de ventas: crear, actualizar y validar errores.

Script ejecutado en CI con:

npx cypress run --browser chrome --headless


Ver captura en evidencia

☁️ SonarCloud

El análisis verifica:

Duplicaciones

Vulnerabilidades

Smells

Cobertura global combinada Go + React


🧩 Pipeline completo


📚 Tecnologías clave
Componente	Tecnología
Backend	Go 1.24, Gin, GORM
Frontend	React, Vite, TypeScript
CI/CD	GitLab CI
Testing	Go test, Jest, Cypress
QA	SonarCloud, Cobertura XML/LCOV