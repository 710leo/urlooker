## [urlooker](https://github.com/710leo/urlooker)
监控web服务可用性及访问质量，采用go语言编写，易于安装和二次开发    
[English](https://github.com/710leo/urlooker/readme.md)|[中文](https://github.com/710leo/urlooker/blob/master/readme_zh.md)

## Feature
- 返回状态码检测
- 页面响应时间检测
- 页面关键词匹配检测
- 自定义Header
- GET、POST、PUT访问
- 自定义POST BODY
- 检测结果支持向open-falcon推送

## Architecture
![Architecture](img/urlooker_arch.png)

## ScreenShot

![ScreenShot](img/urlooker1.png)
![ScreenShot](img/urlooker2.png)
<img src="img/urlooker_stra.png" style="zoom:45%;" />

## 常见问题
- [wiki手册](https://github.com/710leo/urlooker/wiki)
- [常见问题](https://github.com/710leo/urlooker/wiki/FAQ)
- 初始用户名密码：admin/password

## Install
##### 安装依赖
```bash
yum install -y redis
yum install -y mysql-server
```
##### 导入数据库
```bash
wget https://raw.githubusercontent.com/710leo/urlooker/master/sql/schema.sql
mysql -h 127.0.0.1 -u root -p < schema.sql
```

##### 安装组件
```bash
# set $GOPATH and $GOROOT
mkdir -p $GOPATH/src/github.com/710leo
cd $GOPATH/src/github.com/710leo
git clone https://github.com/710leo/urlooker.git
go get ./...
./control build
./control start all
```

打开浏览器访问 http://127.0.0.1:1984 即可


## 答疑
知识星球   
<img src="img/urlooker_zsxq.png" style="zoom:50%;" />

## Thanks
一些功能参考了open-falcon，感谢 [UlricQin](http://ulricqin.com) & [laiwei](https://github.com/laiwei)