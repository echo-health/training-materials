FROM golang:1.13-alpine AS build
WORKDIR /support
COPY . /support
RUN apk update && apk upgrade && \
    apk add --no-cache bash git openssh build-base ca-certificates tzdata
RUN echo "nobody:x:65534:65534:Nobody:/:" > /etc/passwd
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -mod vendor -o build/hello-world .


FROM scratch
COPY --from=build /etc/passwd /etc/passwd
COPY --from=build /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --chown=0:0 --from=build /support/build/hello-world /go/bin/hello-world
USER nobody
ENTRYPOINT ["/go/bin/hello-world"]
CMD ["up"]