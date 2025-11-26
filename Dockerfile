# ---------- STAGE 1: Build Backend ----------
FROM golang:1.24 AS backend_builder

WORKDIR /app

COPY ventas-app/go.mod ventas-app/go.sum ./ventas-app/
RUN cd ventas-app && go mod download

COPY ventas-app ./ventas-app
RUN cd ventas-app && go build -o /app/backend-app ./cmd/main.go

# ---------- STAGE 2: Build Frontend ----------
FROM node:20 AS frontend_builder

WORKDIR /app

COPY ventas-frontend/package.json ventas-frontend/package-lock.json ./ventas-frontend/
RUN cd ventas-frontend && npm ci

COPY ventas-frontend ./ventas-frontend
RUN cd ventas-frontend && npm run build

# ---------- STAGE 3: Final Image ----------
FROM node:20-slim

WORKDIR /app

# Ya no instalamos 'serve': el backend sirve el frontend estático

# Copiar backend
COPY --from=backend_builder /app/backend-app ./backend-app

# Copiar frontend build
COPY --from=frontend_builder /app/ventas-frontend/dist ./dist

# Script para iniciar el backend (que también sirve el frontend)
COPY start.sh ./start.sh
RUN chmod +x start.sh

# Render expone un solo puerto (PORT). El backend lo leerá y escuchará ahí.
EXPOSE 8080

CMD ["./start.sh"]
