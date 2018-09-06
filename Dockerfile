FROM golang:1.10.1 AS build

WORKDIR /go/src/unbc.ca/go-sample-project

RUN go get -d -v \
    github.com/gorilla/mux \
    gopkg.in/goracle.v2

COPY app.go .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

FROM alpine:latest  
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=build /go/src/unbc.ca/go-sample-project/app .
CMD ["./app"]

# FROM scratch
# COPY --from=build /go/src/unbc.ca/go-sample-project/app .
# EXPOSE 80
# ENTRYPOINT ["/app"]