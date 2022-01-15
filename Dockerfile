FROM golang:1.17-alpine as builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .
RUN go build -o ./bin/app cmd/main.go

FROM alpine:3.14

WORKDIR /app
COPY --from=builder /app/config/config-local.yml ./config/config-local.yml
COPY --from=builder /app/bin/app /app/serve

ENTRYPOINT /app/serve
