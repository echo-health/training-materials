FROM golang:1.13-alpine AS build
WORKDIR /goroutine
COPY . /goroutine
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -mod vendor -o build/goroutine .


FROM scratch
COPY --chown=0:0 --from=build /goroutine/build/goroutine /go/bin/goroutine
ENTRYPOINT ["/go/bin/goroutine"]