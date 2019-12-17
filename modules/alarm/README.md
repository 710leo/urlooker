urlooker-alarm
============

alarm是用于判断是否触发报警条件的组件

alarm会定期从web端获取策略列表，接收到web端发送的检测数据后，对数据进行判断，若触发则产生event数据，将event数据存到 redis 中

## 发送短信
编辑 script/send.sms.sh 适配公司的短信网关