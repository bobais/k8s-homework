FROM golang:1.12-alpine as builder

RUN apk update && apk add git

ENV GO_PROJECT="github.com/bobais/eventstats-watcher" CGO_ENABLED="0" GOOS="linux"

WORKDIR ${GOPATH}/src/${GO_PROJECT}
COPY . ./

RUN go get -d ${GO_PROJECT} \
    && go test \
    && go build -o /eventstats-watcher \
    && ls -lh /eventstats-watcher

FROM alpine:latest

# 65534=nobody
COPY --chown=65534 --from=builder /eventstats-watcher /
EXPOSE 10080

USER 65534
CMD ["/eventstats-watcher"]
