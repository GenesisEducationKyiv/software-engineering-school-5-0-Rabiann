FROM golang:1.24-alpine AS base
WORKDIR /app
RUN apk add --no-cache build-base
COPY go.mod go.sum ./
RUN go mod download
COPY . .

FROM base AS unit
CMD ["go", "test", "./internal/..."]

FROM base AS integration
CMD ["go", "test", "-tags=integration", "./tests/integration/..."]

FROM base AS e2e
CMD ["go", "test", "-tags=e2e", "./tests/e2e/..."]