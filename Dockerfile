FROM golang:alpine as builder

RUN apk --no-cache add git ca-certificates

WORKDIR /build

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o main -trimpath -ldflags "-w -s" ./server/

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /build/main .

ENV PORT "8080"

EXPOSE 8080

CMD ["./main"]