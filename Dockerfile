# Указываем базовый образ с Go
FROM golang:1.22 AS builder

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем go.mod и go.sum для кэширования зависимостей
COPY go.mod go.sum ./
RUN go mod download

# Копируем остальные файлы проекта
COPY . .

# Собираем приложение
RUN go build -o main ./cmd/web

# Создаем финальный образ
FROM alpine:latest

# Устанавливаем необходимые зависимости (например, для SQLite)
RUN apk --no-cache add sqlite sqlite-dev

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем бинарник из стадии сборки
COPY --from=builder /app/main .

# Копируем SQL-скрипт new_forum.sql
COPY ./docs/new_forum.sql  /app/docs/

# Запускаем приложение и инициализируем базу данных
CMD ["sh", "-c", "sqlite3 forum.db < /app/docs/new_forum.sql && ./main"]