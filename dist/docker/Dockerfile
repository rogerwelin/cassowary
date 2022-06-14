FROM alpine:3.12.4
RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*
COPY cassowary /usr/bin/cassowary
ENTRYPOINT ["cassowary", "run"]
