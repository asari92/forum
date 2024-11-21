# Используем официальный образ Go 1.22 как базовый
FROM golang:1.22 as builder

# Устанавливаем рабочую директорию внутри контейнера
WORKDIR /app

# Копируем файлы проекта в контейнер
COPY go.mod go.sum ./
RUN go mod download

# Копируем остальной код проекта
COPY . .

# Сборка приложения
RUN go build -o forum ./cmd/web/main.go

# Создаём минимальный образ для запуска приложения
FROM debian:bullseye-slim

# Устанавливаем зависимости для SQLite3
RUN apt-get update && apt-get install -y \
    ca-certificates \
    sqlite3 \
    && rm -rf /var/lib/apt/lists/*

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем собранный бинарник из стадии сборки
COPY --from=builder /app/forum .

# Экспонируем порт 4000
EXPOSE 4000

# Устанавливаем команду запуска
CMD ["./forum"]
