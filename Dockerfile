FROM golang:1.11

WORKDIR /app

RUN go get gotest.tools/gotestsum
COPY scripts/.bashrc /root/.bashrc
