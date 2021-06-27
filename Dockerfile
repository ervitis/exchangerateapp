ENV APP_NAME exchangerateapp

FROM golang:1.16-buster as builder

WORKDIR /output
COPY . .

RUN go mod download && CGO_ENABLED=0 GOOS=linux GOARCH=amd54 go build -a -o ${APP_NAME} .

FROM gcr.io/distroless/static:nonroot

WORKDIR /app
COPY --from=builder /output/${APP_NAME} .
USER 1001:1001

ENTRYPOINT ["/app/${APP_NAME}"]
