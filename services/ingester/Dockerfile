# Start with a lightweight Go image for the build stage
FROM golang:1.25-alpine AS builder

# Set the working directory
WORKDIR /app

# Copy and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy all source code from the project root
COPY . .

# Build the ingestion-worker binary
RUN CGO_ENABLED=0 go build -o /ingester ./services/ingester/cmd

# --- Final Stage ---
# Use a minimal 'scratch' image for the final container
FROM alpine:latest

# Copy the compiled binary from the builder stage
COPY --from=builder /ingester /ingester

# The command to run when the container starts. This is the background worker.
CMD ["/ingester"]