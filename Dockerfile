# First stage: build the application
FROM golang:1.17-alpine AS builder

ARG GOARCH
ARG GOARM

# Install build utilities
RUN apk --no-cache add --virtual .build-deps \
    bash \
    make \
    git

# Don't use ca-certificates as .build-deps so that we use import certificates
# from builder, in the final image
RUN apk --no-cache add \
    ca-certificates \
    && update-ca-certificates

# Set the Current Working Directory inside the container
WORKDIR /go/src/github.com/CESARBR/knot-thing-copergas

# Copy the source code from the current directory to $WORKDIR (inside the container)
COPY . .

# Port that must be opened
EXPOSE 80
EXPOSE 5672

# Install project development tools and dependencies
RUN go install github.com/ahmetb/govvv@latest
RUN go install github.com/spf13/viper
RUN go install gopkg.in/yaml.v2

# Build the application
RUN make bin

# Remove build dependencies
RUN apk del .build-deps

# Second stage: create the entrypoint to the application binary generated in the previous stage
FROM scratch

WORKDIR /root/

# Copy the configuration files from the build stage
COPY --from=builder /go/src/github.com/CESARBR/knot-thing-copergas/internal/ ./internal

# Copy SSL CA certificates
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy the binary file from the build stage
COPY --from=builder /go/src/github.com/CESARBR/knot-thing-copergas/app-linux-amd64 app

VOLUME /root/internal/config/

ENTRYPOINT ["./app"]
