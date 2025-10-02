# syntax=docker/dockerfile:1
FROM oven/bun:1 AS dashboard
WORKDIR /usr/src/app
COPY . .
RUN bun install --production
RUN bun run build

FROM golang:1.24-alpine AS builder
WORKDIR /app
RUN go install github.com/swaggo/swag/cmd/swag@latest
COPY go.mod go.sum ./
RUN go mod download
COPY . .
# Generate Swagger docs
RUN export PATH=$PATH:$(go env GOPATH)/bin && swag init --generalInfo main.go --output docs
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./kinctl -tags prod -mod=readonly -ldflags "-s -w" ./main.go

FROM alpine:3.20
WORKDIR /var/www
COPY --from=dashboard /usr/src/app/dist .
COPY --from=builder /app/kinctl /usr/local/bin/kinctl
COPY --from=builder /app/docs /usr/local/bin/docs
EXPOSE 3000
CMD ["kinctl", "--server", "--migrate", "--port", "3000"]
