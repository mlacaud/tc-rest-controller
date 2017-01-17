FROM alpine
MAINTAINER mlacaud@msstream.net

RUN apk update
RUN apk add --no-cache iproute2

ADD bin/tc-rest-controller /usr/bin/tc-rest-controller

ENTRYPOINT ["tc-rest-controller"]
