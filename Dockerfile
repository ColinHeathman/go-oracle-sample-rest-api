FROM golang:1.10.1 AS build

# Install dep for go-dependencies
RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh

WORKDIR /go/src/unbc.ca/go-sample-project

COPY main/* ./
COPY Gopkg.* ./
RUN dep ensure
RUN GOOS=linux go build -a -installsuffix cgo -o app .

FROM alpine:latest AS tls

WORKDIR /tls

RUN apk add --no-cache openssl

# Generate self signed root CA cert
RUN openssl req -nodes -x509 -newkey rsa:2048 -keyout ca.key -out ca.crt -days 3650 -subj '/CN=localhost'
# Generate server cert to be signed
RUN openssl req -nodes -newkey rsa:2048 -keyout server.key -out server.csr -days 3650 -subj '/CN=localhost'
# Sign the server cert
RUN openssl x509 -req -in server.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out server.crt

# FROM alpine:latest  
# RUN apk --no-cache add ca-certificates
# WORKDIR /root/
# COPY --from=build /go/src/unbc.ca/go-sample-project/app .
# CMD ["./app"]

FROM scratch

COPY --from=build  /go/src/unbc.ca/go-sample-project/app .
COPY --from=tls    /tls/ca.key      /tls/ca.key
COPY --from=tls    /tls/server.crt  /tls/server.crt

EXPOSE 80
EXPOSE 443
ENTRYPOINT ["/app"]