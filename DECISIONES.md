# Decisiones Técnicas - TP Final IS3

## 🎯 Solución al Problema de Race Conditions

### **Problema Identificado**
Si dos pipelines corren en paralelo y ambos modifican la misma imagen (ej: `main-latest`):
- Pipeline 1 llega a QA y comienza tests
- Pipeline 2 sobrescribe `main-latest` con código nuevo
- Pipeline 1 termina tests y despliega a PROD
- **PROD recibe código del Pipeline 2 (no testeado)** ❌

### **Solución Implementada: Concurrency Queue + SHA Tags**

**Combinación de dos estrategias:**

#### 1. **Concurrency Queue** - Ejecución Secuencial
```yaml
concurrency:
  group: deploy-${{ github.ref }}
  cancel-in-progress: false
```
- Los pipelines se ejecutan **uno por uno**, no en paralelo
- Si llega un nuevo push, **espera en cola** hasta que termine el anterior
- Garantiza que cada pipeline completa: Build → QA → Tests → PROD

#### 2. **SHA Tags** - Trazabilidad Perfecta
- Cada commit genera su propia imagen con tag único: `{SHA-corto}` (ej: `abc1234`)
- Tag adicional `latest` como alias al más reciente
- En los logs se muestra claramente qué SHA se desplegó en cada ambiente

**Ejemplo de ejecución:**
```
09:00 - Push commit abc1234
09:01 - Pipeline 1 INICIA: Build imagen abc1234
09:05 - Push commit def5678
09:05 - Pipeline 2 ESPERA en cola (no ejecuta nada)
09:06 - Pipeline 1: Deploy QA con abc1234
09:07 - Pipeline 1: Tests E2E pasan ✓
09:09 - Pipeline 1: Aprobación manual ✓
09:10 - Pipeline 1: Deploy PROD con abc1234 ✓
09:11 - Pipeline 2 COMIENZA: Build imagen def5678
```

**Ventajas:**
- ✅ **Sin race conditions**: Cola secuencial previene sobrescrituras
- ✅ **Trazabilidad**: Cada imagen identificada por SHA del commit
- ✅ **Simple**: Solo una línea de config + output del SHA
- ✅ **Compatible con Render**: Usa webhooks estándar (no requiere API)
- ✅ **Inmutabilidad**: Cada SHA nunca cambia

**Limitaciones conocidas:**
- ⏱️ Si hay muchos pushes consecutivos, se forma cola (espera ~5-10 min por pipeline)
- ⚠️ Cancelar manualmente un pipeline en GitHub rompe la protección
- ℹ️ Suficiente para este proyecto (pushes no son tan frecuentes)

---

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
  1. Extrae SHA corto del commit: `$(echo ${{ github.sha }} | cut -c1-7)`
  2. Login a GitHub Container Registry (ghcr.io)
  3. Build imagen multi-stage (Dockerfile)
  4. Push con tags:
     - `ghcr.io/margarita0912/final-isw3:{SHA}` (único por commit, ej: `abc1234`)
     - `ghcr.io/margarita0912/final-isw3:latest` (alias al más reciente)
- **Outputs:** `image_sha` - SHA corto para trazabilidad en deploys

#### **Job 3: Deploy to QA**
- **Duración:** ~3-5 minutos
- **Condición:** `if: github.ref == 'refs/heads/qa' || github.ref == 'refs/heads/main'`
- **Acciones:**
  1. Recibe `image_sha` del Job 2 (ej: `abc1234`)
  2. Log: "🚀 Deploying to QA with image SHA: abc1234"
  3. Trigger deploy webhook de Render QA
  4. Wait 60s para que Render complete el deploy
  5. Wait-on hasta que QA responda (timeout 120s)
  6. Ejecutar Cypress E2E tests contra QA_URL

#### **Job 4: Deploy to Production**
- **Duración:** ~5 segundos + tiempo de aprobación manual
- **Condición:** Requiere que QA haya pasado exitosamente
- **Acciones:**
  1. **⏸️ Espera aprobación manual** (GitHub environment: production)
  2. Recibe `image_sha` del Job 2 (mismo SHA que QA testeó)
  3. Log: "✅ Production deployment approved for SHA: abc1234"
  4. Trigger deploy webhook de Render PROD

---

## 🏷️ Estrategia de Tags de Imágenes Docker

### **SHA Tags + Latest**

Cada build genera **2 tags**:

1. **`{SHA-corto}`** - Identificador único inmutable (ej: `abc1234`)
2. **`latest`** - Alias al build más reciente

### **Ventajas del SHA Tag:**

| Característica | Beneficio |
|----------------|-----------|
| **Inmutabilidad** | Cada SHA nunca cambia, siempre apunta al mismo código |
| **Trazabilidad** | Logs muestran exactamente qué SHA se desplegó |
| **Auditoría** | Puedes ver en GitHub Actions qué commit está en cada ambiente |
| **Rollback** | Fácil volver a cualquier versión anterior desde Render |

### **Flujo de Despliegue:**

```
Push a main (commit abc1234)
  ↓
Build crea 2 tags:
  - ghcr.io/.../final-isw3:abc1234 (SHA único)
  - ghcr.io/.../final-isw3:latest (alias)
  ↓
Pipeline output: "📦 Image tag: abc1234"
  ↓
QA despliega: "🚀 Deploying to QA with image SHA: abc1234"
  (Render webhook usa tag 'latest', pero sabemos que es abc1234)
  ↓
Tests E2E pasan ✓
  ↓
Aprobación manual en GitHub
  ↓
PROD despliega: "✅ Production deployment approved for SHA: abc1234"
  (Render webhook usa tag 'latest', mismo código que QA testeó)
```

### **Configuración de Render:**

Ambos servicios (QA y PROD) configurados con:
- Image URL: `ghcr.io/margarita0912/final-isw3:latest`
- Auto-Deploy: Enabled (responde a webhooks)

**Nota:** Aunque ambos usan tag `latest`, la **concurrency queue** garantiza que:
- Solo un pipeline modifica `latest` a la vez
- QA termina de testear antes de que otro pipeline actualice la imagen
- PROD despliega el mismo `latest` que QA aprobó (protegido por cola)

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
