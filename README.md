## [urlooker](https://github.com/710leo/urlooker)
enterprise-level websites monitoring system
[English](https://github.com/710leo/urlooker/readme.md)|[中文](https://github.com/710leo/urlooker/readmd_zh.md)

## Feature
- status code
- respose time
- page keyword 
- customize header
- customize post body
- support get post put method
- send to open-falcon

## Architecture
![Architecture](img/urlooker_arch.png)

## ScreenShot

![](img/urlooker_en1.png)
![](img/urlooker_en2.png)
<img src="img/urlooker_stra.png" style="zoom:45%;" />

## FAQ
- [wiki](https://github.com/710leo/urlooker/wiki)
- [FAQ](https://github.com/710leo/urlooker/wiki/FAQ)
- default user/password：admin/password

## Install
##### dependence
```
yum install -y redis
yum install -y mysql-server
```
##### import mysql database
```
wget https://raw.githubusercontent.com/710leo/urlooker/master/sql/schema.sql
mysql -h 127.0.0.1 -u root -p < schema.sql
```

##### install modules
```bash
# set $GOPATH and $GOROOT
mkdir -p $GOPATH/src/github.com/710leo
cd $GOPATH/src/github.com/710leo
git clone https://github.com/710leo/urlooker.git
go get ./...
./control build
./control start all
```

open http://127.0.0.1:1984 in browse

## Q&A
Gitter: [urlooker](https://gitter.im/urllooker/community
