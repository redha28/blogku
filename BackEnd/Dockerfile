FROM golang:latest AS builder

# Create working directory
WORKDIR /app

# Copy everything
COPY . .

# Fix go.mod directly (remove Go version requirement)
RUN sed -i 's/go 1.23.5/go 1.20/g' go.mod

# Run go mod tidy as suggested by the error message
RUN go mod tidy

# Build the application with -mod=mod to force using the go.mod file
RUN go build -o blogku ./cmd/main.go

# Use a clean Alpine for the final image
FROM alpine:latest

# Set working directory
WORKDIR /app

# Install necessary dependencies
RUN apk --no-cache add ca-certificates tzdata

# Copy only what's needed from the builder
COPY --from=builder /app/blogku .
COPY --from=builder /app/migrations ./migrations
COPY --from=builder /app/docs ./docs

# Create logs directory
RUN mkdir -p /app/logs && \
    mkdir -p /app/public/uploads

# Copy environment file
COPY .env.docker /app/.env

EXPOSE 8080

CMD ["./blogku"]
