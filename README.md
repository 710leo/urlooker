## [urlooker][1]
监控web服务可用性及访问质量，采用go语言编写，易于安装和二次开发

## Feature
- 返回状态码检测
- 页面响应时间检测
- 页面关键词匹配检测
- 自定义Header
- GET、POST、PUT访问
- 自定义POST BODY
- 检测结果支持向open-falcon推送

## Architecture
![此处输入图片的描述][2]

## ScreenShot

![看图][3]

![此处输入图片的描述][4]

![添加监控项][5]

## 常见问题
- [wiki手册][6]
- [常见问题][7]
- 初始用户名密码：admin/password

## Install

#### 源码安装
###### 安装依赖
```
yum install -y redis
yum install -y mysql-server
```
###### 导入数据库
```
wget https://raw.githubusercontent.com/710leo/urlooker/master/sql/schema.sql
mysql -h 127.0.0.1 -u root -p < schema.sql
```

###### 安装组件
```bash
# set $GOPATH and $GOROOT
mkdir -p $GOPATH/src/github.com/710leo
cd $GOPATH/src/github.com/urlooker
git clone https://github.com/710leo/urlooker.git
go get ./...
./control build
./control start all
```

打开浏览器访问 http://127.0.0.1:1984 即可


## 答疑
QQ交流群：556988374

## Thanks
一些功能参考了open-falcon，感谢 [UlricQin][9] & [laiwei][10]


  [1]: https://github.com/urlooker
  [2]: https://github.com/urlooker/wiki/raw/master/img/urlooker4.png
  [3]: https://github.com/urlooker/wiki/raw/master/img/urlooker1.png
  [4]: https://github.com/urlooker/wiki/raw/master/img/urlooker3.png
  [5]: https://github.com/urlooker/wiki/raw/master/img/urlooker2.png
  [6]: https://github.com/URLooker/web
  [7]: https://github.com/URLooker/agent
  [8]: https://github.com/URLooker/alarm
  [9]: http://ulricqin.com/
  [10]: https://github.com/laiwei