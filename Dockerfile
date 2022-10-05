# syntax=docker/dockerfile:1
FROM golang:1.18.7-alpine3.16 AS build
WORKDIR /go/src/app
COPY go.* ./
RUN --mount=type=cache,mode=0777,target=/go/pkg/mod \
    --mount=type=cache,mode=0777,target=/root/.cache/build \
    go get -d -v ./...
COPY . ./
RUN --mount=type=cache,mode=0777,target=/go/pkg/mod \
    --mount=type=cache,mode=0777,target=/root/.cache/build \
    go build -o /go/bin/app

FROM alpine:3.16
RUN addgroup -S -g 1001 appgroup && adduser -S -u 1001 -G appgroup appuser
USER appuser
COPY --from=build /go/bin/app /k8s-ecr-login-renew
CMD ["/k8s-ecr-login-renew"]
