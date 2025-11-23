# -------- STAGE 1: Build --------
FROM node:20 AS builder

WORKDIR /app
COPY package*.json ./
RUN npm ci

COPY . .
RUN npm run build


# -------- STAGE 2: Serve --------
FROM node:20-slim

WORKDIR /app

# Instalar un servidor estático
RUN npm install -g serve

COPY --from=builder /app/dist ./dist

EXPOSE 3000
CMD ["serve", "-s", "dist", "-l", "3000"]