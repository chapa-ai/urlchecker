FROM golang:1.21.5-alpine as build-stage

RUN mkdir -p /checker

WORKDIR /checker

COPY . /checker
RUN go mod download

RUN go build -o info cmd/main.go

ENTRYPOINT [ "/checker/info" ]