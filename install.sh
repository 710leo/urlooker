if [ "x$GOPATH" = "x" ];then
	echo "GOPATH is not set, please set $GOPATH!"
	exit 1
fi

export GO111MODULE="off"
echo "set GO111MODULE=off"
mkdir -p $GOPATH/src/github.com/710leo
cd $GOPATH/src/github.com/710leo && git clone https://github.com/710leo/urlooker.git
cd $GOPATH/src/github.com/710leo/urlooker && ./control.sh build
echo "install ok! run cd $GOPATH/src/github.com/710leo/urlooker && ./control.sh start all"