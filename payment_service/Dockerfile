# Используем многоступенчатую сборку
FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /payment-service ./cmd/main.go

FROM alpine:3.18

WORKDIR /root/

COPY --from=builder /payment-service .
COPY --from=builder /app/migrations ./migrations

# Устанавливаем tzdata для работы с временными зонами
RUN apk --no-cache add tzdata

# Настройки
ENV TZ=Europe/Moscow
EXPOSE 50053

# Запуск приложения
CMD ["./payment-service"]