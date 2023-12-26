# This is a multi-stage Dockerfile. The first part executes a build in a Golang
# container, and the second retrieves the binary from the build container and
# inserts it into a "scratch" image.

FROM golang:1.20-alpine as builder

# Set the Current Working Directory inside the container.
WORKDIR /go/src/app

# Copy go mod and sum files to download dependencies.
# Adjust the path to copy from the project root.
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire project into the container.
# Adjust the path to include the entire project directory.
COPY ../ ./


# Build the binary. Note the flags that we use here:
#  CGO_ENABLED=0 --> Do not use CGO; compile statically
#  GOOS=linux    --> Compile for Linux OS
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app ./cmd/main.go

FROM ubuntu:latest as user
RUN useradd -u 10001 app-user
RUN touch  /transactions.log && chown app-user /transactions.log

# Note that we use a "scratch" image, which contains no distribution files. The
# resulting image and containers will have only one file: our service binary.
FROM scratch as image

# Copy the binary from the builder container.
COPY --from=builder /go/src/app/app .

COPY --from=user /etc/passwd /etc/passwd

COPY --from=user /transactions.log /transaction.log

USER app-user

EXPOSE 8080

CMD ["/app"]

