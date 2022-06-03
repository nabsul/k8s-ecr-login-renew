# syntax=docker/dockerfile:1
FROM golang:1.18-alpine AS build
WORKDIR /go/src/app
ADD . /go/src/app
RUN go get -d -v ./...
RUN go build -o /go/bin/app

FROM alpine:3.16
COPY --from=build /go/bin/app /k8s-ecr-login-renew
CMD ["/k8s-ecr-login-renew"]
