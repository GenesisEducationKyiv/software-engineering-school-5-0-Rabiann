FROM golang:alpine AS builder
WORKDIR /build
COPY . .

RUN go build -o api .

FROM alpine

WORKDIR /build

COPY --from=builder /build/api       api
COPY --from=builder /build/templates templates
COPY --from=builder /build/static    static
COPY --from=builder /build/.env      .env
CMD ["./api"]