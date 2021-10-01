FROM golang:1.17.0-alpine as builder

RUN apk --no-cache add ca-certificates gcc libc-dev

VOLUME /root/go
COPY ./ /app
RUN cd /app &&\
  go test ./... &&\
  go build

FROM alpine:3.14.2
WORKDIR /app
COPY --from=builder /app/mapcleaner /bin/mapcleaner

CMD ["/bin/mapcleaner"]
