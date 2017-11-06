## [urlooker][1]
监控web服务可用性及访问质量，采用go语言编写，易于安装和二次开发

## Feature
- 返回状态码检测
- 页面响应时间检测
- 页面关键词匹配检测
- 带cookie访问
- agent多机房部署，指定机房访问
- 检测结果支持向open-falcon推送
- 支持短信和邮件告警

## Architecture
![此处输入图片的描述][2]

## ScreenShot

![看图][3]

![此处输入图片的描述][4]

![添加监控项][5]

## 常见问题
- [wiki手册][6]
- [常见问题][7]

## Install

**环境依赖**   
安装mysql & redis      
wget http://x2know.qiniudn.com/schema.sql      
将schema.sql 导入数据库  

二进制安装(Ubuntu 14.4 Go1.6下编译)：

    wget http://x2know.qiniudn.com/urlooker.tar.gz
    tar xzvf urlooker.tar.gz
    cd urlooker
    # 修改下cfg.json中的mysql和redis配置
    web/control start
    alarm/control start
    agent/control start

打开浏览器访问 http://127.0.0.1:1984 即可


源码安装及详细介绍见：   
web 组件[安装][8]   
agent 组件[安装][9]   
alarm 组件[安装][10]   

## 答疑
QQ交流群：556988374   

## Thanks
一些功能参考了open-falcon，感谢 [UlricQin][11] & [laiwei][12]


  [1]: https://github.com/urlooker
  [2]: https://github.com/urlooker/wiki/raw/master/img/urlooker4.png
  [3]: https://github.com/urlooker/wiki/raw/master/img/urlooker1.png
  [4]: https://github.com/urlooker/wiki/raw/master/img/urlooker3.png
  [5]: https://github.com/urlooker/wiki/raw/master/img/urlooker2.png
  [6]: https://github.com/URLooker/wiki
  [7]: https://github.com/URLooker/wiki/wiki/FAQ
  [8]: https://github.com/URLooker/web
  [9]: https://github.com/URLooker/agent
  [10]: https://github.com/URLooker/alarm
  [11]: http://ulricqin.com/
  [12]: https://github.com/laiwei