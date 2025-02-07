# Build stage
FROM node:18-alpine AS node-builder
WORKDIR /app
COPY package*.json ./
COPY tailwind.config.js ./
COPY static/css/styles.css ./static/css/
RUN npm install
RUN npm run build:css

# Go build stage
FROM golang:1.21-alpine AS go-builder
WORKDIR /app
COPY . .
COPY --from=node-builder /app/static/css/tailwind.css ./static/css/
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Final stage
FROM alpine:latest
WORKDIR /app
COPY --from=go-builder /app/main .
COPY --from=go-builder /app/views ./views
COPY --from=go-builder /app/static ./static

# Instalar certificados para conexiones SSL
RUN apk --no-cache add ca-certificates

EXPOSE 8080
CMD ["./main"]
