# TC go server

HTTP server used as an interface to tc command. Can limit upload/download bandwidth, delay and packet loss.

### Use with golang
* Install [golang](https://golang.org/)

* Build:
```
make
```

* Run (may need root privilege):
```
./bin/tc-rest-controller -iu <interface name | default:eth0> -p <port | default:9010>
```
* For help:
```
./bin/tc-rest-controller -h
```

### Use with the docker image
* Install [docker](https://docs.docker.com/engine/installation/linux/)

* Run:
```
docker run --cap-add=ALL -d --net=host --name=tc-rest-controller mlacaud/tc-rest-controller -iu <interface name> -p <port | default:9010>
```

* For help:
```
docker run --cap-add=ALL -d --net=host --name=tc-rest-controller mlacaud/tc-rest-controller -f
```

### API
```
GET, DELETE /api/{upload,download}/limit
```

```
POST, PUT /api/{upload,download}/limit/{limit value in kbps}
```

```
GET, PUT,DELETE /api/{upload,download}/netem   PUT_JSON:{"delay":"50", "loss":"10"}
```

### Use case

At the beginning, use POST requests on **/api/{upload, download}/limit/{limit value in kbps}** before using **/api/{upload,download}/netem** because the netem is created and deleted in **/api/{upload, download}/limit**. The **DELETE** request on **/api/{upload,download}/netem** just reset the delay and percent of packet loss to the minimal values.

### Build the docker image


* Build the docker image:
```
make docker-build
```

* Build and push the docker image:
```
make docker-push
```

* Clean the go build:
```
make clean
```

* Clean the go dependencies:
```
make clean-go
```

* Clean the docker images:
```
make clean-docker
```

* Clean the go build, stop the containers and remove the docker images:
```
make mrproper
```

## License

MIT License

Copyright (c) 2017 mlacaud
