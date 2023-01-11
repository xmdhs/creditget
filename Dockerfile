FROM golang:alpine as builder

RUN apk --no-cache add git ca-certificates

WORKDIR /build

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o main -trimpath -ldflags "-w -s" ./server/

FROM alpine:latest

RUN set -ex \
    &&  apk --no-cache add ca-certificates tzdata \
    &&  mkdir /server \
    &&  adduser -H -D server\
    &&  chown -R server /server

USER server
WORKDIR /server
COPY --from=builder /build/main .

ENV PORT "8080"
EXPOSE 8080

CMD ["./main"]