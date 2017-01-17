FROM golang:1.6.3-alpine
MAINTAINER mlacaud@viotech.net

RUN apk update
RUN apk add --no-cache iproute2
RUN apk add --no-cache git

RUN go get github.com/gorilla/mux


ADD tcserver.go tcserver.go

ENTRYPOINT ["go", "run", "tcserver.go"]
