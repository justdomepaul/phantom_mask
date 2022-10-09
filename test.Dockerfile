FROM golang:1.19.2-buster AS base
ENV GO111MODULE=on
ADD ./go.mod ./go.sum /workspace/
WORKDIR /

WORKDIR /workspace
RUN go mod download

COPY . /workspace
