# TC go server

HTTP server used as an interface to tc command. Can limit upload/download bandwidth, delay and packet loss.

### Use

```
docker run --cap-add=ALL -d --net=host --name=tcserver msstream/tcserver -iu <interface name> -p <port | default:9010>
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
