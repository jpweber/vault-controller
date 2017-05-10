FROM alpine:3.5
MAINTAINER Kelsey Hightower <kelsey.hightower@gmail.com> Jim Weber <jpweber@gmail.com>
ADD vault-init /vault-init
RUN apk --update upgrade && \
       apk add --no-cache ca-certificates && \
       update-ca-certificates && \
       rm -rf /var/cache/apk/*
       
ENTRYPOINT ["/vault-init"]
