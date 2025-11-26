# ---------- STAGE 1: Build ----------
FROM node:20 AS builder

WORKDIR /app

# Instalamos dependencias con caché correcto
COPY package.json package-lock.json ./
RUN npm ci

# Copiamos el resto del proyecto
COPY . .

# Build del frontend
RUN npm run build

# ---------- STAGE 2: Static Server ----------
FROM node:20-slim

WORKDIR /app

# Servidor estático
RUN npm install -g serve

# Copiamos build
COPY --from=builder /app/dist ./dist

# Copiamos archivo config.json (el que se sobrescribe en runtime)
COPY public/config.json ./dist/config.json

EXPOSE 3000

CMD ["serve", "-s", "dist", "-l", "3000"]
