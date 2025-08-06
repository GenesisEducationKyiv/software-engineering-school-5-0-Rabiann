FROM golang:alpine AS builder
WORKDIR /build/pkg/go/weather
COPY pkg/go/weather/v1 ./v1
WORKDIR /build
COPY lib ./lib
COPY protos ./protos
COPY .env .
WORKDIR /build/pkg/go/weather/v1
RUN go build -o app .

FROM alpine

WORKDIR /build

COPY --from=builder /build/pkg/go/weather/v1/app   app
COPY --from=builder /build/.env                         .env
CMD ["./app"]
