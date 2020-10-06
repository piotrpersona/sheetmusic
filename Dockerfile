FROM golang@sha256:2b3ca6f02d74eaf6f2d1788a16c1ccf551fe2407cb457636f3826f0108fed8ff AS stage-build

WORKDIR "/go/src/github.com/piotrpersona/sheetmusic"

RUN apk update && apk add dep git
ENV GO111MODULE=on
COPY go.mod .
RUN go mod donwload
COPY main.go .
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 \
    go build \
    -ldflags="-w -s" \
    -o /go/bin/package

FROM scratch

COPY --from=stage-build \
    /go/bin/package /usr/local/bin/package
