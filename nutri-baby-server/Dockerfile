FROM golang:1.24-alpine AS builder

WORKDIR /app

# Install git for fetch dependencies
RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o nutri-baby-server cmd/server/main.go

FROM alpine:latest

WORKDIR /app

# Install basic dependencies (optional, e.g., for timezones)
RUN apk --no-cache add ca-certificates tzdata

COPY --from=builder /app/nutri-baby-server .
COPY --from=builder /app/config ./config

# Expose port
EXPOSE 8080

# Run the binary
CMD ["./nutri-baby-server"]
