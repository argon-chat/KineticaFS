# syntax=docker/dockerfile:1

# ----------- Build Stage -----------
FROM golang:1.24-alpine AS builder
WORKDIR /app

# Install swag for swagger doc generation
RUN go install github.com/swaggo/swag/cmd/swag@latest

# Copy go mod and sum files
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Generate Swagger docs
RUN export PATH=$PATH:$(go env GOPATH)/bin && swag init --generalInfo main.go --output docs

# Build the Go app
RUN go build -o /app/kineticafs ./main.go

# ----------- Run Stage -----------
FROM alpine:3.20
WORKDIR /app

# Install ca-certificates for HTTPS
RUN apk add --no-cache ca-certificates

# Copy built binary and docs
COPY --from=builder /app/kineticafs .
COPY --from=builder /app/docs ./docs

# Expose default port
EXPOSE 3000

CMD ["./kineticafs", "--server", "--port", "3000"]
