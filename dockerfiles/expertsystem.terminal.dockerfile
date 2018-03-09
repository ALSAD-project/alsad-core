FROM golang:alpine
RUN \
     apk update \
  && apk add curl git \
  && rm -rf /var/cache/apk/*

RUN curl https://glide.sh/get | sh

# Switch to our app directory
RUN mkdir -p /go/src/github.com/ALSAD-project/alsad-core
WORKDIR /go/src/github.com/ALSAD-project/alsad-core

# Copy the source files
COPY . /go/src/github.com/ALSAD-project/alsad-core

RUN glide install
RUN go build ./cmd/expertsystem/terminal/main.go
RUN cp main /usr/local/bin/expertsystem-terminal
WORKDIR /

CMD ["expertsystem-terminal"]
