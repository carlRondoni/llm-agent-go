# -------- BUILD STAGE --------
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Cache deps
COPY go.mod ./
RUN go mod download

# Copy source
COPY . .

# Build
ARG TARGETOS
ARG TARGETARCH

RUN CGO_ENABLED=0 \
    go build -o app ./cmd/api/main.go

# -------- RUNTIME STAGE --------
FROM alpine:3.23

WORKDIR /app

COPY --from=builder /app/app /app/app

EXPOSE 8080

RUN chmod +x /app/app

ENTRYPOINT ["/app/app"]
