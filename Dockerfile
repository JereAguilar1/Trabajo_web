# Etapa 1: compilación
FROM golang:1.27-alpine AS builder
WORKDIR /app
# sqlc genera db/sqlc/, que no se versiona en el repo
RUN go install github.com/sqlc-dev/sqlc/cmd/sqlc@v1.31.1
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN sqlc generate
RUN CGO_ENABLED=0 go build -o /app/api .

# Etapa 2: imagen final (solo el binario)
FROM alpine:3.24
RUN adduser -D -u 1000 app
USER app
COPY --from=builder /app/api /usr/local/bin/api
EXPOSE 8090
CMD ["api"]
