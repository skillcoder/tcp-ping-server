# Manual install on linux
In src dir:  
```
go build
sudo install -s -T tcp-ping-server /usr/local/bin/tcping
sudo cp tcping.service /etc/systemd/system/
sudo groupadd -r tcping
sudo useradd -r -g tcping tcping
sudo systemctl enable tcping.service
systemctl -l status tcping
sudo systemctl start tcping
```

Your service run on port 8888  
