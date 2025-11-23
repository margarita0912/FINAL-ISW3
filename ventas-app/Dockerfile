# -------- STAGE 1: Build --------
FROM golang:1.22 AS builder

WORKDIR /app

# Copiamos go.mod y go.sum primero para aprovechar cache
COPY go.mod go.sum ./
RUN go mod download

# Copiamos el resto
COPY . .

# Compilar binario estático
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/main.go


# -------- STAGE 2: Runtime --------
FROM gcr.io/distroless/base-debian12

WORKDIR /app
COPY --from=builder /app/server .

EXPOSE 8080
CMD ["./server"]