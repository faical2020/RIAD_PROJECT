FROM golang:1.26-alpine AS builder
WORKDIR /app
# Copy proto to the root of the builder to satisfy the '../proto/sync' replacement
COPY proto /proto
COPY RIAD_SERVER/go.mod RIAD_SERVER/go.sum ./
RUN go mod download
COPY RIAD_SERVER .
RUN CGO_ENABLED=0 GOOS=linux go build -o riad-server cmd/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/riad-server .
EXPOSE 8080
CMD ["./riad-server"]
