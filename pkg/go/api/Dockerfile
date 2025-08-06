FROM golang:alpine AS builder
WORKDIR /build/pkg/go/api
COPY pkg/go/api/v1 ./v1
WORKDIR /build
COPY lib ./lib
COPY protos ./protos
COPY .env .
COPY templates ./templates
COPY static ./static
WORKDIR /build/pkg/go/api/v1
RUN go build -o app .

FROM alpine

WORKDIR /build

COPY --from=builder /build/pkg/go/api/v1/app   app
COPY --from=builder /build/templates           templates
COPY --from=builder /build/static              static
COPY --from=builder /build/.env                .env
CMD ["./app"]
