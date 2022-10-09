FROM golang:1.19.2-buster AS base
ENV GO111MODULE=on
ADD ./go.mod ./go.sum /workspace/

WORKDIR /workspace
RUN go mod download

FROM base AS build
ADD . /workspace
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go install ./cmd/restful

FROM debian:bullseye-slim AS runtime
RUN apt-get update && apt-get install -y ca-certificates wget
ENV DOCKERIZE_VERSION v0.6.1
RUN wget https://github.com/jwilder/dockerize/releases/download/$DOCKERIZE_VERSION/dockerize-linux-amd64-$DOCKERIZE_VERSION.tar.gz \
    && tar -C /usr/local/bin -xzvf dockerize-linux-amd64-$DOCKERIZE_VERSION.tar.gz \
    && rm dockerize-linux-amd64-$DOCKERIZE_VERSION.tar.gz \
    && rm -rf /var/lib/apt/lists/*

FROM runtime
RUN mkdir /app
COPY --from=build /go/bin/restful /app/restful
WORKDIR /app
ENTRYPOINT ["dockerize"]
CMD ["bash"]
