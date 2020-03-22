FROM golang:1.14 as build

WORKDIR /app
COPY . .
RUN go build -o k8s-ecr-login-renew .

FROM alpine:latest
RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*
COPY --from=build /app/k8s-ecr-login-renew .
CMD ["./k8s-awsreg-renew"]
