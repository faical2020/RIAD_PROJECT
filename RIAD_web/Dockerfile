FROM golang:1.26-alpine AS builder
WORKDIR /app
# Copy go.mod from the RIAD_web directory
COPY RIAD_web/go.mod ./
RUN go mod download
# Copy main.go from the RIAD_web directory
COPY RIAD_web/main.go .
RUN CGO_ENABLED=0 GOOS=linux go build -o riad-web main.go

FROM alpine:latest
WORKDIR /root/
# Copy binary from builder
COPY --from=builder /app/riad-web .
# Copy static assets from RIAD_APP/frontend/dist
COPY RIAD_APP/frontend/dist ./dist

EXPOSE 8080
CMD ["./riad-web"]
