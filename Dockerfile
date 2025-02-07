# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY . .

# Instalar dependencias de node para Tailwind
RUN apk add --no-cache nodejs npm
RUN npm install
RUN npm run build:css

# Compilar la aplicación Go
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Final stage
FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/views ./views
COPY --from=builder /app/static ./static
COPY --from=builder /app/.env.example ./.env

# Instalar certificados para conexiones SSL
RUN apk --no-cache add ca-certificates

EXPOSE 8080

CMD ["./main"]
