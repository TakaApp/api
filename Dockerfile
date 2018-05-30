FROM alpine

ADD server /

RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*

CMD ["/server"]
