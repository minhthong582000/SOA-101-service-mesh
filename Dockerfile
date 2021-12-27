# Builder
FROM golang:1.17-alpine3.15 as builder

RUN apk update && apk upgrade && \
    apk --update add git make

WORKDIR /app

COPY . .

RUN make engine

# Distribution
FROM alpine:latest

ARG APP_USER="app"
ARG APP_GROUP="app"
ARG APP_UID="50000"
ARG APP_GID="50000"

USER root

RUN set -x \
    # create app user/group first, to be consistent throughout docker variants
    && addgroup -g $APP_GID -S $APP_GROUP \
    && adduser -S -D -H -u $APP_UID \
        -G $APP_GROUP \
        -g $APP_USER \
        $APP_USER

# Tini is installed as a helpful container entrypoint that reaps zombie
# processes and such of the actual executable we want to start, see
# https://github.com/krallin/tini#why-tini for details.
RUN apk add --no-cache --virtual build_deps tini &&  \
    cp /sbin/tini /usr/local/bin/tini && \
    apk del build_deps

RUN mkdir /app
WORKDIR /app

COPY --from=builder /app/bin/engine /app
RUN chmod +x /app/engine

USER $APP_USER
EXPOSE 8080
ENTRYPOINT ["tini", "-g", "--"]
CMD ["/app/engine"]