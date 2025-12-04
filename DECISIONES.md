# Decisiones Técnicas - TP Final IS3

## 🛠️ Stack Tecnológico

### **Backend**
- **Lenguaje:** Go 1.22
- **Framework Web:** Gin (router HTTP de alto rendimiento)
- **ORM:** GORM (para interacción con base de datos)
- **Autenticación:** JWT con `golang-jwt/jwt/v5` + bcrypt
- **CORS:** gin-contrib/cors
- **Driver DB:** go-sql-driver/mysql

### **Frontend**
- **Framework:** React 19.1.1
- **Lenguaje:** TypeScript
- **Build Tool:** Vite 7.1.7
- **Routing:** React Router 7.1.1
- **HTTP Client:** Axios
- **Styling:** Tailwind CSS (inline classes)

### **Base de Datos**
- **Motor:** MySQL 8
- **Hosting:** Railway
- **Conexión:** TLS habilitado

### **Testing**
- **Backend:** Go testing + testify
- **Frontend:** Jest 29.7.0 + React Testing Library + MSW
- **E2E:** Cypress 15.6.0

### **CI/CD**
- **Pipeline:** GitHub Actions
- **Registry:** GitHub Container Registry (ghcr.io)
- **Quality Gate:** SonarCloud
- **Deployment:** Render.com

### **Containerización**
- **Docker:** Multi-stage builds
- **Base Images:** 
  - Backend: `golang:1.24` → `node:20-slim`
  - Frontend: `node:20` (build) → incluido en imagen unificada

---

## 📋 Testing Strategy

### 1. **Unit Tests - Backend (Go)**

**Ubicación:** `ventas-app/**/*_test.go`

**Framework:** Go testing + testify

**Archivos:**
- `config/config_test.go` - Tests de carga de configuración
- `controllers/*_controller_test.go` - Tests de controladores con mocks
- `middleware/auth_test.go` - Tests de autenticación JWT
- `models/*_test.go` - Tests de modelos de datos
- `routes/routes_test.go` - Tests de registro de rutas
- `utils/jwt_test.go` - Tests de generación y validación de tokens

**Cobertura:** 96.1% en controllers, 100% en middleware/routes/utils

**Cómo ejecutar:**
```bash
cd ventas-app
go test ./... -v -coverprofile=coverage.out
```

**Características:**
- Usan mocks para aislar la base de datos
- Validan lógica de negocio sin dependencias externas
- Rápidos de ejecutar (~1 segundo)

**Integración con SonarCloud:**
- Genera `coverage.out` en formato Go coverage
- SonarCloud lee este archivo vía propiedad `sonar.go.coverage.reportPaths=ventas-app/coverage.out`
- Reporta coverage por paquete y líneas cubiertas/no cubiertas

---

### 2. **Unit Tests - Frontend (React)**

**Ubicación:** `ventas-frontend/src/tests/*.test.tsx`

**Framework:** Jest + React Testing Library + MSW (Mock Service Worker)

**Archivos:**
- `App.test.tsx` - Tests del componente principal
- `Login.test.tsx` - Tests de autenticación
- `Ventas.test.tsx` - Tests de registro de ventas
- `CrearProductos.test.tsx` - Tests de creación de productos
- `CrearUsuario.test.tsx` - Tests de creación de usuarios
- `RoleRoute.test.tsx` - Tests de protección de rutas por rol
- `EntornoSelector.test.tsx` - Tests de selector de entorno
- `useValidacion.test.ts` - Tests de hook de validación

**Cobertura:** 67.77% (64 tests)

**Cómo ejecutar:**
```bash
cd ventas-frontend
npm run test:ci
```

**Características:**
- Mockean llamadas HTTP con MSW
- Validan rendering, interacciones de usuario, y estados
- Ejecutan en ~15 segundos

**Integración con SonarCloud:**
- Jest configurado con `coverageReporters: ['text', 'lcov', 'clover', 'json']`
- Genera `lcov.info` en formato LCOV estándar
- SonarCloud lee vía `sonar.javascript.lcov.reportPaths=ventas-frontend/coverage/lcov.info`
- Incluye cobertura de líneas, branches, y funciones

---

### 3. **Integration Tests E2E (Cypress)**

**Ubicación:** `ventas-frontend/cypress/e2e/`

**Framework:** Cypress

**Archivos:**
- `acceptance.cy.js` - Tests de aceptación básicos
- `ventas_flow.cy.js` - Flujo completo de ventas

**Cómo ejecutar:**
```bash
cd ventas-frontend
npx cypress run --spec "cypress/e2e/ventas_flow.cy.js,cypress/e2e/acceptance.cy.js"
```

**Características:**
- Se ejecutan contra el ambiente de QA desplegado
- Validan el sistema completo (frontend + backend + base de datos)
- Ejecutan después del deploy a QA en el pipeline
- Credenciales de test: `julio/julio123`

**Casos de prueba:**
- Login con credenciales válidas/inválidas
- Navegación entre páginas según rol
- Flujo completo de registro de venta
- Validaciones de formularios

---

## 🔄 Pipeline CI/CD

**Archivo:** `.github/workflows/main.yml`

### **Flujo General:**

```
Push a main/qa/prod
    ↓
[Job 1] Unit Tests & SonarCloud Analysis
    ├─ Backend: Tests + Build
    ├─ Frontend: Tests + Build
    └─ SonarCloud: Análisis de cobertura y calidad
    ↓
[Job 2] Build & Push Docker Images
    ├─ Build imagen unificada (backend + frontend)
    ├─ Tag: {SHA}, {branch}-latest
    └─ Push a ghcr.io
    ↓
[Job 3] Deploy to QA (solo si branch = main/qa)
    ├─ Deploy a Render QA (webhook)
    ├─ Wait 60s para estabilización
    └─ Cypress E2E tests contra QA
    ↓
[Job 4] Deploy to Production (solo si branch = main/prod)
    ├─ ⏸️ Requiere aprobación manual (environment: production)
    ├─ Tag imagen como prod-release
    └─ Deploy a Render PROD (webhook)
```

### **Jobs Detallados:**

#### **Job 1: Unit Tests & SonarCloud Analysis**
- **Duración:** ~2 minutos
- **Acciones:**
  1. Checkout código
  2. Setup Go 1.22 → tests backend → build binario
  3. Setup Node 20 → tests frontend → build dist
  4. SonarCloud scan con reportes de cobertura
- **Outputs:** Coverage reports (lcov.info, coverage.out)

#### **Job 2: Build & Push Docker Images**
- **Duración:** ~1-2 minutos
- **Acciones:**
  1. Login a GitHub Container Registry (ghcr.io)
  2. Build imagen multi-stage (Dockerfile)
  3. Push con múltiples tags:
     - `ghcr.io/margarita0912/final-isw3:{SHA}` (commit específico, ej: `abc1234`)
     - `ghcr.io/margarita0912/final-isw3:main-latest` (última versión de branch main)

#### **Job 3: Deploy to QA**
- **Duración:** ~3-5 minutos
- **Condición:** `if: github.ref == 'refs/heads/qa' || github.ref == 'refs/heads/main'`
- **Acciones:**
  1. Trigger deploy webhook de Render QA
  2. Wait 60s para que Render complete el deploy
  3. Wait-on hasta que QA responda (timeout 120s)
  4. Ejecutar Cypress E2E tests contra QA_URL

#### **Job 4: Deploy to Production**
- **Duración:** ~5 segundos + tiempo de aprobación manual
- **Condición:** Requiere que QA haya pasado exitosamente
- **Acciones:**
  1. **⏸️ Espera aprobación manual** (GitHub environment: production)
  2. Trigger deploy webhook de Render PROD

---

## 🏷️ Estrategia de Tags de Imágenes Docker

### **Dos Tags - Mismo SHA**

Cada build genera **2 tags apuntando a la misma imagen**:

1. **`main-latest`** - Para QA (siempre actualizado)
2. **`prod-release`** - Para PROD (misma imagen, nombre diferente)

### **Tags por Ambiente:**

| Ambiente | Tag | Actualización | Uso |
|----------|-----|---------------|-----|
| **QA** | `main-latest` | Cada push a main | Deploy automático |
| **PROD** | `prod-release` | Cada push a main | Deploy tras aprobación manual |

### **Flujo de Despliegue:**

```
Push a main (commit abc123)
  ↓
Build crea 2 tags de la MISMA imagen:
  - ghcr.io/.../final-isw3:main-latest
  - ghcr.io/.../final-isw3:prod-release
  (ambos apuntan al mismo SHA de imagen)
  ↓
QA despliega: main-latest (abc123)
  ↓
Tests E2E pasan ✓
  ↓
Aprobación manual en GitHub
  ↓
PROD despliega: prod-release (abc123, misma imagen que QA)
```

**Ventajas:**
- ✅ Tags separados por ambiente (claridad)
- ✅ Ambos tags siempre sincronizados (mismo SHA)
- ✅ Render configuración diferenciada pero misma imagen
- ✅ Simple: Solo 2 tags, sin re-taggeo manual

**Consideración:**
- ℹ️ Ambos tags se actualizan con cada push (apuntan a la misma imagen nueva)
- ✅ **Protección:** Concurrency queue evita que pipelines se ejecuten en paralelo

---

## 🔐 Solución al Problema de Race Conditions

### **Problema Original:**
Si dos pipelines corren en paralelo:
- Pipeline 1 despliega a QA → corre tests
- Pipeline 2 sobrescribe la imagen `main-latest` mientras Pipeline 1 testea
- Pipeline 1 aprueba a PROD → despliega imagen incorrecta (del Pipeline 2)

### **Solución Implementada: Concurrency Queue**

```yaml
concurrency:
  group: deploy-${{ github.ref }}
  cancel-in-progress: false  # No cancelar, hacer cola secuencial
```

**Cómo funciona:**
- Los pipelines se ejecutan **uno por uno en cola**, no en paralelo
- Si llega un push mientras otro está corriendo, el nuevo **espera en cola**
- Garantiza que cada pipeline completa QA → Tests → Aprobación antes del siguiente

**Ejemplo de ejecución:**
```
Timeline:
─────────────────────────────────────────────────
09:00 - Push commit A
09:01 - Pipeline 1 inicia: Build → Deploy QA (main-latest = A)
09:05 - Push commit B
09:05 - Pipeline 2 ESPERA (en cola, no ejecuta nada)
09:06 - Pipeline 1: Tests E2E en QA (con imagen A) ✓
09:08 - Pipeline 1: Aprobación manual ✓
09:09 - Pipeline 1: Deploy PROD (main-latest = A) ✓
09:10 - Pipeline 2 COMIENZA: Build → Deploy QA (main-latest = B)
        QA ahora tiene B, PROD sigue con A hasta nueva aprobación
```

**Ventajas:**
- ✅ Evita race conditions mediante cola secuencial
- ✅ QA siempre testea la imagen correcta
- ✅ PROD despliega lo que fue aprobado (mientras no haya nuevo push)
- ✅ Muy simple de implementar (una sola línea de config)
- ✅ No requiere tags adicionales ni infraestructura extra

**Consideración:**
- ⚠️ Si alguien **cancela manualmente** un pipeline en GitHub Actions, la protección se rompe
- ⚠️ Si hay un push nuevo **después de aprobar pero antes de deploy**, PROD podría tomar la imagen nueva
- ✅ En la práctica esto es raro porque el deploy a PROD es inmediato tras aprobación

---

## 🗄️ Base de Datos

**Proveedor:** Railway MySQL

**Ambientes:**
- **QA:** `witchyard.proxy.rlwy.net:20665` (db: `railway`)
- **PROD:** `amanote.proxy.rlwy.net:50180` (db: `railway`)

**Configuración:** Variables de entorno en Render

---

## 🚀 Deployment

**Plataforma:** Render.com

**Servicios:**
- **QA:** `https://tp8-front-qa.onrender.com`
- **PROD:** `https://tp8-front-prod.onrender.com`

**Estrategia:**
- Imagen Docker unificada (backend + frontend en un solo contenedor)
- Backend (Go) sirve el frontend estático y expone API en `/api/*`
- Auto-deploy habilitado al detectar nuevo tag en ghcr.io

---

## 📊 Code Quality

**Herramienta:** SonarCloud

**Métricas actuales:**
- Backend: ~96% cobertura
- Frontend: ~68% cobertura
- Quality Gate: Configurado para requerir >60% coverage

**Integración:** Job en pipeline analiza código después de tests unitarios
