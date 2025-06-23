FROM golang:1.18-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum first, then download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

RUN CGO_ENABLED=0 go build -o loggerServiceApp ./cmd/api

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/loggerServiceApp .

CMD ["/app/loggerServiceApp"]