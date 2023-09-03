FROM golang:1.18-alpine AS builder

WORKDIR /go/src/app

ENV GOPROXY https://goproxy.io

ADD . .

RUN go mod tidy \
    && go build -ldflags "-s -w" -o rsearch main.go

FROM alpine:3.16

ARG TZ="Asia/Shanghai"

COPY --from=builder /go/src/app/rsearch /usr/local/bin/

RUN apk add bash tzdata \
    && ln -sf /usr/share/zoneinfo/${TZ} /etc/localtime \
    && echo ${TZ} > /etc/timezone \
    && chmod +x /usr/local/bin/rsearch

ENTRYPOINT ["/usr/local/bin/rsearch"]

