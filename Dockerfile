FROM golang:1.22.1-alpine AS base

RUN apk update && apk add --no-cache ca-certificates

FROM base AS builder
WORKDIR /app
COPY . .
RUN go build -o server cmd/server/main.go

FROM alpine:3.19 AS final
WORKDIR /app
COPY --from=builder /app/server .
COPY .env .

EXPOSE 3530

CMD ["/app/server"]

