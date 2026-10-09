FROM golang:latest AS builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0  go build -o onte-server ./cmd

FROM cloudflare/cloudflared:latest AS cloudflared

FROM alpine:latest AS runner
WORKDIR /app
COPY --from=builder /app/onte-server .
COPY --from=cloudflared /usr/local/bin/cloudflared /usr/local/bin/cloudflared
EXPOSE 3333
CMD ["./onte-server"]