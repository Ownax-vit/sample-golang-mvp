FROM golang:1.25-alpine AS modules

COPY go.mod go.sum /modules/

WORKDIR /modules

RUN go mod download

FROM golang:1.25-alpine AS debug
# dev stage with delve debugger

COPY --from=modules /go/pkg /go/pkg
COPY . /app

WORKDIR /app

RUN go install github.com/go-delve/delve/cmd/dlv@latest
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -gcflags=all="-N -l" -v -o ./app ./cmd/main.go
EXPOSE 2345 80
CMD ["dlv", "--headless=true", "--accept-multiclient", "--continue", "--listen=0.0.0.0:2345", "--api-version=2", "exec", "./app", "--log"]
FROM golang:1.25-alpine AS builder

COPY --from=modules /go/pkg /go/pkg
COPY . /app

WORKDIR /app

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -v -o /bin/app ./cmd/main.go

FROM scratch

COPY --from=builder /app/config /config
COPY --from=builder /bin/app /app
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

CMD ["/app"]
