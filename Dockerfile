FROM golang:1.14.1
RUN mkdir /app
ADD . /app
WORKDIR /app
RUN go build -o klp ./cmd/kiddyLineProcessor/.