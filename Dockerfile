FROM golang:1.27.1-alpine AS builder

WORKDIR /app

COPY go.mod go.sum* ./
COPY cmd/ ./cmd/
COPY internal/ ./internal/

RUN CGO_ENABLED=0 GOOS=linux go build \
    -trimpath \
    -ldflags="-s -w" \
    -o /app/server ./cmd/server

FROM alpine:3.22 AS runner

RUN apk add --no-cache ca-certificates \
    && addgroup -S app \
    && adduser -S app -G app

WORKDIR /app

COPY --from=builder /app/server ./server

USER app

EXPOSE 8003

CMD ["./server"]
