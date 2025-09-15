# Step 1: Build stage
FROM golang:1.24.5 AS build

WORKDIR /app

# Copy go.mod and go.sum for caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the code
COPY . .

# Build the binary from main.go
RUN CGO_ENABLED=0 GOOS=linux go build -o out ./cmd/server/main.go

# Step 2: Minimal runtime image
FROM debian:bookworm-slim

WORKDIR /app
COPY --from=build /app/out /out

# Ensure executable
RUN chmod +x /out

EXPOSE 8080
ENV PORT=8080

CMD ["/out"]
