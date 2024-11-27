FROM golang:1.22-alpine AS builder
RUN apk add --no-cache git gcc musl-dev
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-w -s" \
    -o bin/gotiny ./cmd/main.go

FROM alpine:3.19
RUN apk add --no-cache ca-certificates tzdata curl
RUN adduser -D -g '' appuser
WORKDIR /app
COPY --from=builder /app/bin/gotiny .
COPY scripts/startup.sh .
RUN chown -R appuser:appuser /app && \
    chmod +x /app/gotiny && \
    chmod +x /app/startup.sh
USER appuser
EXPOSE 8080
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD curl -f http://localhost:${GOTINY_PORT}/health || exit 1
ENTRYPOINT ["/app/startup.sh"]
