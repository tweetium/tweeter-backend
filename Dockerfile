FROM golang:1.11

WORKDIR /app

RUN go get gotest.tools/gotestsum
RUN curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s -- -b $GOPATH/bin v1.15.0

COPY scripts/.bashrc /root/.bashrc
