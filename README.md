monicore
========
```Shell
cd /data
export GOPATH=/data
export PATH=$PATH:$GOPATH/bin
go get "github.com/gorilla/mux"
go get "github.com/codegangsta/gin"
mkdir -p $GOPATH/src/github.com/playgrunge
cd $GOPATH/src/github.com/playgrunge
git clone https://github.com/playgrunge/monicore.git
cd monicore
gin -p 3001 -a 3000 r server.go
```
```Shell
http://192.171.0.120:3001/
```
