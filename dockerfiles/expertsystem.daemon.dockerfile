FROM golang:alpine

# Switch to our app directory
RUN mkdir -p /go/src/github.com/ALSAD-project/alsad-core
WORKDIR /go/src/github.com/ALSAD-project/alsad-core

# Copy the source files
COPY ./cmd/expertsystem /go/src/github.com/ALSAD-project/alsad-core/cmd/expertsystem
COPY ./pkg/expertsystem /go/src/github.com/ALSAD-project/alsad-core/pkg/expertsystem

# Save ENV
# ENV REQUEST_PORT 4000

RUN go build ./cmd/expertsystem/daemon/main.go
RUN cp main /usr/local/bin/expertsystem-daemon
WORKDIR /

# EXPOSE 4000
CMD ["expertsystem-daemon"]