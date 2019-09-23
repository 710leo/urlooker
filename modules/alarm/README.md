urlooker-alarm
============

alarm是用于判断是否触发报警条件的组件

alarm会定期从web端获取策略列表，接收到web端发送的检测数据后，对数据进行判断，若触发则产生event数据，将event数据存到 redis 中

## Installation

```bash
# set $GOPATH and $GOROOT
mkdir -p $GOPATH/src/github.com/urlooker
cd $GOPATH/src/github.com/urlooker
git clone https://github.com/URLooker/alarm.git
cd alarm
./control build
./control start
```

## 发送短信
编辑 script/send.sms.sh 适配公司的短信网关

## Configuration

```

{
    "debug": false,
	"remain":10,  #配置策略中支持的最大连续次数
	"rpc":{
		"listen":"0.0.0.0:1986"
	},
    "web": {
        "addrs": ["127.0.0.1:1985"], #可以填多个web地址
        "timeout": 300,
        "interval": 60
    },
    "alarm": {
        "enabled": true,
        "minInterval": 180,
        "queuePattern": "event",
        "redis": {
            "dsn": "127.0.0.1:6379",
            "maxIdle": 5,
            "connTimeout": 20000,
            "readTimeout": 20000,
            "writeTimeout": 20000
        }
    },
    "queue": {
        "sms": "/sms",
        "mail": "/mail"
    },
    "worker": {
        "sms": 10,
        "mail": 50
    },
    "smtp": {
        "addr": "mail.addr:25",
        "username": "mail@mail.com",
        "password": "",
        "from": "mail@mail.com"
    }
}


```

