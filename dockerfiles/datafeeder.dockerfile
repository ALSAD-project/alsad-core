FROM golang:1.9-alpine

RUN \
     apk update \
  && apk add curl git \
  && rm -rf /var/cache/apk/*

RUN curl https://glide.sh/get | sh

# Switch to our app directory
RUN mkdir -p /go/src/github.com/ALSAD-project/alsad-core
WORKDIR /go/src/github.com/ALSAD-project/alsad-core

# Copy the source files
COPY ./cmd/datafeeder /go/src/github.com/ALSAD-project/alsad-core/cmd/datafeeder
COPY ./pkg/datafeeder /go/src/github.com/ALSAD-project/alsad-core/pkg/datafeeder
COPY ./glide.yaml /go/src/github.com/ALSAD-project/alsad-core/glide.yaml


# As gonum is not yet a package
RUN go get -u -t gonum.org/v1/gonum/...
RUN go get github.com/kelseyhightower/envconfig
# RUN glide install


RUN go build ./cmd/datafeeder/datafeeder.go
RUN cp datafeeder /usr/local/bin/datafeeder
WORKDIR /


CMD ["datafeeder"]