FROM golang:1.16-alpine as build

ENV GO111MODULE=on \
    CGO_ENABLED=0

WORKDIR /build

COPY . .

RUN \
    go mod download \
    && \
    go build -a -tags netgo -ldflags '-w -extldflags "-static"' -o sabnzbd_prometheus /build/cmd/sabnzbd_prometheus/. \
    && \
    chmod +x sabnzbd_prometheus

FROM alpine:3.14
ENV PORT="8081"

RUN \
    apk add --no-cache \
        ca-certificates \
        tzdata \
        tini \
    && \
    addgroup -S sabnzbd_prometheus \
    && \
    adduser -S sabnzbd_prometheus -G sabnzbd_prometheus
USER sabnzbd_prometheus:sabnzbd_prometheus
COPY  --from=build /build/sabnzbd_prometheus /usr/local/bin/sabnzbd_prometheus

ENTRYPOINT [ "/sbin/tini", "--" ]
CMD [ "sabnzbd_prometheus" ]

LABEL org.opencontainers.image.source https://github.com/triggity/sabnzbd_prometheus