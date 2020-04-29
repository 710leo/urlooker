if [ "x$GOPATH" = "x" ];then
	echo "GOPATH is not set, please install golang & set $GOPATH!"
	exit 1
fi

export GO111MODULE="off"
echo "set GO111MODULE=off"
mkdir -p $GOPATH/src/github.com/710leo
   
if [ ! -d "$GOPATH/src/github.com/710leo/urlooker" ]; then
	cd $GOPATH/src/github.com/710leo && git clone https://github.com/710leo/urlooker.git
else
	cd $GOPATH/src/github.com/710leo/urlooker && git pull
fi
cd $GOPATH/src/github.com/710leo/urlooker && ./control build
echo "install ok! run cd $GOPATH/src/github.com/710leo/urlooker && ./control start all"