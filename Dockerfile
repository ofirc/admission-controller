# Build stage
FROM golang as builder

WORKDIR /app
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o webhook-server ./cmd/webhook-server

# Production image using Google's Distroless base
FROM gcr.io/distroless/static as distroless
COPY --from=builder /app/webhook-server /
CMD ["/webhook-server"]

# Development image using Alpine
FROM alpine:latest as development
COPY --from=builder /app/webhook-server /
RUN apk --no-cache add bash curl
CMD ["/webhook-server"]