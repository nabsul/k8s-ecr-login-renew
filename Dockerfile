# syntax=docker/dockerfile:1
FROM golang:1.20-alpine AS build
WORKDIR /go/src/app
COPY go.mod /go/src/app/
COPY go.sum /go/src/app/
RUN go mod download

COPY . .
RUN go build -o /go/bin/app

FROM alpine:3.17
RUN addgroup -S -g 1001 appgroup && adduser -S -u 1001 -G appgroup appuser
USER appuser

COPY --from=build /go/bin/app /k8s-ecr-login-renew
CMD ["/k8s-ecr-login-renew"]
