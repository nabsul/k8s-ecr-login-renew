FROM golang:1.14 as build

WORKDIR /app
COPY . .

RUN go build -o k8s-awsreg-renew .

FROM debian
RUN apt-get update
RUN apt-get install -y ca-certificates
RUN update-ca-certificates

COPY --from=build /app/k8s-awsreg-renew .
CMD ["./k8s-awsreg-renew"]
