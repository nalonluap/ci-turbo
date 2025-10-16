# Start with a lightweight Go image, ideal for production
FROM golang:1.25-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy Go modules and download dependencies first to leverage Docker layer caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire project source code. This is necessary to access the shared
# /internal/domain directory.
COPY . .

# Build the turbo-observer binary, ensuring it's static for a smaller final image
RUN CGO_ENABLED=0 go build -o /observer ./services/observer/cmd

# --- Final Stage ---
# Use a minimal 'scratch' image for the final container to reduce attack surface
FROM alpine:latest

# Copy the compiled binary from the builder stage
COPY --from=builder /observer /observer

# Expose the port the server will run on
EXPOSE 8080

# The command to run when the container starts
CMD ["/observer"]