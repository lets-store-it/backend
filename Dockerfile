FROM golang:1.24-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/server ./cmd/server.go
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/init_db ./cmd/init_db/main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/server .
COPY --from=builder /app/init_db .
COPY schema.sql .
EXPOSE 8080
ENV PORT=8080

# Default command runs the server
ENTRYPOINT ["/app/server"]

# To run init_db instead, override both ENTRYPOINT and CMD:
# docker run --rm image:tag /app/init_db -schema /app/schema.sql 
