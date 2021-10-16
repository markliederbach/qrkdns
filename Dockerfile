################################
# STEP 1 build binary
################################
FROM golang:1.17-alpine as builder
ARG VERSION=latest

RUN adduser --uid 1500 -D qrkdns

WORKDIR $GOPATH/src/github.com/markliederbach/qrkdns/
COPY . /go/src/github.com/markliederbach/qrkdns

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build --ldflags "-s -w -X main.Version=${VERSION}" \
    -o /qrkdns cmds/qrkdns/main.go

############################
# STEP 2 build a tiny image
############################
FROM scratch

USER qrkdns
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /qrkdns /qrkdns

ENTRYPOINT ["/qrkdns"]
