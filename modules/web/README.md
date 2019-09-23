urlooker-web
============

Web组件主要用于

- 监控项的添加
- 告警组人员管理
- 查看url访问质量绘图

## 常见问题
- [wiki说明][1]
- [常见问题][2]
- 初始用户名密码：admin/password

## Installation

```bash
# set $GOPATH and $GOROOT
mkdir -p $GOPATH/src/github.com/urlooker
cd $GOPATH/src/github.com/urlooker
git clone https://github.com/URLooker/web.git
cd web
./control build
./control start
```

## Configuration

```

    "debug": true,
    "salt": "have fun!",
    "admin":["admin"], #这里的用户是会变成admin用户
    "past": 30, #查看最近几分钟内的报警历史和绘图，默认为30分钟
    "http": {
        "listen": "0.0.0.0:1984",
        "secret": "secret"
    },
    "rpc": {
        "listen": "0.0.0.0:1985"
    },
    "mysql": {
        "addr": "root:123456@tcp(127.0.0.1:3306)/urlooker?charset=utf8&&loc=Asia%2FShanghai",
        "idle": 10,
        "max": 20
    },
    "alarm":{
        "enable": true,
        "batch": 200,
        "replicas": 500,
        "connTimeout": 1000,
        "callTimeout": 5000,
        "maxConns": 32,
        "maxIdle": 32,
        "sleepTime":30,
        "cluster":{
            "node-1":"127.0.0.1:1986"
        }
    },
    "monitorMap": { 
        "default":["hostname.1"], #监控指标多了之后agent地址可以填多个
    },
    "falcon":{
        "enable": false, # 为true表示向falcon推送数据
        "addr":"http://falcon.transfer.addr/api/push",
        "interval": 60
    },
    "internalDns":{ #通过公司内部接口获取url对应ip所在机房
        "enable": false,
        "addr":""
    }

```


  [1]: https://github.com/URLooker/wiki
  [2]: https://github.com/URLooker/wiki/wiki/FAQ