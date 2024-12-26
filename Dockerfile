FROM golang:1.22-alpine3.21 AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod tidy
COPY . .
RUN go build -o migrate ./cmd/migrator/
RUN go build -o main ./cmd/s3-mini/

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/main /app/main
COPY --from=builder /app/migrate /app/migrate
COPY --from=builder /app/internal/adapters/repository/db/migrations /app/internal/adapters/repository/db/migrations

ENV HTTP_PORT=8080
ENV FIBER_SERVER_PORT=8081
ENV GRPC_PORT=50051

EXPOSE 8080
EXPOSE 50051
EXPOSE 8081

COPY .env .env
COPY docker_entrypoint.sh /app/docker_entrypoint.sh

RUN chmod +x /app/docker_entrypoint.sh

ENTRYPOINT ["/app/docker_entrypoint.sh"]