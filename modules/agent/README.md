urlooker-agent
============

agent会定时从web组件获取待监控url列表，发起模拟访问，然后将访问结果回报给web组件

## Installation

```bash
# set $GOPATH and $GOROOT
mkdir -p $GOPATH/src/github.com/urlooker
cd $GOPATH/src/github.com/urlooker
git clone https://github.com/710leo/urlooker/modules/agent.git
cd agent
go get ./...
./control build
./control start
```

## Configuration

```

{
    "debug": false,
    "hostname": "hostname.1", #hostname.1 和 web组件配置文件中monitorMap的值对应
    "worker": 1000, # 同时访问url的并发数
    "web": {
        "addrs": ["127.0.0.1:1985"],
        "interval": 60,
        "timeout": 1000
    }
}

```

