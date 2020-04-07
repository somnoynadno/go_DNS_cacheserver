FROM golang:1.9 as builder

RUN mkdir -p /go/src/go_DNS_cacheserver
WORKDIR /go/src/go_DNS_cacheserver

COPY . .

RUN go get -d
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o app .


FROM alpine:latest
WORKDIR /

COPY --from=builder /go/src/go_DNS_cacheserver/app .

CMD ["/app"]

EXPOSE 5153