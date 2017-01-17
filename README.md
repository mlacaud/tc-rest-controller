# TC go server

HTTP server used as an interface to tc command. Can limit upload/download bandwidth, delay and packet loss.

### Use with golang
* Install [golang](https://golang.org/)

* Add dependencies:
```
go get
```
* Run (may need root privilege):
```
go run tc-rest-controller.go -iu <interface name> -p <port | default:9010>
```
* For help:
```
go run tc-rest-controller.go -h
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
First, **add dependencies:**
```
go get
```
Then, use the **Makefile to**:

* Compile the go code:
```
make build
```

* Build the docker image:
```
make
```

* Build and push
```
make docker-push
```

* Clean the go build:
```
make clean
```

* Clean the go build, stop the containers and remove the docker images:
```
make mrproper
```

## License

MIT License

Copyright (c) 2017 mlacaud
