FROM golang:1.25.1-alpine

WORKDIR /app



# Копируем исходный код
COPY . .

# Собираем приложение
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /build ./cmd/service/main.go \
    && go clean -cache -modcache

RUN chmod +x /build

EXPOSE 8080

CMD ["/build"]