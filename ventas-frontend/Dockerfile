# -------------------------------------------------------
#  DOCKERFILE FINAL PARA RENDER
#  (Render SOLO sirve /dist, NO construye)
# -------------------------------------------------------

FROM node:20-slim

WORKDIR /app

# Instalar servidor estático
RUN npm install -g serve

# Copiar la carpeta dist generada por GitHub Actions
COPY dist ./dist

EXPOSE 3000

# Servir dist
CMD ["serve", "-s", "dist", "-l", "3000"]
