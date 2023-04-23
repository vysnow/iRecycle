FROM golang:alpine as builder

WORKDIR /gowordspace/src/mego/
COPY . .

RUN go env -w GO111MODULE=on \
    && go env -w GOPROXY=https://goproxy.cn,direct \
    && go env -w CGO_ENABLED=0 \
    && go env \
    && go mod tidy \
    && go build -o mego .

FROM alpine:latest
LABEL MAINTAINER="cqins@163.com"

WORKDIR /gowordspace/src/mego/

COPY --from=0 /gowordspace/src/mego ./
COPY --from=0 /gowordspace/src/mego/resource ./resource/
COPY --from=0 /gowordspace/src/mego/config.yaml ./

EXPOSE 8080
ENTRYPOINT ./mego -c config.yaml
