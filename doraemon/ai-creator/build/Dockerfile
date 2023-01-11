FROM golang:1.19 as builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
RUN go build -o bin/ai-creator ./cmd/main.go

# Command to run the executable
#CMD bin/ai-creator 


FROM nicolaka/netshoot
RUN mkdir -p /app && \
    chown -R nobody:nogroup /app
COPY --from=builder /app/bin/ai-creator /app
COPY --from=builder /app/templates /app/templates
USER nobody
WORKDIR /app
# Expose the 8123 port
EXPOSE 8123
ENTRYPOINT  ["/app/ai-creator"]