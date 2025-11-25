# ==== STAGE 1: BUILD ====
FROM node:20-slim AS builder

WORKDIR /app

COPY package.json package-lock.json ./
RUN npm ci

COPY . .

# Recibe el valor desde Render → build-arg
ARG VITE_API_URL
ENV VITE_API_URL=$VITE_API_URL

# Crear archivo .env para Vite
RUN echo "VITE_API_URL=$VITE_API_URL" > .env

RUN npm run build


# ==== STAGE 2: SERVE DIST ====
FROM node:20-slim

WORKDIR /app

# Servidor estático
RUN npm install -g serve

COPY --from=builder /app/dist ./dist

EXPOSE 3000
CMD ["serve", "-s", "dist", "-l", "3000"]
