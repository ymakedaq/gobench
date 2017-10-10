#!/bin/bash -e
export GOROOT=/usr/bin/go
export PATH=$PATH:/usr/bin/go
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
BASENAME=`basename $DIR`
export GOPATH=$DIR
export GOBIN=$DIR/bin/ 

app=$BASENAME
conf=src/conf/app.conf
pidfile=$DIR/$BASENAME.pid
logfile=$DIR/$BASENAME.log

 function check_pid() {
    if [ -f $pidfile ];then
        pid=`cat $pidfile`
        if [ -n $pid ]; then
            running=`ps -p $pid|grep -v "PID TTY" |wc -l`
            return $running
        fi
    fi
    return 0
}

function build() {
    gofmt -w src/
    cd src/ 
    go build -o $BASENAME
    if [ $? -ne 0 ]; then
        exit $?
    fi
}
function pack() {
    build
    cd  $DIR
    rm  -rf src/logs && rm  -rf src/cache 
    tar zcvf $app.tar.gz control src/$app src/conf   src/views src/swagger  src/logs  src/cache 
}
function start() {
    check_pid
    running=$?
    if [ $running -gt 0 ];then
        echo -n "$app now is running already, pid="
        cat $pidfile
        return 1
    fi

    if ! [ -f $conf ];then
        echo "Config file $conf doesn't exist, creating one." 
    fi
    cd src/ 
    nohup  ./$BASENAME  >$logfile 2>&1 &
    sleep 1
    running=`ps -p $! | grep -v "PID TTY" | wc -l`
    if [ $running -gt 0 ];then
        echo $! > $pidfile
        echo "$app started..., pid=$!"
    else
        echo "$app failed to start."
        return 1
    fi


}
function killall() {
    pid=`cat $pidfile`
    ps -ef|grep $BASENAME|grep -v grep|awk '{print $2}'|xargs kill -9 
    rm -f $pidfile
    echo "$app killed..., pid=$pid"
}
function stop() {
    #ps -ef|grep $BASENAME|grep -v grep|awk '{print $2}'|xargs kill -9
    pid=`cat $pidfile`
    kill $pid
    rm -f $pidfile
    echo "$app stoped..., pid=$pid"
}
function restart() {
    stop
    sleep 1
    start 
}
function reload() { 
    pid=`cat $pidfile`
    kill -HUP $pid
    sleep 1
    newpid=`ps -ef|grep $BASENAME|grep -v grep|awk '{print $2}'` 
    echo "$app reload..., pid=$newpid"
    echo $newpid > $pidfile 
}
function status() {
    check_pid
    running=$?
    if [ $running -gt 0 ];then
        echo started
    else
        echo stoped
    fi
}
function run() {
   cd src/
   ./$BASENAME -docker
   #go run main.go 
}
function beerun() {
   cd src/
   bee run 
}

function tailf() {
   tail -f $logfile
}
function docs() {
   cd src/ 
   bee generate docs 
}

function sslkey() {
   cd src/conf/ssl
   ###CA:
   #私钥文件
   openssl genrsa -out ca.key 2048
   #数字证书
   openssl req -x509 -new -nodes -key ca.key -subj "/C=CN/ST=Guangdong/L=Shenzhen/O=Linc/OU=uworkapi.juanpi.com/CN=uworkapi.juanpi.com"  -days 36500 -out ca.crt
   ###Server:
   #私钥文件
   openssl genrsa -out server.key 2048
   openssl req -new -key server.key -subj "/C=CN/ST=Guangdong/L=Shenzhen/O=Linc/OU=uworkapi.juanpi.com/CN=uworkapi.juanpi.com"    -out server.csr
   #数字证书
   openssl x509 -req -in server.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out server.crt -days 36500
   ###Client:
   #私钥文件
   openssl genrsa -out client.key 2048
   openssl req -new -key client.key -subj "/C=CN/ST=Guangdong/L=Shenzhen/O=Linc/OU=uworkapi.juanpi.com/CN=uworkapi.juanpi.com"  -out client.csr
   #数字证书
   openssl x509 -req -in client.csr -CA ca.crt -CAkey ca.key -CAcreateserial -extfile client.ext -out client.crt -days 36500
}  

 
function help() {
    echo "$0 build|start|stop|kill|restart|reload|run|tail|docs|pack|beerun|sslkey"
}
if [ "$1" == "" ]; then
    help
elif [ "$1" == "build" ];then
    build
elif [ "$1" == "pack" ];then
    pack
elif [ "$1" == "start" ];then
    start
elif [ "$1" == "stop" ];then
    stop
elif [ "$1" == "kill" ];then
    killall
elif [ "$1" == "restart" ];then
    restart
elif [ "$1" == "reload" ];then
    reload
elif [ "$1" == "status" ];then
    status
elif [ "$1" == "run" ];then
    run 
elif [ "$1" == "beerun" ];then
    beerun 
elif [ "$1" == "tail" ];then
    tailf
elif [ "$1" == "docs" ];then
    docs 
elif [ "$1" == "sslkey" ];then
    sslkey
else
    help
fi
#bee api  gopub -driver=mysql -conn="mysql:mysql888@tcp(192.168.143.62:3306)/walle"
#bee api  uwrok-api -driver=mysql -conn="mysql:mysql888@tcp(127.0.0.1:3306)/wf_workflow1"
