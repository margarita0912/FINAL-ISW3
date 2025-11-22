# Local Pipeline Runner

These scripts let you run the CI pipeline locally (tests first, builds after).

Files added by this repo:
- `run_pipeline.ps1` - PowerShell script (Windows)
- `run_pipeline.sh` - Bash script (Unix/macOS)

Prerequisites
- Go 1.24 installed and on `PATH`.
- Node (v20 recommended) and `npm` installed.
- Git (optional) and standard build tools.
- If you want to run Cypress E2E: install Cypress (`npm ci` will install it) and have QA services running.

Usage (PowerShell - Windows)
```powershell
# from project root
powershell -NoProfile -ExecutionPolicy Bypass -File .\run_pipeline.ps1
# to also run Cypress E2E (QA services must be running):
powershell -NoProfile -ExecutionPolicy Bypass -File .\run_pipeline.ps1 -RunE2E
```

Usage (Bash - Unix/macOS)
```bash
# make script executable once
chmod +x ./run_pipeline.sh
# run
./run_pipeline.sh
```

What the scripts do
1. Run backend unit tests with coverage and generate `coverage.xml` (Cobertura).
2. Run frontend unit tests using `npm run test:ci`.
3. If both test steps pass, build the backend binary and build the frontend (`npm run build`).
4. Optional: run Cypress E2E (PowerShell supports `-RunE2E`).

Notes
- Scripts are "fail-fast": if tests fail the script exits with non-zero code and skips builds.
- Scripts call the same commands as `.gitlab-ci.yml` to keep behavior similar to CI.
- Adjust ports/URLs for QA in `ventas-frontend/cypress.config.ts` if needed.

## Variables CI recomendadas
En GitLab (Project → Settings → CI/CD → Variables) definan al menos:

- `QA_URL` : URL base del entorno QA (ej: `https://qa.example.com`).
- `QA_DB_HOST`, `QA_DB_USER`, `QA_DB_PASS`, `QA_DB_NAME` : si los jobs necesitan conectarse a la base de datos QA.
- `SONAR_TOKEN`, `SONAR_ORG`, `SONAR_PROJECT_KEY` : para análisis en SonarCloud.

Al definir `QA_URL`, los jobs manuales `smoke` y `acceptance` podrán usarse desde la UI.

## Pruebas y compilación - qué hace cada script/job
- Pruebas unitarias backend: `go test ./... -v -coverprofile=coverage.out -covermode=atomic` (genera `coverage.out` y luego `coverage.xml`).
- Pruebas unitarias frontend: `npm run test:ci` (config de Jest en `ventas-frontend`).
- Pruebas de compilación (build): `go build -o ../backend-app ./cmd/main.go` y `npm run build` para frontend.
- Pruebas de aceptación: job `acceptance` (manual) que ejecuta `npx cypress run --spec "cypress/e2e/acceptance.cy.js"` apuntando a `QA_URL`.
- Smoke checks: job `smoke` (manual) que hace `curl -sSf $QA_URL/health`.

## Cómo ejecutar solo pasos concretos localmente
- Solo pruebas backend (local):
```powershell
cd ventas-app
go test ./... -v -coverprofile=coverage.out -covermode=atomic
```
- Solo compilación backend:
```powershell
cd ventas-app
go build -o ../backend-app ./cmd/main.go
```
- Solo pruebas frontend (local):
```bash
cd ventas-frontend
npm ci
npm run test:ci
```
- Ejecutar spec de aceptación local (requiere frontend/server QA o usar `--config baseUrl`):
```bash
cd ventas-frontend
npx cypress run --spec "cypress/e2e/acceptance.cy.js" --config baseUrl=http://localhost:5174
```

## Ejecutar jobs manuales en GitLab
- Desde la UI: CI / Pipelines → abrir pipeline → en la etapa `smoke` o `acceptance` hacer clic en ▶️ (Play) para ejecutar el job manual.
- Usando la API (avanzado):
```bash
curl --request POST \
  --header "PRIVATE-TOKEN: $GITLAB_TOKEN" \
  "https://gitlab.com/api/v4/projects/<project_id>/jobs/<job_id>/play"
```

Recomendación: definas `QA_URL` en Variables del proyecto antes de ejecutar los jobs manuales.

Troubleshooting
- If `gocover-cobertura` is not found, the scripts will attempt to `go install` it.
- Ensure `npm ci` can run (network, registry auth, etc.).
- If running E2E, make sure QA backend/frontend services are available at the QA URLs.

If you'd like, I can:
- Add flags to run only backend or only frontend tests.
- Integrate Sonar scan step (requires SONAR_TOKEN + network access).