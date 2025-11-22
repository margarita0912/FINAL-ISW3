# Testing checklist

Este archivo resume las verificaciones realizadas sobre el repo y marca qué está listo y qué queda para el futuro.

## Resumen rápido
- Entornos (QA / PROD): IMPLEMENTADO
  - Backend: `ventas-app` carga `.env.qa` y `.env.prod` con `godotenv`.
  - Frontend: `EntornoSelector` en `ventas-frontend/src/components/EntornoSelector.tsx` usa `localStorage` y por defecto `qa`.
  - Axios añade header `X-Env` en `ventas-frontend/src/api/axios.ts`.

- Smoke tests: POSPUESTOS (se ejecutarán después del deploy a QA)
  - No hay smoke tests explícitos en el repo actualmente.
  - Recomendación: crear un spec `cypress/e2e/smoke.cy.js` o un job `smoke` que haga `curl -f $QA_URL/health`.
  - También: `acceptance` (specs Cypress ligeros) POSPUESTOS hasta QA/PROD.

- Pruebas unitarias: IMPLEMENTADO
  - Backend (Go): tests en `ventas-app/controllers/` — archivos `*_test.go` detectados.
  - Frontend (Jest): tests en `ventas-frontend/tests/` y `ventas-frontend/src/tests/`.
  - Cypress: tests E2E presentes en `ventas-frontend/cypress/e2e/`.

## Checklist (estado actual)
- [x] Detectar configuraciones por entorno (backend y frontend)
- [ ] Encontrar smoke tests (no existen, posponer hasta QA)
 - [ ] Encontrar smoke tests (no existen, posponer hasta QA)
 - [ ] Añadir job `acceptance` en `.gitlab-ci.yml` (futuro, ejecutar tras deploy a QA)
- [x] Verificar pruebas unitarias (backend + frontend)
- [x] Crear scripts locales para ejecutar pipeline (`run_pipeline.ps1`, `run_pipeline.sh`)
- [ ] Añadir job `smoke` en `.gitlab-ci.yml` (futuro, ejecutar tras deploy a QA)

## Comandos útiles (local)
- Ejecutar pipeline local (PowerShell):

```powershell
pwsh .\run_pipeline.ps1
# o, si no tienes pwsh instalado con PowerShell 5.1
powershell -NoProfile -ExecutionPolicy Bypass -File .\run_pipeline.ps1
```

- Ejecutar pipeline local (Bash/WSL):

```bash
chmod +x ./run_pipeline.sh
./run_pipeline.sh
```

- Ejecutar sólo tests backend (local):

```bash
cd ventas-app
go test ./... -v -coverprofile=coverage.out -covermode=atomic
```

- Ejecutar sólo tests frontend (local):

```bash
cd ventas-frontend
npm ci
npm run test:ci
```

- Ejecutar compilación (comprobación de compilación):
```bash
cd ventas-app
go build ./... # devuelve error si no compila
cd ../ventas-frontend
npm run build # para verificar compilación frontend
```

## Acceptance y Smoke (documentación rápida)

- `smoke` (job manual): verifica que `QA_URL/health` responda 2xx mediante `curl`. Debe definirse `QA_URL` en Variables del proyecto para que funcione.

- `acceptance` (job manual): ejecuta un spec ligero de Cypress `cypress/e2e/acceptance.cy.js` contra `QA_URL`. Se ejecuta manualmente desde la UI.

## Snippet para re-agregar jobs `smoke` y `acceptance` (pegar en `.gitlab-ci.yml` cuando QA/PROD esté disponible)

Copiar y pegar el siguiente bloque YAML en `.gitlab-ci.yml` (ajusta `QA_URL` y dependencias según tu flujo de despliegue):

```yaml
# --- smoke (manual) ---
smoke:
  stage: smoke
  image: curlimages/curl:8.3.0
  needs:
    - job: build_backend
      artifacts: true
    - job: build_frontend
      artifacts: true
  variables:
    QA_URL: ""
  script:
    - |
      if [ -z "$QA_URL" ]; then
        echo "ERROR: QA_URL is not set. Set QA_URL in CI variables to point to QA environment." >&2
        exit 1
      fi
    - curl -sSf "$QA_URL/health" -o /dev/null
  when: manual
  only:
    - main

# --- acceptance (manual, Cypress) ---
acceptance:
  stage: acceptance
  image: node:20
  needs:
    - job: build_backend
      artifacts: true
    - job: build_frontend
      artifacts: true
  variables:
    QA_URL: ""
  script:
    - |
      if [ -z "$QA_URL" ]; then
        echo "ERROR: QA_URL is not set. Set QA_URL in CI variables to point to QA environment." >&2
        exit 1
      fi
    - cd ${FRONTEND_DIR}
    - npm ci
    - npx cypress run --spec "cypress/e2e/acceptance.cy.js" --config baseUrl=$QA_URL --env apiUrl=$QA_URL
  when: manual
  only:
    - main
```

## Cómo disparar los jobs manuales

- En GitLab UI: CI / Pipelines → seleccionar pipeline → en la etapa `smoke` o `acceptance` pulsar ▶️ (Play).
- Por API: usar el endpoint `POST /projects/:id/jobs/:job_id/play` con `PRIVATE-TOKEN`.

## Notas sobre pruebas de aceptación

- Actualmente el spec de acceptance es mínimo (ver `ventas-frontend/cypress/e2e/acceptance.cy.js`) y comprueba que la app carga y que el selector de entorno existe.
- Si querés pruebas de aceptación reales que incluyan lógica de negocio (crear venta, comprobar stock, etc.), necesito:
  - Endpoints / payloads de ejemplo para crear ventas en QA
  - Credenciales o usuario de test para QA (o configurar mocks en CI)


## Recomendaciones / Próximos pasos
1. (Corto plazo) Documentar las variables CI necesarias para QA y PROD en `README_LOCAL_PIPELINE.md` o `README.md`:
   - `QA_URL`, `QA_DB_HOST`, `QA_DB_USER`, `QA_DB_PASS`, `QA_DB_NAME`, `SONAR_TOKEN`, `SONAR_ORG`, `SONAR_PROJECT_KEY`.

2. (Medio plazo) Crear job `smoke` en CI que se ejecute tras el deploy a QA. Opciones:
   - `curl -f $QA_URL/health` (rápido) 
   - `npx cypress run --spec "cypress/e2e/smoke.cy.js"` (más robusto)

3. (Opcional) Añadir un pequeño spec Cypress `cypress/e2e/smoke.cy.js` con un test que valide login y una petición mínima.

4. (Opcional) Integrar el job `smoke` en el pipeline `.gitlab-ci.yml` dentro de la etapa `e2e_tests` o como etapa separada `smoke` que dependa del despliegue QA.

Si querés, implemento ahora la tarea 2 o 3 (crear job `smoke` o crear el spec Cypress mínimo). Si preferís dejarlo para cuando se despliegue QA, lo dejamos pospuesto en la lista.
