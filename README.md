# TCP Ping Server

### GET
```
mkdir -p ~/go/src/github.com/skillcoder
cd ~/go/src/github.com/skillcoder
git clone https://github.com/skillcoder/tcp-ping-server.git
```

### BUILD
#### Amazon Linux 2 (AWS)
```
sudo yum install glide
cd ~/go/src/github.com/skillcoder/tcp-ping-server
glide update
go build
```

### RUN
```
setenv TCPPINGSERVER_SERVICE_LISTEN 127.0.0.1:8888
./tcp-ping-server
```
