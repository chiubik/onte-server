FROM golang:1.24.4 AS builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 go build -o onte-server ./cmd

FROM cloudflare/cloudflared:latest AS cloudflared

FROM alpine:3.22 AS runner
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=builder /app/onte-server .
COPY --from=cloudflared /usr/local/bin/cloudflared /usr/local/bin/cloudflared
CMD ["./onte-server"]
