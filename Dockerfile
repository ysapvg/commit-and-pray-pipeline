# STAGE 1 : Build the Go binary
FROM golang:1.23 AS builder

WORKDIR /app

COPY go.mod ./
COPY main.go ./
COPY main_test.go ./


# Set the default version 
ARG VERSION=dev

# Build the Go binary with the specified version
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-s -w -X main.version=${VERSION}" \ 
    -o app main.go  

# STAGE 2 : Create a minimal image using the scratch base image
FROM scratch
COPY --from=builder /app/app /app
EXPOSE 8080
ENTRYPOINT ["/app"]

