# Используем официальный образ Go как базовый образ
FROM golang:1.22.1-alpine

# Устанавливаем рабочую директорию внутри контейнера
WORKDIR /app

# Копируем файлы проекта в рабочую директорию
COPY . .

# Скачиваем зависимости
RUN go mod download

# Сборка приложения
RUN go build -o main .

# Указываем команду для запуска контейнера
CMD ["./main"]
