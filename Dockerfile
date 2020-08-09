FROM golang:1.14-alpine as build

COPY ./ /go/src/github.com/meyskens/wap.maartje.dev

WORKDIR /go/src/github.com/meyskens/wap.maartje.dev

RUN go build -o server ./

FROM alpine:3.12

RUN apk add --no-cache ca-certificates

RUN mkdir /opt/wap.maartje.dev
WORKDIR /opt/wap.maartje.dev

COPY --from=build /go/src/github.com/meyskens/wap.maartje.dev/server /usr/local/bin
COPY --from=build /go/src/github.com/meyskens/wap.maartje.dev/static /opt/wap.maartje.dev/static

ENTRYPOINT [ "server" ]