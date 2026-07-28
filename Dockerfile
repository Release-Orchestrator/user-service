# Build stage
FROM golang:1.26-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /user-service ./cmd

# Runtime stage
FROM gcr.io/distroless/static-debian12:nonroot

WORKDIR /app

COPY --from=builder /user-service .

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["./user-service"]