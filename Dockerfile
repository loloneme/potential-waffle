FROM golang:1.25.1-alpine

WORKDIR /app

COPY . .

RUN go build -o /build ./cmd/service/main.go \
    && go clean -cache -modcache

EXPOSE 8080

CMD ["/build"]