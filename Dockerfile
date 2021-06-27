ARG API_KEY

FROM golang:1.16-buster as builder

WORKDIR /output
COPY . .

RUN go mod download && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -gcflags="all=-N -l" -a -o exchangerateapp .

FROM debian:buster

WORKDIR /app
COPY --from=builder /output/exchangerateapp .
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

ARG API_KEY
ENV API_KEY ${API_KEY}

EXPOSE 8080

ENTRYPOINT ["/app/exchangerateapp"]
