monicore
========
```Shell
cd /data
export GOPATH=/data
export PATH=$PATH:$GOPATH/bin
go get github.com/gorilla/mux
mkdir -p $GOPATH/src/github.com/playgrunge
cd $GOPATH/src/github.com/playgrunge
git clone https://github.com/playgrunge/monicore.git
go build
go run server.go
```
```Shell
http://192.171.0.120:3000/
```
