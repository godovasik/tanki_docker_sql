FROM golang:1.21-alpine AS builder

WORKDIR /app

# Копируем go.mod и go.sum для эффективного кэширования зависимостей
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходный код
COPY . .

# Собираем приложение
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/main.go

# Финальный образ
FROM alpine:3.18

WORKDIR /app

# Установка зависимостей времени выполнения
RUN apk --no-cache add ca-certificates tzdata

# Копируем исполняемый файл из builder
COPY --from=builder /app/main .

# Открываем порт, который использует приложение
EXPOSE 8080

CMD ["./main"]