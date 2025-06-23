FROM golang:1.24-alpine
RUN apk add --no-cache build-base
WORKDIR /app
COPY . .
CMD ["sh", "-c", "go test ./internal/services && go test ./internal/persistance"]