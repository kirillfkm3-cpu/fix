FROM golang:1.24-alpine

# Устанавливаем зависимости для SQLite
RUN apk add --no-cache gcc musl-dev sqlite-dev

WORKDIR /app

# Копируем и скачиваем зависимости
COPY go.mod go.sum ./
RUN go mod download

# Копируем код
COPY . .

# Собираем приложение и кладем его ВНЕ папки /app
# Мы сохраняем его прямо в системный путь /usr/local/bin/social-server
RUN CGO_ENABLED=1 GOOS=linux go build -o /usr/local/bin/social-server .

# Даем права на запуск
RUN chmod +x /usr/local/bin/social-server

EXPOSE 8080

# Запускаем из нового места
CMD ["/usr/local/bin/social-server"]